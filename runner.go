package empire

import (
	"io"

	"golang.org/x/net/context"
)

const (
	// GenericProcessName is the process name for `emp run` processes not defined in the procfile.
	GenericProcessName = "run"
)

// RunRecorder is a function that returns an io.Writer that will be written to
// to record Stdout and Stdin of interactive runs.
type RunRecorder func() (io.Writer, error)

type runnerService struct {
	*Empire
}

func (r *runnerService) Run(ctx context.Context, opts RunOpts) error {
	app := opts.App

	procName := opts.Command[0]
	var proc Process

	// First, let's check if the command we're running matches a defined
	// process in the Procfile/Formation. If it does, we'll replace the
	// command, with the one in the procfile and expand it's arguments.
	//
	// For example, given a procfile like this:
	//
	//	psql:
	//	  command: ./bin/psql
	//
	// Calling `emp run psql DATABASE_URL` will expand the command to
	// `./bin/psql DATABASE_URL`.
	if cmd, ok := app.Formation[procName]; ok {
		proc = cmd
		proc.Command = append(cmd.Command, opts.Command[1:]...)
		proc.NoService = false
	} else {
		// If we've set the flag to only allow `emp run` on commands
		// defined in the procfile, return an error since the command is
		// not defined in the procfile.
		if r.AllowedCommands == AllowCommandProcfile {
			return commandNotInFormation(Command{procName}, app.Formation)
		}

		// This is an unnamed command, fallback to a generic proc name.
		procName = GenericProcessName
		proc.Command = opts.Command
		proc.SetConstraints(DefaultConstraints)
	}

	proc.Quantity = 1

	// Set the size of the process.
	if opts.Constraints != nil {
		proc.SetConstraints(*opts.Constraints)
	}

	app.Formation = Formation{procName: proc}

	return r.engine.Run(ctx, app, opts.IO)
}
