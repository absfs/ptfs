package ptfs_test

import (
	"testing"

	"github.com/absfs/absfs"
	"github.com/absfs/memfs"
	"github.com/absfs/ptfs"
)

func TestPtfs(t *testing.T) {
	mfs, err := memfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}
	pfs, err := ptfs.NewSymlinkFS(mfs)
	if err != nil {
		t.Fatal(err)
	}
	var fs absfs.SymlinkFileSystem
	fs = pfs
	_ = fs
}
