package hiking

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	os.RemoveAll("testdir")
	if err := os.Mkdir("testdir", 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Create(filepath.Join("testdir", "a.txt")); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("testdir")

	tests := []struct {
		name    string
		workDir string
		command string
		args    []string
	}{
		{
			name:    "ls",
			workDir: "testdir",
			command: "ls",
			args:    nil,
		},
		{
			name:    "ls -a",
			workDir: "testdir",
			command: "ls",
			args:    []string{"-a"},
		},
		{
			name:    "file not found",
			workDir: "testdir",
			command: "ls",
			args:    []string{"-1", "xxxxx.txt"},
		},
		{
			name:    "workdir not exists",
			workDir: "xxxxxx",
			command: "ls",
			args:    nil,
		},
		{
			name:    "command timeout",
			workDir: "testdir",
			command: "sleep",
			args:    []string{"3"},
		},
		{
			name:    "command cancel",
			workDir: "testdir",
			command: "sleep",
			args:    []string{"3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				stdout bytes.Buffer
				stderr bytes.Buffer
			)
			o := NewOSCommand(tt.workDir, tt.command, tt.args, &stdout, &stderr)

			switch tt.name {
			case "ls":
				err := o.Run(context.Background())

				assert.Nil(t, err)
				assert.Equal(t, []string{"a.txt"},
					strings.Split(strings.TrimSpace(stdout.String()), "\n"))
				assert.Equal(t, "", stderr.String())
			case "ls -a":
				err := o.Run(context.Background())

				assert.Nil(t, err)
				assert.Equal(t, []string{".", "..", "a.txt"},
					strings.Split(strings.TrimSpace(stdout.String()), "\n"))
				assert.Equal(t, "", stderr.String())
			case "file not found":
				err := o.Run(context.Background())

				assert.NotNil(t, err)
				assert.Equal(t, "", stdout.String())
				assert.Equal(t, "ls: xxxxx.txt: No such file or directory\n",
					stderr.String())
			case "workdir not exists":
				err := o.Run(context.Background())

				assert.NotNil(t, err)
				assert.Equal(t, "", stdout.String())
				assert.Equal(t, "", stderr.String())
			case "command timeout":
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
				defer cancel()

				err := o.Run(ctx)

				assert.NotNil(t, err)
				assert.Equal(t, "", stdout.String())
				assert.Equal(t, "", stderr.String())
			case "command cancel":
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				err := o.Run(ctx)

				assert.NotNil(t, err)
				assert.Equal(t, "", stdout.String())
				assert.Equal(t, "", stderr.String())
			default:
				t.Fatal("no such testing")
			}
		})
	}
}
