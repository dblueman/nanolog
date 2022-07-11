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

type Microlog struct {
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
func New(userPrefix string, filter int) (*Microlog, error) {
   // check if context is interactive or non-interactive
   _, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)

   ml := Microlog{userPrefix: userPrefix}

   if err == nil {
      ml.interactive = true
   } else if err != unix.ENOTTY {
      return nil, err
   }

   if filter > 7 {
      return nil, fmt.Errorf("invalid log level %d", filter)
   }

   // default filter threshold
   if filter == 0 {
      if ml.interactive {
         ml.filter = 6
      } else {
         ml.filter = 7
      }
   }

   return &ml, nil
}

func (ml *Microlog) Filter(filter int) {
   ml.filter = filter
}

func (ml *Microlog) NamedFilter(filter string) error {
   for i := range levels {
      if filter != levels[i] {
         continue
      }

      ml.filter = i + 3 // starts at crit
      return nil
   }

   return fmt.Errorf("unknown level %s", filter)
}

func (ml *Microlog) Fatal(format string, args ...any) {
   message := fmt.Sprintf(format, args...)
   prefix := inFatalPrefix
   suffix := inSuffix

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

func (ml *Microlog) Error(format string, args ...any) {
   if ml.filter < LevelError {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inErrorPrefix
   suffix := inSuffix

   if !ml.interactive {
      // journalctl splits long messages, losing loglevel
      if len(message) > maxLine {
         message = message[:maxLine] + " (truncated)"
      }

      message = strings.ReplaceAll(message, "\n", "\n" + unErrorPrefix)
      prefix = unErrorPrefix
      suffix = unSuffix
   }

   fmt.Print(prefix + ml.userPrefix + message + suffix + "\n")
}

func (ml *Microlog) Warn(format string, args ...any) {
   if ml.filter < LevelWarn {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inWarnPrefix
   suffix := inSuffix

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

func (ml *Microlog) Info(format string, args ...any) {
   if ml.filter < LevelInfo {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inInfoPrefix
   suffix := inSuffix

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

func (ml *Microlog) Debug(format string, args ...any) {
   if ml.filter < LevelDebug {
      return
   }

   message := fmt.Sprintf(format, args...)
   prefix := inDebugPrefix
   suffix := inSuffix

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
