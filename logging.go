// Package cloudglog is a logger that outputs to stdout. It is strongly based on glog but
// without any kind of buffering.
//
// LogFile
//
// By default logging goes to stdout, use LogFile(file)
//
// Example:
//    f, err := os.open("filename")
//    if err != nil {
//        cloudglog.Fatal(err)
//    }
//
//    w := bufio.NewWriter(f)
//    defer w.Flush()
//
//    cloudglog.LogFile(w)
//
//
// Format Styles
//
// define the log output format, use FormatStyle(style) to set one of:
//
//  DefaultFormat		: the original glog format
//  ModernFormat		: shorter format, uses brackets to separate Package, File, Line
//
// Example:
//  cloudglog.FormatStyle(cloudglog.ModernFormat)
//
// Color Styles
//
// define coloring schemes, use ColorStyle(style) to set one of:
//
//  NoColor                  	: no colors
//  PrefixColor              	: colorize from prefix until line number
//  PrefixBoldColor          	: colorize from prefix until line number with bold colors
//  FullColor                	: colorize everything
//  FullBoldColor            	: colorize everything with bold colors
//  FullColorWithBoldMessage 	: colorize everything with bold colored message
//  FullColorWithBoldPrefix  	: colorize everything with bold coloring from prefix until line number
//
// Example:
//  cloudglog.ColorStyle(cloudglog.FullColor)
//
// LogFilter
//
// can be used to filter logging of other packages
// that provide a way to set the log output. It takes a io.Writer
// as output and a logType and returns a io.Writer.
//
// TODO: make this function more idiomatic
//
// Example:
//   ERROR = log.New(cloudglog.LogFilter(os.Stdout, cloudglog.ERROR),
//  	"ERROR: ",
//  	log.Ldate|log.Ltime|log.Llongfile)
//
package cloudglog

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"runtime"
)

const CallDepth = 2 // depth to trace the caller file

var LogLevel int // logging level for V() type calls, can also be set by LOG_LEVEL environment variable

type formatStyle int

const (
	// formatStyle controlss the output format
	DefaultFormat formatStyle = iota // PREFIX: YYYY/MM/DD HH:MM:SS log.Llongfile Message
	ModernFormat                     // PREFIX: YYYY/MM/DD HH:MM:SS [Package][File][:Line] Message
)

var currentFormat formatStyle

// FormatStyle changes the formatStyle
func FormatStyle(f formatStyle) {
	currentFormat = f
	setupLogger(logFile, logFile, logFile, logFile, logFile)
}

var logFile io.Writer = os.Stdout

// LogFile sets the logfile to write to
func LogFile(file io.Writer) {
	logFile = file
	setupLogger(logFile, logFile, logFile, logFile, logFile)
}

// can be set to log.Llongfile or log.Lshortfile
var lFileLength = log.Llongfile


// LogFileName will log only file names
func LogFileName() {
	lFileLength = log.Lshortfile
}

// LogFilePath will log path and file name, this is the default
func LogFilePath() {
	lFileLength = log.Llongfile
}


type logType int

const (
	// logType, only used in LogFilter to set the right color
	// TODO: while reworking LogFilter overhink this too
	TRACE   logType = iota // TRACE: ColorCyan
	INFO                   // INFO: ColorGreen
	WARNING                // WARNING: ColorYellow
	ERROR                  // ERROR: ColorRed
	FATAL                  // FATAL: ColorMagenta
)

type colorType int

