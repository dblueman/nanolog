package nanolog

import (
   "fmt"
   "os"
   "strings"

   "golang.org/x/sys/unix"
)

const (
   LevelError = 4
   LevelWarn  = 5
   LevelInfo  = 6
   LevelDebug = 7
)

var (
   levels      = []string{"crit", "error", "warn", "info", "debug"}
   fatalPrefix string
   errorPrefix string
   warnPrefix  string
   infoPrefix  string
   debugPrefix string
   suffix      string
   level       int
   interactive bool
)

func init() {
   _, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
   if err == nil {
      interactive = true
      level       = 6
      fatalPrefix = "\033[1;31m"
      errorPrefix = "\033[1;31m"
      warnPrefix  = "\033[1;33m"
      infoPrefix  = ""
      debugPrefix = "\033[1;36m"
      suffix      = "\033[m"
   } else {
      level       = 7
      fatalPrefix = "" // default output must be set to crit by "[Unit] SyslogLevel=crit" in systemd unit
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

func SetMinimumStr(_level string) error {
   for i := range levels {
      if _level != levels[i] {
         continue
      }

      level = i + 3 // starts at crit
      return nil
   }

   return fmt.Errorf("unknown level %s", _level)
}

func Fatal(format string, args ...interface{}) {
   message := fmt.Sprintf(format, args...)
   if !interactive {
      message = strings.ReplaceAll(message, "\n", "\n" + fatalPrefix)
   }
   panic(fatalPrefix + message + suffix + "\n")
}

func Error(format string, args ...interface{}) {
   if level < LevelError {
      return
   }

   message := fmt.Sprintf(format, args...)
   if !interactive {
      message = strings.ReplaceAll(message, "\n", "\n" + errorPrefix)
   }
   fmt.Print(errorPrefix + message + suffix + "\n")
}

func Warn(format string, args ...interface{}) {
   if level < LevelWarn {
      return
   }

   message := fmt.Sprintf(format, args...)
   if !interactive {
      message = strings.ReplaceAll(message, "\n", "\n" + warnPrefix)
   }
   fmt.Print(warnPrefix + message + suffix + "\n")
}

func Info(format string, args ...interface{}) {
   if level < LevelInfo {
      return
   }

   message := fmt.Sprintf(format, args...)
   if !interactive {
      message = strings.ReplaceAll(message, "\n", "\n" + infoPrefix)
   }
   fmt.Print(infoPrefix + message + suffix + "\n")
}

func Debug(format string, args ...interface{}) {
   if level < LevelDebug {
      return
   }

   message := fmt.Sprintf(format, args...)
   if !interactive {
      message = strings.ReplaceAll(message, "\n", "\n" + debugPrefix)
   }
   fmt.Print(debugPrefix + message + suffix + "\n")
}
