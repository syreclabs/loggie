// Package loggie provides simple unstructured leveled logging.
package loggie

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	colorable "github.com/mattn/go-colorable"
	isatty "github.com/mattn/go-isatty"
)

// Known logging levels. Setting logging level to Lnone silences all logging output
// (including Fatal and Panic logging).
const (
	Ldebug = iota
	Linfo
	Lwarning
	Lerror
	Lfatal
	Lpanic
	Lnone
)

// Logger is the logging interface.
type Logger interface {
	Print(lvl int, v ...interface{})
	Printf(lvl int, format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Level() int
	SetLevel(lvl int)
}

type defaultLogger struct {
	name  string
	level int
	f     Formatter
	mu    sync.Mutex
	w     io.Writer
}

// NewLogger creates a new named logger writing to w and using formatter f.
func NewLogger(w io.Writer, name string, f Formatter) Logger {
	return &defaultLogger{w: w, name: name, f: f}
}

// New returns new named os.Stdout logger with default text formatter.
func New(name string) Logger {
	return NewLogger(os.Stdout, name, DefaultFormatter)
}

func (l *defaultLogger) Print(lvl int, v ...interface{}) {
	if lvl >= l.level {
		l.write(lvl, v...)
		if lvl > Lwarning {
			l.writeStack(lvl)
		}
	}
}

func (l *defaultLogger) Printf(lvl int, format string, v ...interface{}) {
	if lvl >= l.level {
		l.writef(lvl, format, v...)
		if lvl > Lwarning {
			l.writeStack(lvl)
		}
	}
}

func (l *defaultLogger) Debug(v ...interface{}) {
	l.Print(Ldebug, v...)
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	l.Printf(Ldebug, format, v...)
}

func (l *defaultLogger) Info(v ...interface{}) {
	l.Print(Linfo, v...)
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	l.Printf(Linfo, format, v...)
}

func (l *defaultLogger) Warning(v ...interface{}) {
	l.Print(Lwarning, v...)
}

func (l *defaultLogger) Warningf(format string, v ...interface{}) {
	l.Printf(Lwarning, format, v...)
}

func (l *defaultLogger) Error(v ...interface{}) {
	l.Print(Lerror, v...)
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	l.Printf(Lerror, format, v...)
}

// Panic logs with Lpanic level and calls panic().
func (l *defaultLogger) Panic(v ...interface{}) {
	var msg string
	if l.level <= Lpanic {
		msg, _ = l.write(Lpanic, v...)
	}
	panic(msg)
}

// Panicf logs with Lpanic level and calls panic().
func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	var msg string
	if l.level <= Lpanic {
		msg, _ = l.writef(Lpanic, format, v...)
	}
	panic(msg)
}

// Fatal logs with Lfatal level and calls os.Exit(1).
func (l *defaultLogger) Fatal(v ...interface{}) {
	l.Print(Lfatal, v...)
	os.Exit(1)
}

// Fatalf logs with Lfatal level and calls os.Exit(1).
func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	l.Printf(Lfatal, format, v...)
	os.Exit(1)
}

// Level returns current minimal logging level.
func (l *defaultLogger) Level() int {
	return l.level
}

// SetLevel sets minimal output level. All logging with levels below lvl will
// be silenced.
func (l *defaultLogger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *defaultLogger) write(lvl int, v ...interface{}) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msg := fmt.Sprintln(v...)
	return msg, l.f.Format(l.w, lvl, l.name, msg)
}

func (l *defaultLogger) writef(lvl int, format string, v ...interface{}) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msg := fmt.Sprintf(format, v...)
	return msg, l.f.Format(l.w, lvl, l.name, msg)
}

func (l *defaultLogger) writeStack(lvl int) {
	stack := make([]byte, 1024*8)
	stack = stack[:runtime.Stack(stack, false)]
	l.w.Write(stack)
}

var Default Logger

func Print(lvl int, v ...interface{}) {
	Default.Print(lvl, v...)
}

func Printf(lvl int, format string, v ...interface{}) {
	Default.Printf(lvl, format, v...)
}

func Debug(v ...interface{}) {
	Default.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	Default.Debugf(format, v...)
}

func Info(v ...interface{}) {
	Default.Debug(v...)
}

func Infof(format string, v ...interface{}) {
	Default.Debugf(format, v...)
}

func Warning(v ...interface{}) {
	Default.Debug(v...)
}

func Warningf(format string, v ...interface{}) {
	Default.Debugf(format, v...)
}

func Error(v ...interface{}) {
	Default.Debug(v...)
}

func Errorf(format string, v ...interface{}) {
	Default.Debugf(format, v...)
}

func Panic(v ...interface{}) {
	Default.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	Default.Panicf(format, v...)
}

func Fatal(v ...interface{}) {
	Default.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	Default.Fatalf(format, v...)
}

func Level() int {
	return Default.Level()
}

func SetLevel(lvl int) {
	Default.SetLevel(lvl)
}

// Default formatter is default output formatter.
// It produces colorized output if output destination is a tty and plain text otherwise.
// Colorization additionally can be controller by CLICOLOR and CLICOLOR_FORCE environment
// variables: CLICOLOR=0 disables colorization and CLICOLOR_FORCE=1 forces it even when
// output is not a tty.
var DefaultFormatter Formatter

func init() {
	isAtty := isatty.IsTerminal(os.Stdout.Fd())
	clicolor := os.Getenv("CLICOLOR")
	clicolorForce := os.Getenv("CLICOLOR_FORCE")

	var writer io.Writer
	if (isAtty || clicolorForce == "1") && clicolor != "0" {
		writer = colorable.NewColorableStdout()
		DefaultFormatter = NewTextFormatter(Fdefault | Fcolor)
	} else {
		writer = os.Stdout
		DefaultFormatter = NewTextFormatter(Fdefault)
	}

	Default = NewLogger(writer, "", DefaultFormatter)
}
