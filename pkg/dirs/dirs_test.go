package dirs

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestTiltDevDir(t *testing.T) {
	emptyPath := ""
	oldWmdaemonHome := os.Getenv("WMDAEMON_HOME")
	oldHome := os.Getenv("HOME")
	oldTiltDevDir := os.Getenv("WINDMILL_DIR")
	defer os.Setenv("WMDAEMON_HOME", oldWmdaemonHome)
	defer os.Setenv("HOME", oldHome)
	defer os.Setenv("WINDMILL_DIR", oldTiltDevDir)
	tmpHome := os.TempDir()

	f := setup(t)

	os.Setenv("HOME", tmpHome)

	os.Setenv("WMDAEMON_HOME", emptyPath)
	f.assertTiltDevDir(path.Join(tmpHome, ".tilt-dev"), "empty .windmill")

	os.Mkdir(filepath.Join(tmpHome, ".windmill"), 0755)
	f.assertTiltDevDir(path.Join(tmpHome, ".windmill"), "populated .windmill")

	tmpWmdaemonHome := os.TempDir()
	os.Setenv("WMDAEMON_HOME", tmpWmdaemonHome)
	f.assertTiltDevDir(tmpWmdaemonHome, "tmp WMDAEMON_HOME")

	nonExistentWmdaemonHome := path.Join(tmpWmdaemonHome, "foo")
	os.Setenv("WMDAEMON_HOME", nonExistentWmdaemonHome)
	f.assertTiltDevDir(nonExistentWmdaemonHome, "nonexistent WMDAEMON_HOME")

	wmDir := os.TempDir()
	os.Setenv("WINDMILL_DIR", wmDir)
	f.assertTiltDevDir(nonExistentWmdaemonHome, "prefer WMDAEMON_HOME") // prefer WMDAEMON_HOME

	os.Unsetenv("WMDAEMON_HOME")
	f.assertTiltDevDir(wmDir, "no WMDAEMON_HOME")
}

func TestOpenFile(t *testing.T) {
	tmp, _ := ioutil.TempDir("", t.Name())
	defer os.RemoveAll(tmp)
	dir := NewTiltDevDirAt(tmp)

	fp, err := dir.OpenFile("inner/a.txt", os.O_WRONLY|os.O_CREATE, os.FileMode(0700))
	if err != nil {
		t.Fatal(err)
	}

	fp.Write([]byte("hello"))
	fp.Close()

	contents, err := dir.ReadFile("inner/a.txt")
	if err != nil {
		t.Fatal(err)
	}

	if contents != "hello" {
		t.Errorf("Expected %q. Actual: %q", "hello", contents)
	}
}

type fixture struct {
	t *testing.T
}

func setup(t *testing.T) *fixture {
	return &fixture{t: t}
}

func (f *fixture) assertTiltDevDir(expected, testCase string) {
	actual, err := GetTiltDevDir()
	if err != nil {
		f.t.Error(err)
	}

	// NOTE(maia): filepath behavior is weird on macOS, use abs path to mitigate
	absExpected, err := filepath.Abs(expected)
	if err != nil {
		f.t.Error("[filepath.Abs]", err)
	}

	if actual != absExpected {
		f.t.Errorf("[TEST CASE: %s] got windmill dir %q; expected %q", testCase, actual, absExpected)
	}
}
