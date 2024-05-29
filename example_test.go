package environment_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/x-ethr/environment"
)

func Example() {
	ctx, level := context.Background(), slog.LevelInfo

	// Log all environment variables
	environment.Log(ctx, level)

	// Log only specific environment variable(s)
	environment.Log(ctx, level, func(o *environment.Options) {
		o.Variables = []string{
			"PATH",
		}
	})

	// Log only specific environment variable(s), and warn if one of the specified variable(s) was set to an empty string
	environment.Log(ctx, level, func(o *environment.Options) {
		_ = os.Setenv("TEST", "") // --> for example purposes

		o.Variables = []string{
			"TEST",
		}

		o.Warnings.Empty = true
	})

	// Log only specific environment variable(s), and warn if one of the specified variable(s) wasn't found
	environment.Log(ctx, level, func(o *environment.Options) {
		_ = os.Unsetenv("TEST") // --> for example purposes

		o.Variables = []string{
			"TEST",
		}

		o.Warnings.Missing = true
	})
}
