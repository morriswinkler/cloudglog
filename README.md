# cloudglog
--
    import "github.com/morriswinkler/cloudglog"

Package cloudglog is a logger that outputs to stdout. It is strongly based on
glog but without any kind of buffering.


### LogFile

By default logging goes to stdout, use LogFile(file)

Example:

    f, err := os.open("filename")
    if err != nil {
        cloudglog.Fatal(err)
    }

    w := bufio.NewWriter(f)
    defer w.Flush()

    cloudglog.LogFile(w)


Format styles

define the log output format, use FormatStyle(style) to set one of:

    DefaultFormat		: the original glog format
    ModernFormat		: shorter format, uses brackets to separate Package, File, Line

Example:

    cloudglog.FormatStyle(cloudglog.ModernFormat)


### Color Styles

define coloring schemes, use ColorStyle(style) to set one of:

    NoColor                  	: no colors
    PrefixColor              	: colorize from prefix until line number
    PrefixBoldColor          	: colorize from prefix until line number with bold colors
    FullColor                	: colorize everything
    FullBoldColor            	: colorize everything with bold colors
    FullColorWithBoldMessage 	: colorize everything with bold colored message
    FullColorWithBoldPrefix  	: colorize everything with bold coloring from prefix until line number

Example:

    cloudglog.ColorStyle(cloudglog.FullColor)


### LogFilter

can be used to filter logging of other packages that provide a way to set the
log output. It takes a io.Writer as output and a logType and returns a
io.Writer.

TODO: make this function more idiomatic

Example:

     ERROR = log.New(cloudglog.LogFilter(os.Stdout, cloudglog.ERROR),
    	"ERROR: ",
    	log.Ldate|log.Ltime|log.Llongfile)

## Usage

```go
const (
	// formatStyle controlss the output format
	DefaultFormat formatStyle = iota // PREFIX: YYYY/MM/DD HH:MM:SS log.Llongfile Message
	ModernFormat                     // PREFIX: YYYY/MM/DD HH:MM:SS [Package][File][:Line] Message
)
```

```go
const (
	// logType, only used in LogFilter to set the right color
	// TODO: while reworking LogFilter overhink this too
	TRACE   logType = iota // TRACE: ColorCyan
	INFO                   // INFO: ColorGreen
	WARNING                // WARNING: ColorYellow
	ERROR                  // ERROR: ColorRed
	FATAL                  // FATAL: ColorMagenta
)
```

```go
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
```

```go
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
```

```go
const CallDepth = 2 // depth to trace the caller file

```

```go
var LogLevel int // logging level for V() type calls, can also be set by LOG_LEVEL environment variable

```

#### func  ColorsStyle

```go
func ColorsStyle(cStyle colorStyle)
```
SetColor defines the coloring format

#### func  Error

```go
func Error(args ...interface{})
```
Error logs to the ERROR log. Arguments are handled in the manner of fmt.Print; a
newline is appended if missing.

#### func  ErrorDepth

```go
func ErrorDepth(depth int, args ...interface{})
```
ErrorDepth acts as ERROR but uses depth to determine which call frame to log.
ErrorDepth(0, "msg") is the same as Error("msg").

#### func  Errorf

```go
func Errorf(format string, args ...interface{})
```
Errorf logs to the ERROR log. Arguments are handled in the manner of fmt.Printf;
a newline is appended if missing.

#### func  Errorln

```go
func Errorln(args ...interface{})
```
Errorln logs to the ERROR log. Arguments are handled in the manner of
fmt.Println; a newline is appended if missing.

#### func  Exit

```go
func Exit(args ...interface{})
```
Exit logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
Arguments are handled in the manner of fmt.Print; a newline is appended if
missing.

#### func  ExitDepth

```go
func ExitDepth(depth int, args ...interface{})
```
ExitDepth acts as Exit but uses depth to determine which call frame to log.
ExitDepth(0, "msg") is the same as Exit("msg").

#### func  Exitf

```go
func Exitf(format string, args ...interface{})
```
Exitf logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).
Arguments are handled in the manner of fmt.Printf; a newline is appended if
missing.

#### func  Exitln

```go
func Exitln(args ...interface{})
```
Exitln logs to the FATAL, ERROR, WARNING, and INFO logs, then calls os.Exit(1).

#### func  Fatal

```go
func Fatal(args ...interface{})
```
Fatal logs to the FATAL log Arguments are handled in the manner of fmt.Print; a
newline is appended if missing.

#### func  FatalDepth

```go
func FatalDepth(depth int, args ...interface{})
```
FatalDepth acts as FATAL but uses depth to determine which call frame to log.
FatalDepth(0, "msg") is the same as Fatal("msg").

#### func  Fatalf

```go
func Fatalf(format string, args ...interface{})
```
Fatalf logs to the FATAL log. Arguments are handled in the manner of fmt.Printf;
a newline is appended if missing.

#### func  Fatalln

```go
func Fatalln(args ...interface{})
```
Fatalln logs to the FATAL log.

#### func  FormatStyle

```go
func FormatStyle(f formatStyle)
```
FormatStyle changes the formatStyle

#### func  Info

```go
func Info(args ...interface{})
```
Info logs to the INFO log. Arguments are handled in the manner of fmt.Print; a
newline is appended if missing.

#### func  InfoDepth

```go
func InfoDepth(depth int, args ...interface{})
```
InfoDepth acts as Info but uses depth to determine which call frame to log.
InfoDepth(0, "msg") is the same as Info("msg").

#### func  Infof

