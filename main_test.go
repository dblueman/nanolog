package nanolog

import (
   "testing"
)

func Test(t *testing.T) {
   Error("should see this error")
   Warn("should see this %s", "warning")
   Info("should see this info")
   Debug("should see this debug")

   l, err := New("[prefix] ", LevelWarn)
   if err != nil {
      t.Fatal(err)
   }

   l.Debug("should not see this debug message")
   l.Info("should not see this info message")
   l.Warn("should see this warning")
}
