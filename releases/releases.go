package releases

import (
	"fmt"
	"strconv"
	"time"

	"github.com/remind101/empire/configs"
	"github.com/remind101/empire/repos"
	"github.com/remind101/empire/slugs"
)

// Release is a combination of a Config and a Slug, which form a deployable
// release.
type Release struct {
	ID        string
	Repo      repos.Repo
	Version   string
	Config    *configs.Config
	Slug      *slugs.Slug
	CreatedAt time.Time
}

// ReleaseRepository is an interface that can be implemented for storing and
// retrieving releases.
type Repository interface {
	Create(repos.Repo, *configs.Config, *slugs.Slug) (*Release, error)
	FindByRepo(repos.Repo) ([]*Release, error)
	FindByReleaseID(string) (*Release, error)
	Head(repos.Repo) (*Release, error)
}

// repository is an in-memory implementation of a Repository
type repository struct {
	byRepo       map[repos.Repo][]*Release
	byReleaseID  map[string]*Release
	versions     map[repos.Repo]int
	genTimestamp func() time.Time
	id           int
}

// Create a new repository
func newRepository() *repository {
	return &repository{
		byRepo:      make(map[repos.Repo][]*Release),
		byReleaseID: make(map[string]*Release),
		versions:    make(map[repos.Repo]int),
	}
}

// Generates a repository that stubs out the CreatedAt field.
func newFakeRepository() *repository {
	r := newRepository()
	r.genTimestamp = func() time.Time {
		return time.Date(2014, time.January, 1, 0, 0, 0, 0, time.UTC)
	}
	return r
}

func (p *repository) Create(repo repos.Repo, config *configs.Config, slug *slugs.Slug) (*Release, error) {
	p.id++

	createdAt := time.Now()
	if p.genTimestamp != nil {
		createdAt = p.genTimestamp()
	}

	version := 1
	if v, ok := p.versions[repo]; ok {
		version = v
	}

	r := &Release{
		ID:        strconv.Itoa(p.id),
		Repo:      repo,
		Version:   fmt.Sprintf("v%d", version),
		Config:    config,
		Slug:      slug,
		CreatedAt: createdAt.UTC(),
	}

	p.versions[repo] = version + 1
	p.byRepo[r.Repo] = append(p.byRepo[r.Repo], r)
	p.byReleaseID[r.ID] = r

	return r, nil
}

func (p *repository) FindByRepo(repo repos.Repo) ([]*Release, error) {
	if set, ok := p.byRepo[repo]; ok {
		return set, nil
	}

	return []*Release{}, nil
}

func (p *repository) FindByReleaseID(releaseID string) (*Release, error) {
	r, ok := p.byReleaseID[releaseID]
	if !ok {
		r = &Release{}
	}

	return r, nil
}

func (p *repository) Head(repo repos.Repo) (*Release, error) {
	set, ok := p.byRepo[repo]
	if !ok {
		return nil, nil
	}

	return set[len(set)-1], nil
}

type Service struct {
	Repository
}

func (s *Service) Create(config *configs.Config, slug *slugs.Slug) (*Release, error) {
	return s.Repository.Create(slug.Image.Repo, config, slug)
}