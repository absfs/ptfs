package ptfs_test

import (
	"testing"

	"github.com/absfs/absfs"
	"github.com/absfs/fstesting"
	"github.com/absfs/memfs"
	"github.com/absfs/ptfs"
)

func TestPtFSFileSystem(t *testing.T) {
	baseFS, err := memfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}

	suite := &fstesting.WrapperSuite{
		Factory: func(base absfs.FileSystem) (absfs.FileSystem, error) {
			return ptfs.NewFS(base)
		},
		BaseFS: baseFS,
		Name:   "ptfs.FileSystem",
	}
	suite.Run(t)
}

func TestPtFSSymlinkFileSystem(t *testing.T) {
	baseFS, err := memfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}

	suite := &fstesting.WrapperSuite{
		Factory: func(base absfs.FileSystem) (absfs.FileSystem, error) {
			// memfs.FileSystem implements SymlinkFileSystem
			// so we can type assert it
			sfs, ok := base.(absfs.SymlinkFileSystem)
			if !ok {
				t.Fatal("base filesystem does not implement SymlinkFileSystem")
			}
			return ptfs.NewSymlinkFS(sfs)
		},
		BaseFS: baseFS,
		Name:   "ptfs.SymlinkFileSystem",
	}
	suite.Run(t)
}
