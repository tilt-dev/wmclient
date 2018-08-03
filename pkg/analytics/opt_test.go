package analytics_test

import (
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

	analytics.SetOptStr("opt-in")
	f.assertOptStatus(analytics.OptIn)

	analytics.SetOptStr("opt-out")
	f.assertOptStatus(analytics.OptOut)

	analytics.SetOptStr("in")
	f.assertOptStatus(analytics.OptIn)

	analytics.SetOptStr("out")
	f.assertOptStatus(analytics.OptOut)

	analytics.SetOptStr("foo")
	f.assertOptStatus(analytics.OptDefault)
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

	os.Setenv("WINDMILL_DIR", dir)
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
