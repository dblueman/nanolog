package nanolog

var (
   defaultFilter int
   Default, _     = New("", 0)
)

func Filter(filter int) {
   defaultFilter = filter
   Default.Filter(filter)
}

func NamedFilter(filter string) error {
   err := Default.NamedFilter(filter)
   if err != nil {
      return err
   }

   defaultFilter = Default.filter
   return nil
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
