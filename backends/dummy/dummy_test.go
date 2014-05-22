package dummy

import (
	"fmt"
	"github.com/jarcoal/ego/testutils"
	"testing"
)

// TestDummyLogging checks the logging functionality
func TestDummyLogging(t *testing.T) {
	logs := make([]string, 0)

	logger := func(format string, vars ...interface{}) {
		logs = append(logs, fmt.Sprintf(format, vars...))
	}

	b := NewBackend(logger)
	e := testutils.TestEmail()

	b.SendEmail(e)

	if len(logs) == 0 {
		t.FailNow()
	}
}
