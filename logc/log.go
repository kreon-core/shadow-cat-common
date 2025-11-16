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
	zerolog.CallerSkipFrameCount = 3
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
		return fmt.Sprintf("%s:", i)
	}
	fmtFieldValueFunc := func(i any) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
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
		TimeFormat:       "2006-01-02 15:04:05",
		NoColor:          false,
		FormatLevel:      fmtLevelFunc,
		FormatMessage:    fmtMessageFunc,
		FormatFieldName:  fmtFieldNameFunc,
		FormatFieldValue: fmtFieldValueFunc,
	}
}

// Helper to add fields to zerolog Event.

func addFields(e *zerolog.Event, fields ...any) *zerolog.Event {
	i, n := 0, len(fields)
	for i < n {
		key, ok := fields[i].(string)
		if ok {
			if i+1 < n {
				e = addKeyValueField(e, key, fields[i+1])
				i += 2

				continue
			}
		} else {
			switch v := fields[i].(type) {
			case error:
				e.Err(v)
			default:
				e.Any("field", v)
			}
		}

		i++
	}

	return e
}

func addKeyValueField(e *zerolog.Event, key string, value any) *zerolog.Event {
	switch v := value.(type) {
	case string:
		e.Str(key, v)
	case []string:
		e.Strs(key, v)
	case int:
		e.Int(key, v)
	case []int:
		e.Ints(key, v)
	case int64:
		e.Int64(key, v)
	case []int64:
		e.Ints64(key, v)
	case float64:
		e.Float64(key, v)
	case []float64:
		e.Floats64(key, v)
	case bool:
		e.Bool(key, v)
	case []bool:
		e.Bools(key, v)
	case time.Duration:
		e.Dur(key, v)
	case []time.Duration:
		e.Durs(key, v)
	default:
		e.Interface(key, v)
	}

	return e
}

// --- Logging functions ---

func Trace(msg string, fields ...any) {
	addFields(outLogger.Trace(), fields...).Msg(msg)
}

func Debug(msg string, fields ...any) {
	addFields(outLogger.Debug(), fields...).Msg(msg)
}

func Info(msg string, fields ...any) {
	addFields(outLogger.Info(), fields...).Msg(msg)
}

func Warn(msg string, fields ...any) {
	addFields(errLogger.Warn(), fields...).Msg(msg)
}

func Error(msg string, fields ...any) {
	addFields(errLogger.Error().Stack(), fields...).Msg(msg)
}

func Fatal(msg string, fields ...any) {
	addFields(errLogger.Fatal().Stack(), fields...).Msg(msg)
}

func Panic(msg string, fields ...any) {
	addFields(errLogger.Panic().Stack(), fields...).Msg(msg)
}
