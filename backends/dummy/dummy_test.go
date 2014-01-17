package dummy

import (
	"fmt"
	"github.com/jarcoal/ego/testutils"
	"testing"
)

func TestDummyLogging(t *testing.T) {
	logs := make([]string, 0)

	logger := func(format string, vars ...interface{}) {
		logs = append(logs, fmt.Sprintf(format, vars...))
	}

	b := NewBackend(logger)
	e := testutils.TestEmail()

	b.DispatchEmail(e)

	if len(logs) == 0 {
		t.FailNow()
	}
}