```go
func Infof(format string, args ...interface{})
```
Infof logs to the INFO log. Arguments are handled in the manner of fmt.Printf; a
newline is appended if missing.

#### func  Infoln

```go
func Infoln(args ...interface{})
```
Infoln logs to the INFO log. Arguments are handled in the manner of fmt.Println;
a newline is appended if missing.

#### func  LogFile

```go
func LogFile(file io.Writer)
```
LogFile sets the logfile to write to

#### func  LogFileName

```go
func LogFileName()
```
LogFileName will log only file names

#### func  LogFilePath

```go
func LogFilePath()
```
LogFilePath will log path and file name, this is the default

#### func  LogFilter

```go
func LogFilter(out io.Writer, l logType) io.Writer
```
LogFilter can be used to filter logging of other packages that provide a way to
set the log output. It takes a io.Writer as output and a logType and returns a
io.Writer.

#### func  Warning

```go
func Warning(args ...interface{})
```
Warning logs to the WARNING log. Arguments are handled in the manner of
fmt.Print; a newline is appended if missing.

#### func  WarningDepth

```go
func WarningDepth(depth int, args ...interface{})
```
WarningDepth acts as WARNING but uses depth to determine which call frame to
log. WarningDepth(0, "msg") is the same as Warning("msg").

#### func  Warningf

```go
func Warningf(format string, args ...interface{})
```
Warningf logs to the WARNING log. Arguments are handled in the manner of
fmt.Printf; a newline is appended if missing.

#### func  Warningln

```go
func Warningln(args ...interface{})
```
Warningln logs to the WARNING log. Arguments are handled in the manner of
fmt.Println; a newline is appended if missing.

#### type Verbosity

```go
type Verbosity bool
```

Verbosity is a boolean type that implements Infof (like Printf) etc. See the
documentation of V for more information.

#### func  V

```go
func V(level int) Verbosity
```
V reports whether verbosity at the call site is at least the requested level.
The returned value is a boolean of type Verbosity, which implements Info, Infoln
and Infof. These methods will write to the Info log if called. Thus, one may
write either

    if cloudglog.V(2) { cloudglog.Info("log this") }

or

    cloudglog.V(2).Info("log this")

The second form is shorter but the first is cheaper if logging is off because it
does not evaluate its arguments.

#### func (Verbosity) Error

```go
func (v Verbosity) Error(args ...interface{})
```
Error is equivalent to the global Error function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) ErrorDepth

```go
func (v Verbosity) ErrorDepth(depth int, args ...interface{})
```
ErrorDepth is equivalent to the global ErrorDepth function, guarded by the value
of v. See the documentation of V for usage.

#### func (Verbosity) Errorf

```go
func (v Verbosity) Errorf(format string, args ...interface{})
```
Errorf is equivalent to the global Errorf function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) Errorln

```go
func (v Verbosity) Errorln(args ...interface{})
```
Errorln is equivalent to the global Errorln function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) Exit

```go
func (v Verbosity) Exit(args ...interface{})
```
Exit is equivalent to the global Exit function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) ExitDepth

```go
func (v Verbosity) ExitDepth(depth int, args ...interface{})
```
c is equivalent to the global Exitln function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) Exitf

```go
func (v Verbosity) Exitf(format string, args ...interface{})
```
Exitf is equivalent to the global Exitf function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) Exitln

```go
func (v Verbosity) Exitln(args ...interface{})
```
Exitln is equivalent to the global Exitln function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) Fatal

```go
func (v Verbosity) Fatal(args ...interface{})
```
Fatal is equivalent to the global Fatal function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) FatalDepth

```go
func (v Verbosity) FatalDepth(depth int, args ...interface{})
```
FatalDepth is equivalent to the global FatalDepth function, guarded by the value
of v. See the documentation of V for usage.

#### func (Verbosity) Fatalf

```go
func (v Verbosity) Fatalf(format string, args ...interface{})
```
Fatalf is equivalent to the global Fatalf function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) Fatalln

```go
func (v Verbosity) Fatalln(args ...interface{})
```
Fatalln is equivalent to the global Fatalln function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) Info

```go
func (v Verbosity) Info(args ...interface{})
```
Info is equivalent to the global Info function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) InfoDepth

```go
func (v Verbosity) InfoDepth(depth int, args ...interface{})
```
InfoDepth is equivalent to the global InfoDepth function, guarded by the value
of v. See the documentation of V for usage.

#### func (Verbosity) Infof

```go
func (v Verbosity) Infof(format string, args ...interface{})
```
Infof is equivalent to the global Infof function, guarded by the value of v. See
the documentation of V for usage.

#### func (Verbosity) Infoln

```go
func (v Verbosity) Infoln(args ...interface{})
```
Infoln is equivalent to the global Infoln function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) Warning

```go
func (v Verbosity) Warning(args ...interface{})
```
Warning is equivalent to the global Warning function, guarded by the value of v.
See the documentation of V for usage.

#### func (Verbosity) WarningDepth

```go
func (v Verbosity) WarningDepth(depth int, args ...interface{})
```
WarningDepth is equivalent to the global WarningDepth function, guarded by the
value of v. See the documentation of V for usage.

#### func (Verbosity) Warningf

```go
func (v Verbosity) Warningf(format string, args ...interface{})
```
Warningf is equivalent to the global Warningf function, guarded by the value of
v. See the documentation of V for usage.

#### func (Verbosity) Warningln

```go
func (v Verbosity) Warningln(args ...interface{})
```
Warningln is equivalent to the global Warningln function, guarded by the value
of v. See the documentation of V for usage.
