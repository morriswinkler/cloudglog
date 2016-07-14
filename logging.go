package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const CallDepth = 3

var (
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger

	LogLevel int
)

func initLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

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

}

func init() {
	initLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	// get LogLevel from env
	getLogLevel := os.Getenv("LOG_LEVEL")
	if len(getLogLevel) == 0 {
		// loglevel is not in env, set to default
		LogLevel = 0
	} else {
		// loglevel is in env convert to int
		var err error
		LogLevel, err = strconv.Atoi(getLogLevel)
		if err != nil {
			// sorry there was an error, fallback to default level
			Warning("reading loglevel from envieronment variable, falling back to default level 0")
			LogLevel = 0
		}

	}

}

// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	info.Output(CallDepth, fmt.Sprint(args...))
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Infoln(args ...interface{}) {
	info.Output(CallDepth, fmt.Sprintln(args...))
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {
	info.Output(CallDepth, fmt.Sprintf(format, args...))
}

// Warning logs to the Warning log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	warning.Output(CallDepth, fmt.Sprint(args...))

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

// Errorln logs to the Error log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Errorln(args ...interface{}) {
	err.Output(CallDepth, fmt.Sprintln(args...))
}

// Errorf logs to the Error log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {
	err.Output(CallDepth, fmt.Sprintf(format, args...))
}

// Panic logs to the Error log followed by a call to panic.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Panic(args ...interface{}) {
	log.Output(CallDepth, fmt.Sprint(args...))
}

// Panicln logs to the Error log followed by a call to panic.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Panicln(args ...interface{}) {
	log.Output(CallDepth, fmt.Sprintln(args...))
}

// Panic	f logs to the Error log followed by a call to panic.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Panicf(format string, args ...interface{}) {
	log.Output(CallDepth, fmt.Sprintf(format, args...))
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
