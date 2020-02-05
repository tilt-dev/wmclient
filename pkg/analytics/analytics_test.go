package analytics

import (
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestGlobalTags(t *testing.T) {
	f := newAnalyticsFixture(t, WithGlobalTags(map[string]string{"fruit": "pomelo"}))
	f.a.Incr("event", map[string]string{"season": "summer"})
	f.a.Flush(time.Second)

	if len(f.reqs) != 1 {
		t.Fatalf("Expected 1 event sent. Actual: %d", len(f.reqs))
	}

	expected := `{"fruit":"pomelo","machine":"random-machine","name":"test-app.event","season":"summer","user":"random-user"}`
	body, err := ioutil.ReadAll(f.reqs[0].Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != expected {
		t.Errorf("Request body did not match\nExpected: %s\nActual: %s", expected, string(body))
	}

	tag, ok := f.a.GlobalTag("fruit")
	assert.True(t, ok)
	assert.Equal(t, "pomelo", tag)
}

func TestIncrWithoutGlobalTags(t *testing.T) {
	f := newAnalyticsFixture(t, WithGlobalTags(map[string]string{"fruit": "pomelo"}))
	f.a.WithoutGlobalTags().Incr("event", map[string]string{"season": "summer"})
	f.a.Flush(time.Second)

	if len(f.reqs) != 1 {
		t.Fatalf("Expected 1 event sent. Actual: %d", len(f.reqs))
	}

	expected := `{"name":"test-app.event","season":"summer"}`
	body, err := ioutil.ReadAll(f.reqs[0].Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != expected {
		t.Errorf("Request body did not match\nExpected: %s\nActual: %s", expected, string(body))
	}
}

type analyticsFixture struct {
	t    *testing.T
	a    Analytics
	reqs []*http.Request
	mu   sync.Mutex
}

func newAnalyticsFixture(t *testing.T, fOptions ...Option) *analyticsFixture {
	f := &analyticsFixture{t: t}
	options := []Option{
		WithHTTPClient(f),
		WithReportURL("/report"),
		WithUserID("random-user"),
		WithMachineID("random-machine"),
		WithEnabled(true),
	}
	options = append(options, fOptions...)
	a, err := NewRemoteAnalytics("test-app", options...)
	if err != nil {
		t.Fatal(err)
	}
	f.a = a
	return f
}

func (f *analyticsFixture) Do(req *http.Request) (*http.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reqs = append(f.reqs, req)
	return &http.Response{StatusCode: 200}, nil
}