const (
	// colorType used to set term color
	ColorBlack colorType = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

type colorStyle int

const (
	// colorStyle used to set coloring
	NoColor                  colorStyle = iota // no colors
	PrefixColor                                // colorize from prefix until line number
	PrefixBoldColor                            // colorize from prefix until line number with bold colors
	FullColor                                  // colorize everything
	FullBoldColor                              // colorize everything with bold colors
	FullColorWithBoldMessage                   // colorize everything with bold colored message
	FullColorWithBoldPrefix                    // colorize everything with bold coloring from prefix until line number
)

var (
	colorFormating colorStyle = NoColor

	colors = []string{
		TRACE:   colorSeq(ColorCyan),
		INFO:    colorSeq(ColorGreen),
		WARNING: colorSeq(ColorYellow),
		ERROR:   colorSeq(ColorRed),
		FATAL:   colorSeq(ColorMagenta),
	}

	boldcolors = []string{
		TRACE:   colorSeqBold(ColorCyan),
		INFO:    colorSeqBold(ColorGreen),
		WARNING: colorSeqBold(ColorYellow),
		ERROR:   colorSeqBold(ColorRed),
		FATAL:   colorSeqBold(ColorMagenta),
	}
)

func colorSeq(color colorType) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

func colorSeqBold(color colorType) string {
	return fmt.Sprintf("\033[%d;1m", int(color))
}

func addColor(lType logType, prefixEnd int, message []string) []string {

	var col, bcol string

	switch colorFormating {
	case PrefixColor:
		col = colors[lType]
		message[prefixEnd] = strings.Join([]string{message[prefixEnd], "\033[0m"}, "")
		message = append([]string{col}, message...)
	case PrefixBoldColor:
		col = boldcolors[lType]
		message[prefixEnd] = strings.Join([]string{message[prefixEnd], "\033[0m"}, "")
		message = append([]string{col}, message...)
	case FullColor:
		col = colors[lType]
		message = append([]string{col}, message...)
		message = append(message, "\033[0m")
	case FullBoldColor:
		col = boldcolors[lType]
		message = append([]string{col}, message...)
		message = append(message, "\033[0m")
	case FullColorWithBoldMessage:
		col = colors[lType]
		bcol = boldcolors[lType]
		message[prefixEnd] = strings.Join([]string{message[prefixEnd], bcol}, "")
		message = append([]string{col}, message...)
		message = append(message, "\033[0m")
	case FullColorWithBoldPrefix:
		col = colors[lType]
		bcol = boldcolors[lType]
		message[prefixEnd] = strings.Join([]string{message[prefixEnd], col}, "")
		message = append([]string{bcol}, message...)
		message = append(message, "\033[0m")
	}

	return message
}

// SetColor defines the coloring format
func ColorsStyle(cStyle colorStyle) {
	colorFormating = cStyle
}

var (
	traceLog   *log.Logger
	infoLog    *log.Logger
	warningLog *log.Logger
	errorLog   *log.Logger
	fatalLog   *log.Logger
)

func setupLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	fatalHandle io.Writer) {

	traceLog = log.New(LogFilter(traceHandle, TRACE),
		"TRACE: ",
		log.Ldate|log.Ltime|lFileLength)

	infoLog = log.New(LogFilter(infoHandle, INFO),
		"INFO: ",
		log.Ldate|log.Ltime|lFileLength)

	warningLog = log.New(LogFilter(warningHandle, WARNING),
		"WARNING: ",
		log.Ldate|log.Ltime|lFileLength)

	errorLog = log.New(LogFilter(errorHandle, ERROR),
		"ERROR: ",
		log.Ldate|log.Ltime|lFileLength)
	fatalLog = log.New(LogFilter(fatalHandle, FATAL),
		"Fatal: ",
		log.Ldate|log.Ltime|lFileLength)
}

// LogFilter can be used to filter logging of other packages
// that provide a way to set the log output. It takes a io.Writer
// as output and a logType and returns a io.Writer.
func LogFilter(out io.Writer, l logType) io.Writer {

	// TODO: make this more idiomatic
	switch currentFormat {
	case DefaultFormat:
		return &defaultLogger{out: out, logType: l}
	case ModernFormat:
		return &modernLogger{out: out, logType: l}
	}

	return ioutil.Discard
}

type defaultLogger struct {
	out     io.Writer
	logType logType
}

func (d *defaultLogger) Write(bytes []byte) (int, error) {

	string_ := string(bytes)

	// splitFunc that splits at avery ' '
	stringSplit := func(r rune) bool {
		return r == ' '
	}

	// split to access Llongfile
	format := strings.FieldsFunc(string_, stringSplit)

	// find log prefix (starts with '/' and ends with ':')
	formatByte := []byte(format[3])
	prefixEnd := 2
	if (formatByte[0] == '/') && (formatByte[len(formatByte)-1] == ':') {
		prefixEnd = 3
	}

	// format color
	format = addColor(d.logType, prefixEnd, format)

	// join string
	defaultFormat := strings.Join(format, " ")

	return d.out.Write([]byte(defaultFormat))
}

type modernLogger struct {
	out     io.Writer
	logType logType
}

