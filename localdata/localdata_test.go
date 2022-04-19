package localdata

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeDirIfNotExist(t *testing.T) {

	// skip test if webassembly
	if runtime.GOOS == "js" {
		t.SkipNow()
	}

	t.Run("Test creating directory", func(t *testing.T) {
		// create temp dir
		dir, err := ioutil.TempDir(os.TempDir(), "pw_slippymap_TestMakeDirIfNotExist")
		if err != nil {
			t.Error(err)
		}
		defer os.RemoveAll(dir)

		testDir := path.Join(dir, "testdir")

		// test MakeDirIfNotExist
		err = MakeDirIfNotExist(testDir, 0700)
		require.NoError(t, err)
		defer os.RemoveAll(testDir)

		// ensure dir exists and mode is 0700
		fileInfo, err := os.Stat(testDir)
		require.NoError(t, err)

		// only check mode bits if not windows
		if runtime.GOOS != "windows" {
			assert.Equal(t, fs.FileMode(0700)+fs.ModeDir, fileInfo.Mode())
		}

		// test MakeDirIfNotExist on a directory that already exists
		err = MakeDirIfNotExist(testDir, 0700)
		require.NoError(t, err)
	})
}
