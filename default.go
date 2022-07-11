package nanolog

var (
   defaultFilter int
   simple, _     = New("", 0)
)

func Filter(filter int) {
   simple.Filter(filter)
}

func NamedFilter(filter string) error {
   return simple.NamedFilter(filter)
}

func Fatal(format string, args ...any) {
   simple.Fatal(format, args...)
}

func Error(format string, args ...any) {
   simple.Error(format, args...)
}

func Warn(format string, args ...any) {
   simple.Warn(format, args...)
}

func Info(format string, args ...any) {
   simple.Info(format, args...)
}

func Debug(format string, args ...any) {
   simple.Debug(format, args...)
}
