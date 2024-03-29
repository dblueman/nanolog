// Package nanolog implements a lightweight logger, featuring:
// - severity given by API, ie Infof(), Debugf() and related
// - automatic detection of interactive (terminal) or uninteractive (daemon) use
// - appropriate colour prefix for interactive output
// - appropriate syslog prefix for non-interactive output
package nanolog

import (
   "fmt"
   "os"
   "strings"

   "golang.org/x/sys/unix"
)

type Logger struct {
   userPrefix  string
   filter      int
   interactive bool
}

const (
   LevelError = 4
   LevelWarn  = 5
   LevelInfo  = 6
   LevelDebug = 7
   maxLine    = 47*1024

   // uninteractive strings
   inFatalPrefix = "\033[1;31m"
   inErrorPrefix = "\033[1;31m"
   inWarnPrefix  = "\033[1;33m"
   inInfoPrefix  = ""
   inDebugPrefix = "\033[1;36m"
   inSuffix      = "\033[m"

   // uninteractive strings
   unFatalPrefix = "" // default output must be set to crit by "[Unit] SyslogLevel=crit" in systemd unit
   unErrorPrefix = "<3>"
   unWarnPrefix  = "<4>"
   unInfoPrefix  = "<6>"
   unDebugPrefix = "<7>"
   unSuffix      = ""
)

var (
   levels = [...]string{"crit", "error", "warn", "info", "debug"}
)

// filter 0 means use resasonable defaults
func New(userPrefix string, filter int) (*Logger, error) {
   if filter > 7 {
      return nil, fmt.Errorf("invalid log level %d", filter)
   }

   ml := Logger{userPrefix: userPrefix, filter: filter}

   // check if context is interactive or non-interactive
   _, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
   if err == nil {
      ml.interactive = true
   } else if err != unix.ENOTTY {
      return nil, fmt.Errorf("New: %w", err)
   }

   if defaultFilter == 0 {
      if ml.interactive {
         defaultFilter = 6
      } else {
         defaultFilter = 7
      }
   }

   // default filter threshold
   if filter == 0 {
      ml.filter = defaultFilter
   }

   return &ml, nil
}

func (ml *Logger) Filter(filter int) {
   ml.filter = filter
}

func (ml *Logger) NamedFilter(filter string) error {
   for i := range levels {
      if filter != levels[i] {
         continue
      }

      ml.filter = i + 3 // starts at crit
      return nil
   }

   return fmt.Errorf("unknown level %s", filter)
}

func (ml *Logger) Fatal(format string, args ...any) {
   message := fmt.Sprintf(format, args...)
   prefix := inFatalPrefix
   suffix := inSuffix

   // apply user prefix to remaining lines
   message = strings.ReplaceAll(message, "\n", "\n" + ml.userPrefix)

   if !ml.interactive {
      // journalctl splits long messages, losing loglevel
      if len(message) > maxLine {
         message = message[:maxLine] + " (truncated)"
      }

      message = strings.ReplaceAll(message, "\n", "\n" + unFatalPrefix)
      prefix = unFatalPrefix
      suffix = unSuffix
   }

   panic(prefix + ml.userPrefix + message + suffix + "\n")
}

func (ml *Logger) Error(format string, args ...any) {
   if ml.filter < LevelError {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inErrorPrefix
   suffix := inSuffix

   // apply user prefix to remaining lines
   message = strings.ReplaceAll(message, "\n", "\n" + ml.userPrefix)

   if !ml.interactive {
      // journalctl splits long messages, losing loglevel
      if len(message) > maxLine {
         message = message[:maxLine] + " (truncated)"
      }

      // apply syslog prefix to every line
      message = strings.ReplaceAll(message, "\n", "\n" + unErrorPrefix)
      prefix = unErrorPrefix
      suffix = unSuffix
   }

   fmt.Print(prefix + ml.userPrefix + message + suffix + "\n")
}

func (ml *Logger) Warn(format string, args ...any) {
   if ml.filter < LevelWarn {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inWarnPrefix
   suffix := inSuffix

   // apply user prefix to remaining lines
   message = strings.ReplaceAll(message, "\n", "\n" + ml.userPrefix)

   if !ml.interactive {
      // journalctl splits long messages, losing loglevel
      if len(message) > maxLine {
         message = message[:maxLine] + " (truncated)"
      }

      message = strings.ReplaceAll(message, "\n", "\n" + unWarnPrefix)
      prefix = unWarnPrefix
      suffix = unSuffix
   }

   fmt.Print(prefix + ml.userPrefix + message + suffix + "\n")
}

func (ml *Logger) Info(format string, args ...any) {
   if ml.filter < LevelInfo {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inInfoPrefix
   suffix := inSuffix

   // apply user prefix to remaining lines
   message = strings.ReplaceAll(message, "\n", "\n" + ml.userPrefix)

   if !ml.interactive {
      // journalctl splits long messages, losing loglevel
      if len(message) > maxLine {
         message = message[:maxLine] + " (truncated)"
      }

      message = strings.ReplaceAll(message, "\n", "\n" + unInfoPrefix)
      prefix = unInfoPrefix
      suffix = unSuffix
   }

   fmt.Print(prefix + ml.userPrefix + message + suffix + "\n")
}

func (ml *Logger) Debug(format string, args ...any) {
   if ml.filter < LevelDebug {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inDebugPrefix
   suffix := inSuffix

   // apply user prefix to remaining lines
   message = strings.ReplaceAll(message, "\n", "\n" + ml.userPrefix)

   if !ml.interactive {
      // journalctl splits long messages, losing loglevel
      if len(message) > maxLine {
         message = message[:maxLine] + " (truncated)"
      }

      message = strings.ReplaceAll(message, "\n", "\n" + unDebugPrefix)
      prefix = unDebugPrefix
      suffix = unSuffix
   }

   fmt.Print(prefix + ml.userPrefix + message + suffix + "\n")
}
