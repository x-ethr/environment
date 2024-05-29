package environment

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// Log logs environment variables with the specified log level. It iterates over the environment variables,
// splits them into key-value pairs, and logs them using slog.Log.
func Log(ctx context.Context, level slog.Level, settings ...Variadic) {
	var o = Settings()
	for _, configuration := range settings {
		configuration(o)
	}

	env := os.Environ()

	var variables = make(map[string]string, len(env))
	for index := range env {
		variable := strings.SplitN(env[index], "=", 2)
		if len(variable) != 2 {
			variables[env[index]] = ""
			continue
		}

		variables[strings.TrimSpace(variable[0])] = strings.TrimSpace(variable[1])
	}

	if len(o.Variables) > 0 {
		for index := range o.Variables {
			key := o.Variables[index]
			if v, found := variables[key]; found && v != "" {
				slog.Log(ctx, level, "Environment Variable", slog.String("key", key), slog.String("value", v))
			} else if found && o.Warnings.Empty {
				slog.WarnContext(ctx, "Environment Variable Wasn't Set", slog.String("key", key))
			} else if !(found) && o.Warnings.Missing {
				slog.WarnContext(ctx, "Environment Variable Wasn't Found", slog.String("key", key))
			}
		}

		return
	}

	for key, value := range variables {
		slog.Log(ctx, level, "Environment Variable(s)", slog.String("key", key), slog.String("value", value))
	}
}
