package cloudglog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const CallDepth = 2

var (
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	fatal   *log.Logger

	LogLevel int
)

func initLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	fatalHandle io.Writer) {

	trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Llongfile)

	info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Llongfile)

	warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Llongfile)

	err = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Llongfile)
	fatal = log.New(fatalHandle,
		"Fatal: ",
		log.Ldate|log.Ltime|log.Llongfile)

}

func init() {
	initLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	// get LogLevel from env
	getLogLevel := os.Getenv("LOG_LEVEL")
	if len(getLogLevel) == 0 {
		// loglevel is not in env, set to default
		Info("missing LOG_LEVEL envieronment variable, falling back to default level 0")
		LogLevel = 0
	} else {
		// loglevel is a string, convert it to int
		var err error
		LogLevel, err = strconv.Atoi(getLogLevel)
		if err != nil {
			// sorry there was an error, fallback to default level
			Error("reading loglevel from envieronment variable, falling back to default level 0")
			LogLevel = 0
		}

	}

}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	info.Output(CallDepth, fmt.Sprint(args...))
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func InfoDepth(depth int, args ...interface{}) {
	info.Output(depth, fmt.Sprint(args...))
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Infoln(args ...interface{}) {
	info.Output(CallDepth, fmt.Sprintln(args...))
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	info.Output(CallDepth, buf.String())
}

// Warning logs to the Warning log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	warning.Output(CallDepth, fmt.Sprint(args...))
}

// WarningDepth acts as Warning but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warning("msg").
func WarningDepth(depth int, args ...interface{}) {
	warning.Output(depth, fmt.Sprint(args...))
}

// Warningln logs to the Warning log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Warningln(args ...interface{}) {
	warning.Output(CallDepth, fmt.Sprintln(args...))
}

// Warningf logs to the Warning log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	warning.Output(CallDepth, fmt.Sprintf(format, args...))
}

// Error logs to the Error log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	err.Output(CallDepth, fmt.Sprint(args...))
}

// ErrorDepth acts as Error but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func ErrorDepth(depth int, args ...interface{}) {
	err.Output(depth, fmt.Sprint(args...))
}

// Errorln logs to the Error log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Errorln(args ...interface{}) {
	err.Output(CallDepth, fmt.Sprintln(args...))
}

// Errorf logs to the Error log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	err.Output(CallDepth, buf.String())
}

// Fatal logs to the FATAL log
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Fatal(args ...interface{}) {
	fatal.Output(CallDepth, fmt.Sprint(args...))
}

// FatalDepth acts as Fatal but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func FatalDepth(depth int, args ...interface{}) {
	fatal.Output(depth, fmt.Sprint(args...))
}

// Fatalln logs to the log.
func Fatalln(args ...interface{}) {
	fatal.Output(CallDepth, fmt.Sprintln(args...))
}

// Fatalf logs to the log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Fatalf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	fatal.Output(CallDepth, buf.String())
}

// Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Exit(args ...interface{}) {
	fatal.Output(CallDepth, fmt.Sprint(args...))
	os.Exit(1)
}

// ExitDepth acts as Exit but uses depth to determine which call frame to log.
// ExitDepth(0, "msg") is the same as Exit("msg").
func ExitDepth(depth int, args ...interface{}) {
	fatal.Output(depth, fmt.Sprint(args...))
	os.Exit(1)
}

// Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
func Exitln(args ...interface{}) {
	fatal.Output(CallDepth, fmt.Sprintln(args...))
	os.Exit(1)
}

// Exitf logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Exitf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	fatal.Output(CallDepth, buf.String())
	os.Exit(1)
}

type Verbose bool

// V reports whether verbosity at the call site is at least the requested level.
// The returned value is a boolean of type Verbose, which implements Info, Infoln
// and Infof. These methods will write to the Info log if called.
// Thus, one may write either
//	if glog.V(2) { glog.Info("log this") }
// or
//	glog.V(2).Info("log this")
// The second form is shorter but the first is cheaper if logging is off because it does
// not evaluate its arguments.

func V(level int) Verbose {
	// This function tries hard to be cheap unless there's work to do.
	// The fast path is two atomic loads and compares.

	// Here is a cheap but safe test to see if V logging is enabled globally.
	if LogLevel >= level {
		return Verbose(true)
	}

	return Verbose(false)
}

// Info is equivalent to the global Info function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Info(args ...interface{}) {
	if v {
		Info(args...)
	}
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infoln(args ...interface{}) {
	if v {
		Infoln(args...)
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infof(format string, args ...interface{}) {
	if v {
		Infof(format, args...)
	}
}
