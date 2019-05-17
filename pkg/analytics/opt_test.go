package analytics_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"

	"github.com/windmilleng/wmclient/pkg/analytics"
)

func TestString(t *testing.T) {
	if analytics.OptIn.String() != "opt-in" {
		t.Errorf("Expected opt-in, actual: %s", analytics.OptIn)
	}
}

func TestSetOptStr(t *testing.T) {
	f := setup(t)
	defer f.tearDown()
	f.assertOptStatus(analytics.OptDefault)

	for _, v := range []struct {
		s string
		opt analytics.Opt
	}{
		{"opt-in", analytics.OptIn},
		{"opt-out", analytics.OptOut},
		{"in", analytics.OptIn},
		{"out", analytics.OptOut},
	} {
		opt, err := analytics.SetOptStr(v.s)
		if assert.NoError(t, err) {
			assert.Equal(t, v.opt, opt)
			f.assertOptStatus(v.opt)
		}
	}

	_, err := analytics.SetOptStr("foo")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "unknown analytics opt: \"foo\"")
	}
}

func TestSetOpt(t *testing.T) {
	f := setup(t)
	defer f.tearDown()

	f.assertOptStatus(analytics.OptDefault)

	analytics.SetOpt(analytics.OptIn)
	f.assertOptStatus(analytics.OptIn)

	analytics.SetOpt(analytics.OptOut)
	f.assertOptStatus(analytics.OptOut)

	analytics.SetOpt(99999)
	f.assertOptStatus(analytics.OptDefault)
}

type fixture struct {
	t              *testing.T
	dir            string
	oldWindmillDir string
}

func setup(t *testing.T) *fixture {
	oldWindmillDir := os.Getenv("WINDMILL_DIR")
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatalf("Error making temp dir: %v", err)
	}

	err = os.Setenv("WINDMILL_DIR", dir)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return &fixture{t: t, dir: dir, oldWindmillDir: oldWindmillDir}
}

func (f *fixture) tearDown() {
	os.RemoveAll(f.dir)
	os.Setenv("WINDMILL_DIR", f.oldWindmillDir)
}

func (f *fixture) assertOptStatus(expected analytics.Opt) {
	actual, err := analytics.OptStatus()
	if err != nil {
		f.t.Fatal(err)
	}
	if actual != expected {
		f.t.Errorf("got opt status %v, expected %v", actual, expected)
	}
}