// TODO: check efficiency and maybe reimplement in []byte operations
func (m *modernLogger) Write(bytes []byte) (int, error) {

	string_ := string(bytes)

	// splitFunc that splits at avery ' '
	stringSplit := func(r rune) bool {
		return r == ' '
	}

	// split to access Llongfile
	format := strings.FieldsFunc(string_, stringSplit)

	// find log prefix (starts with '/' and ends with ':')
	formatByte := []byte(format[3])
	prefixEnd := 2
	if (formatByte[0] == '/') && (formatByte[len(formatByte)-1] == ':') {
		prefixEnd = 3
	}

	// splitFunc for log.Llongfile
	longFileSplit := func(r rune) bool {
		return r == '/' || r == ':'
	}
	var modernLongFile = make([]string, 3)
	// split log.Llongfile
	subFormat := strings.FieldsFunc(format[prefixEnd], longFileSplit)
	copy(modernLongFile, subFormat[len(subFormat)-3:]) // package, file, line

	// add []'s and a trailing tab
	format[prefixEnd] = strings.Join([]string{"[", modernLongFile[0], "]", "[", modernLongFile[1], "]", "[:", modernLongFile[2], "]", "\t"}, "")

	// format color
	format = addColor(m.logType, prefixEnd, format)

	// join string
	modernFormat := strings.Join(format, " ")

	return m.out.Write([]byte(modernFormat))
}


// stacks is a wrapper for runtime.Stack that attempts to recover the data for all goroutines.
// Todo: wire this func into fatal and panic
func stacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit. Start large, though.
	n := 10000
	if all {
		n = 100000
	}
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, all)
		if nbytes < len(trace) {
			return trace[:nbytes]
		}
		n *= 2
	}
	return trace
}

func init() {

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

	currentFormat = DefaultFormat // init with DefaultFormat
	setupLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
}


// Info logs to the INFO log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Info(args ...interface{}) {
	infoLog.Output(CallDepth, fmt.Sprint(args...))
}

// InfoDepth acts as Info but uses depth to determine which call frame to log.
// InfoDepth(0, "msg") is the same as Info("msg").
func InfoDepth(depth int, args ...interface{}) {
	infoLog.Output(depth, fmt.Sprint(args...))
}

// Infoln logs to the INFO log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Infoln(args ...interface{}) {
	infoLog.Output(CallDepth, fmt.Sprintln(args...))
}

// Infof logs to the INFO log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Infof(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	infoLog.Output(CallDepth, buf.String())
}

// Warning logs to the WARNING log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Warning(args ...interface{}) {
	warningLog.Output(CallDepth, fmt.Sprint(args...))
}

// WarningDepth acts as WARNING but uses depth to determine which call frame to log.
// WarningDepth(0, "msg") is the same as Warning("msg").
func WarningDepth(depth int, args ...interface{}) {
	warningLog.Output(depth, fmt.Sprint(args...))
}

// Warningln logs to the WARNING log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Warningln(args ...interface{}) {
	warningLog.Output(CallDepth, fmt.Sprintln(args...))
}

// Warningf logs to the WARNING log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Warningf(format string, args ...interface{}) {
	warningLog.Output(CallDepth, fmt.Sprintf(format, args...))
}

// Error logs to the ERROR log.
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Error(args ...interface{}) {
	errorLog.Output(CallDepth, fmt.Sprint(args...))
}

// ErrorDepth acts as ERROR but uses depth to determine which call frame to log.
// ErrorDepth(0, "msg") is the same as Error("msg").
func ErrorDepth(depth int, args ...interface{}) {
	errorLog.Output(depth, fmt.Sprint(args...))
}

// Errorln logs to the ERROR log.
// Arguments are handled in the manner of fmt.Println; a newline is appended if missing.
func Errorln(args ...interface{}) {
	errorLog.Output(CallDepth, fmt.Sprintln(args...))
}

// Errorf logs to the ERROR log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Errorf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	errorLog.Output(CallDepth, buf.String())
}

// Fatal logs to the FATAL log
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Fatal(args ...interface{}) {
	fatalLog.Output(CallDepth, fmt.Sprint(args...))
	// Todo: check if we need to flush here.
	os.Exit(1)
}

// FatalDepth acts as FATAL but uses depth to determine which call frame to log.
// FatalDepth(0, "msg") is the same as Fatal("msg").
func FatalDepth(depth int, args ...interface{}) {
	fatalLog.Output(depth, fmt.Sprint(args...))
	os.Exit(1)
}

// Fatalln logs to the FATAL log.
func Fatalln(args ...interface{}) {
	fatalLog.Output(CallDepth, fmt.Sprintln(args...))
	os.Exit(1)
}

// Fatalf logs to the FATAL log.
// Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func Fatalf(format string, args ...interface{}) {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	fatalLog.Output(CallDepth, buf.String())
	os.Exit(1)
}

// Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print; a newline is appended if missing.
func Exit(args ...interface{}) {
	fatalLog.Output(CallDepth, fmt.Sprint(args...))
	os.Exit(1)
}

// ExitDepth acts as Exit but uses depth to determine which call frame to log.
// ExitDepth(0, "msg") is the same as Exit("msg").
func ExitDepth(depth int, args ...interface{}) {
	fatalLog.Output(depth, fmt.Sprint(args...))
	os.Exit(1)
}

// Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
func Exitln(args ...interface{}) {
	fatalLog.Output(CallDepth, fmt.Sprintln(args...))
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
	fatalLog.Output(CallDepth, buf.String())
	os.Exit(1)
}

// Verbosity is a boolean type that implements Infof (like Printf) etc.
// See the documentation of V for more information.
type Verbosity bool

// V reports whether verbosity at the call site is at least the requested level.
// The returned value is a boolean of type Verbosity, which implements Info, Infoln
// and Infof. These methods will write to the Info log if called.
// Thus, one may write either
//	if cloudglog.V(2) { cloudglog.Info("log this") }
// or
//	cloudglog.V(2).Info("log this")
// The second form is shorter but the first is cheaper if logging is off because it does
// not evaluate its arguments.
func V(level int) Verbosity {
	// This function tries hard to be cheap unless there's work to do.
	// The fast path is two atomic loads and compares.

	// Here is a cheap but safe test to see if V logging is enabled globally.
	if LogLevel >= level {
		return Verbosity(true)
	}

	return Verbosity(false)
}

// Info is equivalent to the global Info function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Info(args ...interface{}) {
	if v {
		infoLog.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// InfoDepth is equivalent to the global InfoDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) InfoDepth(depth int, args ...interface{}) {
	if v {
		infoLog.Output(depth, fmt.Sprint(args...))
	}
}

// Infoln is equivalent to the global Infoln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Infoln(args ...interface{}) {
	if v {
		infoLog.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Infof is equivalent to the global Infof function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Infof(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		infoLog.Output(CallDepth+1, buf.String())
	}
}

// Warning is equivalent to the global Warning function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Warning(args ...interface{}) {
	if v {
		warningLog.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// WarningDepth is equivalent to the global WarningDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) WarningDepth(depth int, args ...interface{}) {
	if v {
		warningLog.Output(depth, fmt.Sprint(args...))
	}
}

// Warningln is equivalent to the global Warningln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Warningln(args ...interface{}) {
	if v {
		warningLog.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Warningf is equivalent to the global Warningf function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Warningf(format string, args ...interface{}) {
	if v {
		warningLog.Output(CallDepth+1, fmt.Sprintf(format, args...))
	}
}

// Error is equivalent to the global Error function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Error(args ...interface{}) {
	if v {
		errorLog.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// ErrorDepth is equivalent to the global ErrorDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) ErrorDepth(depth int, args ...interface{}) {
	if v {
		errorLog.Output(depth, fmt.Sprint(args...))
	}
}

// Errorln is equivalent to the global Errorln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Errorln(args ...interface{}) {
	if v {
		errorLog.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Errorf is equivalent to the global Errorf function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Errorf(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		errorLog.Output(CallDepth+1, buf.String())
	}
}

// Fatal is equivalent to the global Fatal function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Fatal(args ...interface{}) {
	if v {
		fatalLog.Output(CallDepth+1, fmt.Sprint(args...))
	}
}

// FatalDepth is equivalent to the global FatalDepth function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) FatalDepth(depth int, args ...interface{}) {
	if v {
		fatalLog.Output(depth, fmt.Sprint(args...))
	}
}

// Fatalln is equivalent to the global Fatalln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Fatalln(args ...interface{}) {
	if v {
		fatalLog.Output(CallDepth+1, fmt.Sprintln(args...))
	}
}

// Fatalf is equivalent to the global Fatalf function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Fatalf(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		fatalLog.Output(CallDepth+1, buf.String())
	}
}

// Exit is equivalent to the global Exit function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Exit(args ...interface{}) {
	if v {
		fatalLog.Output(CallDepth+1, fmt.Sprint(args...))
		os.Exit(1)
	}
}

// c is equivalent to the global Exitln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) ExitDepth(depth int, args ...interface{}) {
	if v {
		fatalLog.Output(depth, fmt.Sprint(args...))
		os.Exit(1)
	}
}

// Exitln is equivalent to the global Exitln function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Exitln(args ...interface{}) {
	if v {
		fatalLog.Output(CallDepth+1, fmt.Sprintln(args...))
		os.Exit(1)
	}
}

// Exitf is equivalent to the global Exitf  function, guarded by the value of v.
// See the documentation of V for usage.
func (v Verbosity) Exitf(format string, args ...interface{}) {
	if v {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, format, args...)
		if buf.Bytes()[buf.Len()-1] != '\n' {
			buf.WriteByte('\n')
		}
		fatalLog.Output(CallDepth+1, buf.String())
		os.Exit(1)
	}
}
