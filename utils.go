package ptfs

import (
	"github.com/absfs/absfs"
)

// UnwrapFiler unwraps an `absfs.Filer` if it is a pass through filer, and
// returns the underlying `absfs.Filer` object, otherwise it returns the
// argument unmodified.
func UnwrapFiler(fs absfs.Filer) absfs.Filer {
	pfs, ok := fs.(*Filer)
	if ok {
		return pfs.fs
	}
	return fs
}

// UnwrapFS unwraps a `absfs.FileSystem` if it is a pass through filesystem, and
// returns the underlying `absfs.FileSystem` object, otherwise it returns the
// argument unmodified.
func UnwrapFS(fs absfs.FileSystem) absfs.FileSystem {
	pfs, ok := fs.(*FileSystem)
	if ok {
		return pfs.fs
	}
	return fs
}

// SymlinkFileSystem unwraps an `absfs.Filer` if it is a pass through filer, and
// returns the underlying `absfs.Filer` object, otherwise it returns the
// argument unmodified.
func UnwrapSymlinkFS(fs absfs.SymlinkFileSystem) absfs.SymlinkFileSystem {
	pfs, ok := fs.(*SymlinkFileSystem)
	if ok {
		return pfs.sfs
	}
	return fs
}
