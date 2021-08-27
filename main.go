package nanolog

import (
   "fmt"
   "os"

   "golang.org/x/sys/unix"
)

const (
   LevelError = 4
   LevelWarn  = 5
   LevelInfo  = 6
   LevelDebug = 7
)

var (
   level       = 7
   fatalPrefix string
   errorPrefix string
   warnPrefix  string
   infoPrefix  string
   debugPrefix string
   suffix      string
)

func init() {
   _, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
   if err == nil {
      fatalPrefix = "\033[1;31m"
      errorPrefix = "\033[1;31m"
      warnPrefix  = "\033[1;33m"
      infoPrefix  = "\033[1;32m"
      debugPrefix = "\033[1;36m"
      suffix      = "\033[m"
   } else {
      fatalPrefix = "<2>"
      errorPrefix = "<3>"
      warnPrefix  = "<4>"
      infoPrefix  = "<6>"
      debugPrefix = "<7>"
      suffix      = ""
   }
}

func SetMinimum(_level int) {
   level = _level
}

func Fatal(format string, args ...interface{}) {
   message := fmt.Sprintf(fatalPrefix + format + suffix + "\n", args...)
   panic(message)
}

func Error(format string, args ...interface{}) {
   if level < LevelError {
      return
   }

   fmt.Printf(errorPrefix + format + suffix + "\n", args...)
}

func Warn(format string, args ...interface{}) {
   if level < LevelWarn {
      return
   }

   fmt.Printf(warnPrefix + format + suffix + "\n", args...)
}

func Info(format string, args ...interface{}) {
   if level < LevelInfo {
      return
   }

   fmt.Printf(infoPrefix + format + suffix + "\n", args...)
}

func Debug(format string, args ...interface{}) {
   if level < LevelDebug {
      return
   }

   fmt.Printf(debugPrefix + format + suffix + "\n", args...)
}
