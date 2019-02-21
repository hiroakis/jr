package jr

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadJob(t *testing.T) {
	jobFile := `
jobs:
  - name: The job name
    tasks:
      - name: Run batch.bat
        workdir: c:/BatchHome/bin
        command: batch.bat
        timeout: 30s
      - name: Run some-command
        workdir: .
        command: some-command.exe
        args:
          - -opt1 opt1
          - -opt2
          - -opt3=opt3
        timeout: 30s
        env:
          no_proxy: 169.254.169.254
          http_proxy: http://proxy.com:8080
          https_proxy: http://proxy.com:8080
`

	jobs, err := LoadJob(strings.NewReader(jobFile))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(jobs.Jobs))

	job0 := jobs.Jobs[0]
	assert.Equal(t, 2, len(job0.Tasks))
	assert.Equal(t, "The job name", job0.Name)

	task0 := job0.Tasks[0]
	task1 := job0.Tasks[1]
	assert.Equal(t, "Run batch.bat", task0.Name)
	assert.Equal(t, "c:/BatchHome/bin", task0.WorkDir)
	assert.Equal(t, "batch.bat", task0.Command)
	assert.Nil(t, task0.Args)
	assert.Equal(t, 30*time.Second, task0.Timeout)
	assert.Equal(t, 0, len(task0.Env))

	assert.Equal(t, "Run some-command", task1.Name)
	assert.Equal(t, ".", task1.WorkDir)
	assert.Equal(t, "some-command.exe", task1.Command)
	assert.Equal(t, []string{"-opt1 opt1", "-opt2", "-opt3=opt3"}, task1.Args)
	assert.Equal(t, 30*time.Second, task1.Timeout)
	assert.Equal(t, 3, len(task1.Env))
	assert.Equal(t, "169.254.169.254", task1.Env["no_proxy"])
	assert.Equal(t, "http://proxy.com:8080", task1.Env["http_proxy"])
	assert.Equal(t, "http://proxy.com:8080", task1.Env["https_proxy"])
}

func TestTaskRun(t *testing.T) {
	tests := []struct {
		name string
		task Task
	}{
		{
			name: "run echo command",
			task: Task{
				Name:    "echo command",
				WorkDir: ".",
				Command: "echo",
				Args:    []string{"-n", "xxxx"},
				Timeout: 1 * time.Second,
			},
		},
		{
			name: "run echo command",
			task: Task{
				Name:    "echo command",
				WorkDir: ".",
				Command: "echo",
				Args:    []string{"-n", "xxxx"},
				Timeout: 1 * time.Second,
				Env: map[string]string{
					"no_proxy":    "169.254.169.254",
					"http_proxy":  "http://proxy.prd.proc.kanmu:8080",
					"https_proxy": "http://proxy.prd.proc.kanmu:8080",
				},
			},
		},
		{
			name: "command timeout",
			task: Task{
				Name:    "sleep 3 second",
				WorkDir: ".",
				Command: "sleep",
				Args:    []string{"3"},
				Timeout: 100 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "run echo command":
				err := tt.task.Run()
				assert.Nil(t, err)
				assert.Equal(t, "xxxx", tt.task.Stdout.String())
				assert.Equal(t, "", tt.task.Stderr.String())
			case "command timeout":
				err := tt.task.Run()
				assert.NotNil(t, err)
			default:
				t.Fatal("no such testing")
			}
		})
	}
}
