package nanolog

var (
   defaultFilter int
   Default, _     = New("", 0)
)

func Filter(filter int) {
   Default.Filter(filter)
}

func NamedFilter(filter string) error {
   return Default.NamedFilter(filter)
}

func Fatal(format string, args ...any) {
   Default.Fatal(format, args...)
}

func Error(format string, args ...any) {
   Default.Error(format, args...)
}

func Warn(format string, args ...any) {
   Default.Warn(format, args...)
}

func Info(format string, args ...any) {
   Default.Info(format, args...)
}

func Debug(format string, args ...any) {
   Default.Debug(format, args...)
}
