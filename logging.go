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
	TRACE   *log.Logger
	INFO    *log.Logger
	WARNING *log.Logger
	ERROR     *log.Logger
	FATAL   *log.Logger

	LogLevel int
)

func initLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	fatalHandle io.Writer) {

	TRACE = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Llongfile)

	INFO = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Llongfile)

	WARNING = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Llongfile)

	ERROR = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Llongfile)
	FATAL = log.New(fatalHandle,
		"Fatal: ",
		log.Ldate|log.Ltime|log.Llongfile)

}

func init() {
	initLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	// get LogLevel from env
	getLogLevel := os.Getenv("LOG_LEVEL")
	if len(getLogLevel) == 0 {
		// loglevel is not in env, set to default
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
	INFO.Output(CallDepth, fmt.Sprint(args...))
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func InfoDepth(depth int, args ...interface{}) {
	INFO.Output(depth, fmt.Sprint(args...))
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Infoln(args ...interface{}) {
	INFO.Output(CallDepth, fmt.Sprintln(args...))
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	INFO.Output(CallDepth, buf.String())
}

// Warning logs to the WARNING log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	WARNING.Output(CallDepth, fmt.Sprint(args...))
}

// WarningDepth acts as WARNING but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warning("msg").
func WarningDepth(depth int, args ...interface{}) {
	WARNING.Output(depth, fmt.Sprint(args...))
}

// Warningln logs to the WARNING log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Warningln(args ...interface{}) {
	WARNING.Output(CallDepth, fmt.Sprintln(args...))
}

// Warningf logs to the WARNING log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	WARNING.Output(CallDepth, fmt.Sprintf(format, args...))
}

// Error logs to the ERROR log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	ERROR.Output(CallDepth, fmt.Sprint(args...))
}

// ErrorDepth acts as ERROR but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func ErrorDepth(depth int, args ...interface{}) {
	ERROR.Output(depth, fmt.Sprint(args...))
}

// Errorln logs to the ERROR log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Errorln(args ...interface{}) {
	ERROR.Output(CallDepth, fmt.Sprintln(args...))
}

// Errorf logs to the ERROR log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	ERROR.Output(CallDepth, buf.String())
}

// Fatal logs to the FATAL log
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Fatal(args ...interface{}) {
	FATAL.Output(CallDepth, fmt.Sprint(args...))
}

// FatalDepth acts as FATAL but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func FatalDepth(depth int, args ...interface{}) {
	FATAL.Output(depth, fmt.Sprint(args...))
}

// Fatalln logs to the FATAL log.
func Fatalln(args ...interface{}) {
	FATAL.Output(CallDepth, fmt.Sprintln(args...))
}

// Fatalf logs to the FATAL log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Fatalf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	FATAL.Output(CallDepth, buf.String())
}

// Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Exit(args ...interface{}) {
	FATAL.Output(CallDepth, fmt.Sprint(args...))
	os.Exit(1)
}

// ExitDepth acts as Exit but uses depth to determine which call frame to log.
// ExitDepth(0, "msg") is the same as Exit("msg").
func ExitDepth(depth int, args ...interface{}) {
	FATAL.Output(depth, fmt.Sprint(args...))
	os.Exit(1)
}

// Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
func Exitln(args ...interface{}) {
	FATAL.Output(CallDepth, fmt.Sprintln(args...))
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
	FATAL.Output(CallDepth, buf.String())
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
		INFO.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// InfoDepth is equivalent to the global InfoDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) InfoDepth(depth int, args ...interface{}) {
	if v {
		INFO.Output(depth, fmt.Sprint(args...))
	}
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infoln(args ...interface{}) {
	if v {
		INFO.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Infof(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		INFO.Output(CallDepth+1, buf.String())
	}
}

// Warning is equivalent to the global Warning function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Warning(args ...interface{}) {
	if v {
		WARNING.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// WarningDepth is equivalent to the global WarningDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) WarningDepth(depth int, args ...interface{}) {
	if v {
		WARNING.Output(depth, fmt.Sprint(args...))
	}
}

// Warningln is equivalent to the global Warningln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose)Warningln(args ...interface{}) {
	if v {
		WARNING.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Warningf is equivalent to the global Warningf function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Warningf(format string, args ...interface{}) {
	if v {
		WARNING.Output(CallDepth+1, fmt.Sprintf(format, args...))
	}
}

// Error is equivalent to the global Error function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Error(args ...interface{}) {
	if v {
		ERROR.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// ErrorDepth is equivalent to the global ErrorDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) ErrorDepth(depth int, args ...interface{}) {
	if v {
		ERROR.Output(depth, fmt.Sprint(args...))
	}
}

// Errorln is equivalent to the global Errorln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Errorln(args ...interface{}) {
	if v {
		ERROR.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Errorf is equivalent to the global Errorf function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Errorf(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len() - 1] != '\n' {
			buf.WriteByte('\n')
		}
		ERROR.Output(CallDepth+1, buf.String())
	}
}

// Fatal is equivalent to the global Fatal function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Fatal(args ...interface{}) {
	if v {
		FATAL.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// FatalDepth is equivalent to the global FatalDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) FatalDepth(depth int, args ...interface{}) {
	if v {
		FATAL.Output(depth, fmt.Sprint(args...))
	}
}

// Fatalln is equivalent to the global Fatalln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Fatalln(args ...interface{}) {
	if v {
		FATAL.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Fatalf is equivalent to the global Fatalf function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Fatalf(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len() - 1] != '\n' {
			buf.WriteByte('\n')
		}
		FATAL.Output(CallDepth+1, buf.String())
	}
}

// Exit is equivalent to the global Exit function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Exit(args ...interface{}) {
	if v {
		FATAL.Output(CallDepth+1, fmt.Sprint(args...))
		os.Exit(1)
	}
}

// c is equivalent to the global Exitln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) ExitDepth(depth int, args ...interface{}) {
	if v {
		FATAL.Output(depth, fmt.Sprint(args...))
		os.Exit(1)
	}
}

// Exitln is equivalent to the global Exitln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Exitln(args ...interface{}) {
	if v {
		FATAL.Output(CallDepth+1, fmt.Sprintln(args...))
		os.Exit(1)
	}
}

// Exitf is equivalent to the global Exitf  function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbose) Exitf(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len() - 1] != '\n' {
			buf.WriteByte('\n')
		}
		FATAL.Output(CallDepth+1, buf.String())
		os.Exit(1)
	}
}
