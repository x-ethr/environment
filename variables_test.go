package environment_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/x-ethr/environment"
)

func Test(t *testing.T) {
	ctx, level := context.Background(), slog.LevelInfo
	slog.SetLogLoggerLevel(level)

	t.Run("Log", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			var buffer bytes.Buffer
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: level})))
			t.Setenv("TEST", "example-value-1")

			// Log all environment variables
			environment.Log(ctx, level)

			if size := buffer.Len(); size > 0 {
				t.Logf("Success: Output Size - %d", size)
			}
		})

		t.Run("Options-Variables", func(t *testing.T) {
			var buffer bytes.Buffer
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: level})))
			t.Setenv("TEST", "example-value")

			// Log only specific environment variable(s)
			environment.Log(ctx, level, func(o *environment.Options) {
				o.Variables = []string{
					"TEST",
				}
			})

			var mapping map[string]interface{}
			if e := json.NewDecoder(&buffer).Decode(&mapping); e != nil {
				t.Fatalf("Error Decoding SLOG Log Entry: %v", e)
			}

			content, e := json.MarshalIndent(mapping, "", "    ")
			if e != nil {
				t.Fatalf("Unable to marshal mapping: %v", e)
			}

			if mapping["key"] != "TEST" {
				t.Errorf("Expected Key (%s), Received: %s", "TEST", string(content))
			} else if mapping["value"] != "example-value" {
				t.Errorf("Expected Value (%s), Received: %s", "example-value", string(content))
			} else {
				t.Logf("Success: %s", string(content))
			}
		})

		t.Run("Options-Variables-Empty", func(t *testing.T) {
			var buffer bytes.Buffer
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: level})))
			t.Setenv("TEST", "")

			// Log only specific environment variable(s), and warn if one of the specified variable(s) was set to an empty string
			environment.Log(ctx, level, func(o *environment.Options) {
				o.Variables = []string{
					"TEST",
				}

				o.Warnings.Empty = true
			})

			var mapping map[string]interface{}
			if e := json.NewDecoder(&buffer).Decode(&mapping); e != nil {
				t.Fatalf("Error Decoding SLOG Log Entry: %v", e)
			}

			content, e := json.MarshalIndent(mapping, "", "    ")
			if e != nil {
				t.Fatalf("Unable to marshal mapping: %v", e)
			}

			if mapping["key"] != "TEST" {
				t.Errorf("Expected Key (%s), Received: %s", "TEST", string(content))
			} else if mapping["msg"] != "Environment Variable Wasn't Set" {
				t.Errorf("Expected Message (\"%s\"), Received: %s", "Environment Variable Wasn't Set", string(content))
			} else if mapping["level"] != "WARN" {
				t.Errorf("Expected Log Level Warning, Received: %s", mapping["level"])
			} else {
				t.Logf("Success: %s", string(content))
			}
		})

		t.Run("Options-Variables-Missing", func(t *testing.T) {
			var buffer bytes.Buffer
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: level})))

			// Log only specific environment variable(s), and warn if one of the specified variable(s) was set to an empty string
			environment.Log(ctx, level, func(o *environment.Options) {
				o.Variables = []string{
					"TEST",
				}

				o.Warnings.Missing = true
			})

			var mapping map[string]interface{}
			if e := json.NewDecoder(&buffer).Decode(&mapping); e != nil {
				t.Fatalf("Error Decoding SLOG Log Entry: %v", e)
			}

			content, e := json.MarshalIndent(mapping, "", "    ")
			if e != nil {
				t.Fatalf("Unable to marshal mapping: %v", e)
			}

			if mapping["key"] != "TEST" {
				t.Errorf("Expected Key (%s), Received: %s", "TEST", string(content))
			} else if mapping["msg"] != "Environment Variable Wasn't Found" {
				t.Errorf("Expected Message (\"%s\"), Received: %s", "Environment Variable Wasn't Found", string(content))
			} else if mapping["level"] != "WARN" {
				t.Errorf("Expected Log Level Warning, Received: %s", mapping["level"])
			} else {
				t.Logf("Success: %s", string(content))
			}
		})
	})
}
