package hiking

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// LoadJob loads jobs from job.yml
func LoadJob(jobFile io.Reader) (JobFile, error) {
	b := new(bytes.Buffer)
	io.Copy(b, jobFile)

	var jobs JobFile
	if err := yaml.Unmarshal(b.Bytes(), &jobs); err != nil {
		return JobFile{}, errors.Wrap(err, "yaml.Unmarshal failed")
	}
	return jobs, nil
}

// JobFile jobfile
type JobFile struct {
	Jobs []Job `yaml:"jobs"`
}

// Job job
type Job struct {
	Name   string `yaml:"name"`
	Tasks  []Task `yaml:"tasks"`
	Stderr *log.Logger
	Stdout *log.Logger
}

var (
	defaultStdoutLogger = log.New(os.Stdout, "[jr stdout] ", log.LstdFlags|log.LUTC|log.Lmicroseconds)
	defaultStderrLogger = log.New(os.Stderr, "[jr stderr] ", log.LstdFlags|log.LUTC|log.Lmicroseconds)
)

// Run runs a job
func (job Job) Run() error {
	st := time.Now()

	if job.Stdout == nil {
		job.Stdout = defaultStdoutLogger
	}
	if job.Stderr == nil {
		job.Stderr = defaultStderrLogger
	}

	job.Stderr.Printf("START: %s", job.Name)
	for _, task := range job.Tasks {
		job.Stderr.Printf("The task is INPROGRESS: %s", task.Name)
		if err := task.Run(); err != nil {
			job.Stderr.Printf("FAILED: %s, error: %v", task.Name, err)
			return errors.Wrap(err, "task.run failed")
		}

		job.Stderr.Printf("The task is DONE: %s, %s", task.Name, task.Stderr)
	}
	job.Stderr.Printf("FINISH: %s", job.Name)

	job.Stdout.Printf("%s total time: %v", job.Name, time.Now().Sub(st))
	return nil
}

// Task task
type Task struct {
	Name    string            `yaml:"name"`
	WorkDir string            `yaml:"workdir"`
	Command string            `yaml:"command"`
	Args    []string          `yaml:"args"`
	Env     map[string]string `yaml:"env"`
	Timeout time.Duration     `yaml:"timeout"`
	Stdout  *bytes.Buffer
	Stderr  *bytes.Buffer
}

const defaultTaskTimeout = 30 * time.Second

func (t *Task) Run() error {
	for k, v := range t.Env {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range t.Env {
			os.Unsetenv(k)
		}
	}()

	if t.Timeout == 0 {
		t.Timeout = defaultTaskTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	t.Stdout = new(bytes.Buffer)
	t.Stderr = new(bytes.Buffer)

	command := NewOSCommand(t.WorkDir, t.Command, t.Args, t.Stdout, t.Stderr)
	if err := command.Run(ctx); err != nil {
		return errors.Wrap(err, "command.Run failed")
	}
	return nil
}
