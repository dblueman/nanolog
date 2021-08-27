package nanolog

import (
   "testing"
)

func Test(t *testing.T) {
   Error("should see this error")
   Warn("should see this %s", "warning")
   Info("should see this info")
   Debug("should see this debug")
   SetMinimum(LevelWarn)
   Debug("should not see this debug message")
   Info("should not see this info message")
   Warn("should see this warning")
}
