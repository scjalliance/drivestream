package drivestream

import (
	"fmt"
	"io"
	"time"
)

// newTaskLogger returns a task logger that will write to w.
func newTaskLogger(w io.Writer) taskLogger {
	return taskLogger{out: w}
}

// A taskLogger writes logs for a task.
type taskLogger struct {
	out   io.Writer
	task  string
	start time.Time
}

// Task returns a logger for the given subtask.
func (l taskLogger) Task(task string) taskLogger {
	return taskLogger{
		out:   l.out,
		task:  l.task + task + ": ",
		start: time.Now(),
	}
}

// Log formats according to a format specifier and writes to the logger's
// output. If the stream has no log, nothing is printed.
func (l taskLogger) Log(format string, a ...interface{}) {
	if l.out != nil {
		fmt.Fprintf(l.out, "%s%s", l.task, fmt.Sprintf(format, a...))
	}
}

// Duration returns the duration of the task so far.
func (l taskLogger) Duration() time.Duration {
	return time.Now().Sub(l.start)
}
