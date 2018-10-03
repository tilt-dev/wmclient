package analytics

import (
	"net/http"
	"testing"
	"time"
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

func TestUserID(t *testing.T) {
	uname := []byte("Linux Sleepy 4.15.0-34-generic #37-Ubuntu SMP Mon Aug 27 15:21:48 UTC 2018 x86_64 x86_64 x86_64 GNU/Linux")
	hash := hashMD5(uname)
	expected := "39894d36bdd53cfe67fca4e7f570e7ff"
	if hash != expected {
		t.Errorf("Expected %q, actual %q", expected, hash)
	}
}

type analyticsFixture struct {
	t    *testing.T
	a    Analytics
	reqs []*http.Request
}

func newAnalyticsFixture(t *testing.T) *analyticsFixture {
	f := &analyticsFixture{t: t}
	a := NewRemoteAnalytics(f, "test-app", "/report", "random-user", true)
	f.a = a
	return f
}

func (f *analyticsFixture) Do(req *http.Request) (*http.Response, error) {
	f.reqs = append(f.reqs, req)
	return &http.Response{StatusCode: 200}, nil
}
