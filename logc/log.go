package logc

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	outLogger zerolog.Logger //nolint:gochecknoglobals // global logger instance
	errLogger zerolog.Logger //nolint:gochecknoglobals // global logger instance
)

func InitializeLogger() {
	dev := flag.Bool("dev", false, "enable console logs for development")
	trace := flag.Bool("trace", false, "sets log level to trace")
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerSkipFrameCount = 2
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		idx := strings.LastIndex(file, "/")
		if idx == -1 {
			idx = strings.LastIndex(file, "\\")
		}

		if idx != -1 {
			file = file[idx+1:]
		}

		return fmt.Sprintf("%s:%d", file, line)
	}

	switch {
	case *trace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case *debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	outWriter, errWriter := io.Writer(os.Stdout), io.Writer(os.Stderr)

	if *dev {
		enableDevelopmentLogging(&outWriter, &errWriter)
	}

	outLogger = zerolog.New(outWriter).With().Timestamp().Caller().Logger()
	errLogger = zerolog.New(errWriter).With().Timestamp().Caller().Logger()
}

func enableDevelopmentLogging(outWriter, errWriter *io.Writer) {
	fmtLevelFunc := func(i any) string {
		return strings.ToUpper(fmt.Sprintf("| %-5s|", i))
	}
	fmtMessageFunc := func(i any) string {
		return fmt.Sprintf("%s", i)
	}
	fmtFieldNameFunc := func(i any) string {
		return fmt.Sprintf("%s=", i)
	}
	fmtFieldValueFunc := func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	*outWriter = zerolog.ConsoleWriter{
		Out:              os.Stdout,
		TimeFormat:       time.RFC3339,
		NoColor:          false,
		FormatLevel:      fmtLevelFunc,
		FormatMessage:    fmtMessageFunc,
		FormatFieldName:  fmtFieldNameFunc,
		FormatFieldValue: fmtFieldValueFunc,
	}
	*errWriter = zerolog.ConsoleWriter{
		Out:              os.Stderr,
		TimeFormat:       time.RFC3339,
		NoColor:          false,
		FormatLevel:      fmtLevelFunc,
		FormatMessage:    fmtMessageFunc,
		FormatFieldName:  fmtFieldNameFunc,
		FormatFieldValue: fmtFieldValueFunc,
	}
}

//nolint:zerologlint // returning zerolog Event for chaining
func OutLog() *zerolog.Event {
	return outLogger.Log()
}

//nolint:zerologlint // returning zerolog Event for chaining
func ErrLog() *zerolog.Event {
	return errLogger.Log()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Trace() *zerolog.Event {
	return outLogger.Trace()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Debug() *zerolog.Event {
	return outLogger.Debug()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Info() *zerolog.Event {
	return outLogger.Info()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Warn() *zerolog.Event {
	return outLogger.Warn()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Error() *zerolog.Event {
	return outLogger.Error().Stack()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Fatal() *zerolog.Event {
	return errLogger.Fatal().Stack()
}

//nolint:zerologlint // returning zerolog Event for chaining
func Panic() *zerolog.Event {
	return errLogger.Panic().Stack()
}
