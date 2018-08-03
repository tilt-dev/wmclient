package analytics_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/windmilleng/wmclient/pkg/analytics"
)

func TestFlush(t *testing.T) {
	f := newAnalyticsFixture(t)
	for i := 0; i < 10; i++ {
		f.a.Incr("event", nil)
	}

	f.a.Flush(time.Second)
	if len(f.reqs) != 10 {
		t.Errorf("Expected 10 events sent. Actual: %d", len(f.reqs))
	}
}

type analyticsFixture struct {
	t    *testing.T
	a    analytics.Analytics
	reqs []*http.Request
}

func newAnalyticsFixture(t *testing.T) *analyticsFixture {
	f := &analyticsFixture{t: t}
	a := analytics.NewRemoteAnalytics(f, "test-app", "/report", "random-user", true)
	f.a = a
	return f
}

func (f *analyticsFixture) Do(req *http.Request) (*http.Response, error) {
	f.reqs = append(f.reqs, req)
	return &http.Response{StatusCode: 200}, nil
}
