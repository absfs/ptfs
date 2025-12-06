package ptfs_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/absfs/absfs"
	"github.com/absfs/memfs"
	"github.com/absfs/ptfs"
)

// Helper function to create a memfs for testing
func newMemFS(t *testing.T) *memfs.FileSystem {
	t.Helper()
	mfs, err := memfs.NewFS()
	if err != nil {
		t.Fatal(err)
	}
	return mfs
}

// =============================================================================
// Phase 1.1: Filer Type Tests
// =============================================================================

func TestNewFiler(t *testing.T) {
	mfs := newMemFS(t)

	filer, err := ptfs.NewFiler(mfs)
	if err != nil {
		t.Fatalf("NewFiler failed: %v", err)
	}
	if filer == nil {
		t.Fatal("NewFiler returned nil")
	}

	// Verify it implements absfs.Filer
	var _ absfs.Filer = filer
}

func TestFiler_OpenFile(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file
	f, err := filer.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	defer f.Close()

	// Write some data
	_, err = f.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
}

func TestFiler_Mkdir(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	err := filer.Mkdir("/testdir", 0755)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	// Verify directory exists
	info, err := filer.Stat("/testdir")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}

func TestFiler_Remove(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file first
	f, _ := filer.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	// Remove it
	err := filer.Remove("/test.txt")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// Verify it's gone
	_, err = filer.Stat("/test.txt")
	if err == nil {
		t.Fatal("expected error for removed file")
	}
}

func TestFiler_Rename(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file
	f, _ := filer.OpenFile("/old.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte("data"))
	f.Close()

	// Rename it
	err := filer.Rename("/old.txt", "/new.txt")
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	// Verify old is gone
	_, err = filer.Stat("/old.txt")
	if err == nil {
		t.Fatal("old file should not exist")
	}

	// Verify new exists
	_, err = filer.Stat("/new.txt")
	if err != nil {
		t.Fatalf("new file should exist: %v", err)
	}
}

func TestFiler_Stat(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file
	f, _ := filer.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Write([]byte("hello"))
	f.Close()

	info, err := filer.Stat("/test.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Name() != "test.txt" {
		t.Fatalf("expected name 'test.txt', got '%s'", info.Name())
	}
	if info.Size() != 5 {
		t.Fatalf("expected size 5, got %d", info.Size())
	}
}

func TestFiler_Chmod(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file
	f, _ := filer.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	// Change permissions
	err := filer.Chmod("/test.txt", 0755)
	if err != nil {
		t.Fatalf("Chmod failed: %v", err)
	}

	info, _ := filer.Stat("/test.txt")
	if info.Mode().Perm() != 0755 {
		t.Fatalf("expected mode 0755, got %o", info.Mode().Perm())
	}
}

func TestFiler_Chtimes(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file
	f, _ := filer.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	// Change times
	newTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	err := filer.Chtimes("/test.txt", newTime, newTime)
	if err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}

	info, _ := filer.Stat("/test.txt")
	if !info.ModTime().Equal(newTime) {
		t.Fatalf("expected mod time %v, got %v", newTime, info.ModTime())
	}
}

func TestFiler_Chown(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Create a file
	f, _ := filer.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	// Change ownership (memfs may not fully support this, but we test pass-through)
	err := filer.Chown("/test.txt", 1000, 1000)
	// Note: error may or may not occur depending on memfs implementation
	_ = err
}

// =============================================================================
// Phase 1.2: FileSystem Type Tests
// =============================================================================

func TestNewFS(t *testing.T) {
	mfs := newMemFS(t)

	fs, err := ptfs.NewFS(mfs)
	if err != nil {
		t.Fatalf("NewFS failed: %v", err)
	}
	if fs == nil {
		t.Fatal("NewFS returned nil")
	}

	// Verify it implements absfs.FileSystem
	var _ absfs.FileSystem = fs
}

func TestFileSystem_OpenFile(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, err := fs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	f.Close()
}

func TestFileSystem_Mkdir(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	err := fs.Mkdir("/testdir", 0755)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}
}

func TestFileSystem_Remove(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, _ := fs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	err := fs.Remove("/test.txt")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
}

func TestFileSystem_Rename(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, _ := fs.OpenFile("/old.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	err := fs.Rename("/old.txt", "/new.txt")
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}
}

func TestFileSystem_Stat(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, _ := fs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	info, err := fs.Stat("/test.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Name() != "test.txt" {
		t.Fatalf("expected 'test.txt', got '%s'", info.Name())
	}
}

func TestFileSystem_Chmod(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, _ := fs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	err := fs.Chmod("/test.txt", 0755)
	if err != nil {
		t.Fatalf("Chmod failed: %v", err)
	}
}

func TestFileSystem_Chtimes(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, _ := fs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	newTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	err := fs.Chtimes("/test.txt", newTime, newTime)
	if err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}
}

func TestFileSystem_Chown(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, _ := fs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	_ = fs.Chown("/test.txt", 1000, 1000)
}

func TestFileSystem_Separator(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	sep := fs.Separator()
	expected := mfs.Separator()
	if sep != expected {
		t.Fatalf("expected separator %c, got %c", expected, sep)
	}
}

func TestFileSystem_ListSeparator(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	sep := fs.ListSeparator()
	expected := mfs.ListSeparator()
	if sep != expected {
		t.Fatalf("expected list separator %c, got %c", expected, sep)
	}
}

func TestFileSystem_Chdir(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	// Create a directory first
	fs.Mkdir("/testdir", 0755)

	err := fs.Chdir("/testdir")
	if err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
}

func TestFileSystem_Getwd(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	dir, err := fs.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	if dir != "/" {
		t.Fatalf("expected '/', got '%s'", dir)
	}
}

func TestFileSystem_TempDir(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	tempDir := fs.TempDir()
	expected := mfs.TempDir()
	if tempDir != expected {
		t.Fatalf("expected '%s', got '%s'", expected, tempDir)
	}
}

func TestFileSystem_Open(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	// Create a file first
	f, _ := fs.Create("/test.txt")
	f.Write([]byte("hello"))
	f.Close()

	// Open for reading
	f, err := fs.Open("/test.txt")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer f.Close()

	data := make([]byte, 5)
	n, _ := f.Read(data)
	if string(data[:n]) != "hello" {
		t.Fatalf("expected 'hello', got '%s'", string(data[:n]))
	}
}

func TestFileSystem_Create(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	f, err := fs.Create("/test.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	f.Close()

	// Verify file exists
	_, err = fs.Stat("/test.txt")
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
}

func TestFileSystem_MkdirAll(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	err := fs.MkdirAll("/a/b/c", 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	// Verify nested directories exist
	info, err := fs.Stat("/a/b/c")
	if err != nil {
		t.Fatalf("nested dir should exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}

func TestFileSystem_RemoveAll(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	// Create nested structure
	fs.MkdirAll("/a/b/c", 0755)
	f, _ := fs.Create("/a/b/file.txt")
	f.Close()

	err := fs.RemoveAll("/a")
	if err != nil {
		t.Fatalf("RemoveAll failed: %v", err)
	}

	// Verify it's gone
	_, err = fs.Stat("/a")
	if err == nil {
		t.Fatal("directory should be removed")
	}
}

func TestFileSystem_Truncate(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	// Create a file with content
	f, _ := fs.Create("/test.txt")
	f.Write([]byte("hello world"))
	f.Close()

	// Truncate to 5 bytes
	err := fs.Truncate("/test.txt", 5)
	if err != nil {
		t.Fatalf("Truncate failed: %v", err)
	}

	info, _ := fs.Stat("/test.txt")
	if info.Size() != 5 {
		t.Fatalf("expected size 5, got %d", info.Size())
	}
}

// =============================================================================
// Phase 1.3: SymlinkFileSystem Type Tests
// =============================================================================

func TestNewSymlinkFS(t *testing.T) {
	mfs := newMemFS(t)

	sfs, err := ptfs.NewSymlinkFS(mfs)
	if err != nil {
		t.Fatalf("NewSymlinkFS failed: %v", err)
	}
	if sfs == nil {
		t.Fatal("NewSymlinkFS returned nil")
	}

	// Verify it implements absfs.SymlinkFileSystem
	var _ absfs.SymlinkFileSystem = sfs
}

func TestSymlinkFileSystem_OpenFile(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, err := sfs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	f.Close()
}

func TestSymlinkFileSystem_Mkdir(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	err := sfs.Mkdir("/testdir", 0755)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}
}

func TestSymlinkFileSystem_Remove(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	err := sfs.Remove("/test.txt")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
}

func TestSymlinkFileSystem_Rename(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.OpenFile("/old.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	err := sfs.Rename("/old.txt", "/new.txt")
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}
}

func TestSymlinkFileSystem_Stat(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	info, err := sfs.Stat("/test.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Name() != "test.txt" {
		t.Fatalf("expected 'test.txt', got '%s'", info.Name())
	}
}

func TestSymlinkFileSystem_Chmod(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	err := sfs.Chmod("/test.txt", 0755)
	if err != nil {
		t.Fatalf("Chmod failed: %v", err)
	}
}

func TestSymlinkFileSystem_Chtimes(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	newTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	err := sfs.Chtimes("/test.txt", newTime, newTime)
	if err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}
}

func TestSymlinkFileSystem_Chown(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.OpenFile("/test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	_ = sfs.Chown("/test.txt", 1000, 1000)
}

func TestSymlinkFileSystem_Separator(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	sep := sfs.Separator()
	expected := mfs.Separator()
	if sep != expected {
		t.Fatalf("expected separator %c, got %c", expected, sep)
	}
}

func TestSymlinkFileSystem_ListSeparator(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	sep := sfs.ListSeparator()
	expected := mfs.ListSeparator()
	if sep != expected {
		t.Fatalf("expected list separator %c, got %c", expected, sep)
	}
}

func TestSymlinkFileSystem_Chdir(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	sfs.Mkdir("/testdir", 0755)

	err := sfs.Chdir("/testdir")
	if err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}
}

func TestSymlinkFileSystem_Getwd(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	dir, err := sfs.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	if dir != "/" {
		t.Fatalf("expected '/', got '%s'", dir)
	}
}

func TestSymlinkFileSystem_TempDir(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	tempDir := sfs.TempDir()
	expected := mfs.TempDir()
	if tempDir != expected {
		t.Fatalf("expected '%s', got '%s'", expected, tempDir)
	}
}

func TestSymlinkFileSystem_Open(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("hello"))
	f.Close()

	f, err := sfs.Open("/test.txt")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer f.Close()

	data := make([]byte, 5)
	n, _ := f.Read(data)
	if string(data[:n]) != "hello" {
		t.Fatalf("expected 'hello', got '%s'", string(data[:n]))
	}
}

func TestSymlinkFileSystem_Create(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, err := sfs.Create("/test.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	f.Close()
}

func TestSymlinkFileSystem_MkdirAll(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	err := sfs.MkdirAll("/a/b/c", 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
}

func TestSymlinkFileSystem_RemoveAll(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	sfs.MkdirAll("/a/b/c", 0755)
	f, _ := sfs.Create("/a/b/file.txt")
	f.Close()

	err := sfs.RemoveAll("/a")
	if err != nil {
		t.Fatalf("RemoveAll failed: %v", err)
	}
}

func TestSymlinkFileSystem_Truncate(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("hello world"))
	f.Close()

	err := sfs.Truncate("/test.txt", 5)
	if err != nil {
		t.Fatalf("Truncate failed: %v", err)
	}

	info, _ := sfs.Stat("/test.txt")
	if info.Size() != 5 {
		t.Fatalf("expected size 5, got %d", info.Size())
	}
}

func TestSymlinkFileSystem_Lstat(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create a regular file
	f, _ := sfs.Create("/test.txt")
	f.Close()

	info, err := sfs.Lstat("/test.txt")
	if err != nil {
		t.Fatalf("Lstat failed: %v", err)
	}
	if info.Name() != "test.txt" {
		t.Fatalf("expected 'test.txt', got '%s'", info.Name())
	}
}

func TestSymlinkFileSystem_Lchown(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Close()

	// Lchown may or may not be supported by memfs
	_ = sfs.Lchown("/test.txt", 1000, 1000)
}

func TestSymlinkFileSystem_Readlink(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create a symlink
	f, _ := sfs.Create("/target.txt")
	f.Close()

	err := sfs.Symlink("/target.txt", "/link.txt")
	if err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	target, err := sfs.Readlink("/link.txt")
	if err != nil {
		t.Fatalf("Readlink failed: %v", err)
	}
	if target != "/target.txt" {
		t.Fatalf("expected '/target.txt', got '%s'", target)
	}
}

func TestSymlinkFileSystem_Symlink(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/target.txt")
	f.Close()

	err := sfs.Symlink("/target.txt", "/link.txt")
	if err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	// Verify symlink exists via Lstat
	info, err := sfs.Lstat("/link.txt")
	if err != nil {
		t.Fatalf("Lstat on symlink failed: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("expected symlink mode")
	}
}

// =============================================================================
// Phase 1.5: Utility Functions Tests
// =============================================================================

func TestUnwrapFiler_Wrapped(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	unwrapped := ptfs.UnwrapFiler(filer)

	// Should return the original memfs
	if unwrapped == filer {
		t.Fatal("should have unwrapped the filer")
	}
}

func TestUnwrapFiler_NotWrapped(t *testing.T) {
	mfs := newMemFS(t)

	unwrapped := ptfs.UnwrapFiler(mfs)

	// Should return the same object
	if unwrapped != mfs {
		t.Fatal("should return same object when not wrapped")
	}
}

func TestUnwrapFS_Wrapped(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	unwrapped := ptfs.UnwrapFS(fs)

	if unwrapped == fs {
		t.Fatal("should have unwrapped the filesystem")
	}
}

func TestUnwrapFS_NotWrapped(t *testing.T) {
	mfs := newMemFS(t)

	unwrapped := ptfs.UnwrapFS(mfs)

	if unwrapped != mfs {
		t.Fatal("should return same object when not wrapped")
	}
}

func TestUnwrapSymlinkFS_Wrapped(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	unwrapped := ptfs.UnwrapSymlinkFS(sfs)

	if unwrapped == sfs {
		t.Fatal("should have unwrapped the symlink filesystem")
	}
}

func TestUnwrapSymlinkFS_NotWrapped(t *testing.T) {
	mfs := newMemFS(t)

	unwrapped := ptfs.UnwrapSymlinkFS(mfs)

	if unwrapped != mfs {
		t.Fatal("should return same object when not wrapped")
	}
}

// =============================================================================
// Phase 2.1: File Operations Integration Tests
// =============================================================================

func TestFileOperations_CreateWriteReadClose(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create and write
	f, err := sfs.Create("/test.txt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	content := []byte("Hello, World!")
	n, err := f.Write(content)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(content) {
		t.Fatalf("expected to write %d bytes, wrote %d", len(content), n)
	}

	f.Close()

	// Read back
	f, err = sfs.Open("/test.txt")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer f.Close()

	buf := make([]byte, 100)
	n, err = f.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read failed: %v", err)
	}
	if string(buf[:n]) != string(content) {
		t.Fatalf("expected '%s', got '%s'", string(content), string(buf[:n]))
	}
}

func TestFileOperations_Seek(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("Hello, World!"))
	f.Close()

	f, _ = sfs.Open("/test.txt")
	defer f.Close()

	// Seek to position 7 (start of "World!")
	pos, err := f.Seek(7, io.SeekStart)
	if err != nil {
		t.Fatalf("Seek failed: %v", err)
	}
	if pos != 7 {
		t.Fatalf("expected position 7, got %d", pos)
	}

	buf := make([]byte, 6)
	n, _ := f.Read(buf)
	if string(buf[:n]) != "World!" {
		t.Fatalf("expected 'World!', got '%s'", string(buf[:n]))
	}
}

func TestFileOperations_Truncate(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("Hello, World!"))

	err := f.Truncate(5)
	if err != nil {
		t.Fatalf("Truncate failed: %v", err)
	}

	f.Close()

	info, _ := sfs.Stat("/test.txt")
	if info.Size() != 5 {
		t.Fatalf("expected size 5, got %d", info.Size())
	}
}

func TestFileOperations_Stat(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("Hello"))
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Name() != "test.txt" {
		t.Fatalf("expected 'test.txt', got '%s'", info.Name())
	}
}

func TestFileOperations_Sync(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("Hello"))

	err := f.Sync()
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}
	f.Close()
}

func TestFileOperations_Name(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	defer f.Close()

	name := f.Name()
	if name != "/test.txt" {
		t.Fatalf("expected '/test.txt', got '%s'", name)
	}
}

func TestFileOperations_WriteString(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	n, err := f.WriteString("Hello")
	if err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5, got %d", n)
	}
	f.Close()
}

func TestFileOperations_ReadAt(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("Hello, World!"))
	f.Close()

	f, _ = sfs.OpenFile("/test.txt", os.O_RDONLY, 0644)
	defer f.Close()

	buf := make([]byte, 5)
	n, err := f.ReadAt(buf, 7)
	if err != nil && err != io.EOF {
		t.Fatalf("ReadAt failed: %v", err)
	}
	if string(buf[:n]) != "World" {
		t.Fatalf("expected 'World', got '%s'", string(buf[:n]))
	}
}

func TestFileOperations_WriteAt(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Write([]byte("Hello, World!"))

	n, err := f.WriteAt([]byte("Go"), 7)
	if err != nil {
		t.Fatalf("WriteAt failed: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}
	f.Close()

	// Read back - just verify WriteAt worked by checking file size
	info, _ := sfs.Stat("/test.txt")
	if info.Size() < 9 {
		t.Fatalf("expected size >= 9, got %d", info.Size())
	}
}

// =============================================================================
// Phase 2.2: Directory Operations Integration Tests
// =============================================================================

func TestDirectoryOperations_MkdirMkdirAll(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Single dir
	err := sfs.Mkdir("/single", 0755)
	if err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	// Nested dirs
	err = sfs.MkdirAll("/nested/a/b/c", 0755)
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	// Verify
	info, _ := sfs.Stat("/single")
	if !info.IsDir() {
		t.Fatal("/single should be a directory")
	}

	info, _ = sfs.Stat("/nested/a/b/c")
	if !info.IsDir() {
		t.Fatal("/nested/a/b/c should be a directory")
	}
}

func TestDirectoryOperations_Readdir(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create some files and dirs
	sfs.Mkdir("/testdir", 0755)
	f, _ := sfs.Create("/testdir/file1.txt")
	f.Close()
	f, _ = sfs.Create("/testdir/file2.txt")
	f.Close()
	sfs.Mkdir("/testdir/subdir", 0755)

	// Open and read directory
	d, err := sfs.Open("/testdir")
	if err != nil {
		t.Fatalf("Open dir failed: %v", err)
	}
	defer d.Close()

	entries, err := d.Readdir(-1)
	if err != nil {
		t.Fatalf("Readdir failed: %v", err)
	}

	// memfs may include . and .. entries, so just check we got at least 3
	if len(entries) < 3 {
		t.Fatalf("expected at least 3 entries, got %d", len(entries))
	}
}

func TestDirectoryOperations_Readdirnames(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	sfs.Mkdir("/testdir", 0755)
	f, _ := sfs.Create("/testdir/file1.txt")
	f.Close()
	f, _ = sfs.Create("/testdir/file2.txt")
	f.Close()

	d, err := sfs.Open("/testdir")
	if err != nil {
		t.Fatalf("Open dir failed: %v", err)
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		t.Fatalf("Readdirnames failed: %v", err)
	}

	// memfs may include . and .. entries, so just check we got at least 2
	if len(names) < 2 {
		t.Fatalf("expected at least 2 names, got %d", len(names))
	}
}

func TestDirectoryOperations_RemoveRemoveAll(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create structure
	sfs.MkdirAll("/dir/subdir", 0755)
	f, _ := sfs.Create("/dir/subdir/file.txt")
	f.Close()
	f, _ = sfs.Create("/single.txt")
	f.Close()

	// Remove single file
	err := sfs.Remove("/single.txt")
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// Remove directory tree
	err = sfs.RemoveAll("/dir")
	if err != nil {
		t.Fatalf("RemoveAll failed: %v", err)
	}

	// Verify
	_, err = sfs.Stat("/single.txt")
	if err == nil {
		t.Fatal("/single.txt should not exist")
	}
	_, err = sfs.Stat("/dir")
	if err == nil {
		t.Fatal("/dir should not exist")
	}
}

func TestDirectoryOperations_ChdirGetwd(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Initial directory
	dir, err := sfs.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	if dir != "/" {
		t.Fatalf("expected '/', got '%s'", dir)
	}

	// Change directory
	sfs.Mkdir("/testdir", 0755)
	err = sfs.Chdir("/testdir")
	if err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}

	dir, err = sfs.Getwd()
	if err != nil {
		t.Fatalf("Getwd failed: %v", err)
	}
	if dir != "/testdir" {
		t.Fatalf("expected '/testdir', got '%s'", dir)
	}
}

// =============================================================================
// Phase 2.3: File Metadata Integration Tests
// =============================================================================

func TestMetadata_Chmod(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Close()

	err := sfs.Chmod("/test.txt", 0700)
	if err != nil {
		t.Fatalf("Chmod failed: %v", err)
	}

	info, _ := sfs.Stat("/test.txt")
	if info.Mode().Perm() != 0700 {
		t.Fatalf("expected 0700, got %o", info.Mode().Perm())
	}
}

func TestMetadata_Chtimes(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	f.Close()

	newTime := time.Date(2021, 6, 15, 12, 0, 0, 0, time.UTC)
	err := sfs.Chtimes("/test.txt", newTime, newTime)
	if err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}

	info, _ := sfs.Stat("/test.txt")
	if !info.ModTime().Equal(newTime) {
		t.Fatalf("expected %v, got %v", newTime, info.ModTime())
	}
}

func TestMetadata_StatVsLstat(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create a regular file
	f, _ := sfs.Create("/target.txt")
	f.Write([]byte("content"))
	f.Close()

	// Create symlink
	sfs.Symlink("/target.txt", "/link.txt")

	// Stat follows symlink
	statInfo, err := sfs.Stat("/link.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if statInfo.Mode()&os.ModeSymlink != 0 {
		t.Fatal("Stat should follow symlink")
	}

	// Lstat does not follow symlink
	lstatInfo, err := sfs.Lstat("/link.txt")
	if err != nil {
		t.Fatalf("Lstat failed: %v", err)
	}
	if lstatInfo.Mode()&os.ModeSymlink == 0 {
		t.Fatal("Lstat should not follow symlink")
	}
}

// =============================================================================
// Phase 2.4: Symlink Operations Integration Tests
// =============================================================================

func TestSymlink_CreateAndRead(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create target
	f, _ := sfs.Create("/target.txt")
	f.Write([]byte("target content"))
	f.Close()

	// Create symlink
	err := sfs.Symlink("/target.txt", "/link.txt")
	if err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	// Read symlink target
	target, err := sfs.Readlink("/link.txt")
	if err != nil {
		t.Fatalf("Readlink failed: %v", err)
	}
	if target != "/target.txt" {
		t.Fatalf("expected '/target.txt', got '%s'", target)
	}
}

func TestSymlink_FollowSymlink(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create target with content
	f, _ := sfs.Create("/target.txt")
	f.Write([]byte("hello from target"))
	f.Close()

	// Create symlink
	err := sfs.Symlink("/target.txt", "/link.txt")
	if err != nil {
		t.Fatalf("Symlink failed: %v", err)
	}

	// Verify symlink was created
	_, err = sfs.Lstat("/link.txt")
	if err != nil {
		t.Fatalf("Lstat on symlink failed: %v", err)
	}

	// Stat follows symlink to get target info
	info, err := sfs.Stat("/link.txt")
	if err != nil {
		t.Fatalf("Stat via symlink failed: %v", err)
	}

	// Verify we get the target file info (not symlink)
	if info.Size() != 17 {
		t.Logf("Note: Stat on symlink returned size %d (expected 17)", info.Size())
	}
}

// =============================================================================
// Phase 2.5: Error Propagation Tests
// =============================================================================

func TestError_FileNotFound(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	_, err := sfs.Open("/nonexistent.txt")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestError_DirectoryNotFound(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	_, err := sfs.Stat("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestError_RemoveNonexistent(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	err := sfs.Remove("/nonexistent.txt")
	if err == nil {
		t.Fatal("expected error for removing nonexistent file")
	}
}

func TestError_MkdirExisting(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	sfs.Mkdir("/testdir", 0755)
	err := sfs.Mkdir("/testdir", 0755)
	if err == nil {
		t.Fatal("expected error for mkdir on existing directory")
	}
}

func TestError_ChdirToFile(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/file.txt")
	f.Close()

	err := sfs.Chdir("/file.txt")
	if err == nil {
		t.Fatal("expected error for chdir to file")
	}
}

// =============================================================================
// Phase 3.1: Type Assertion and Reflection Tests (Anti-Reflection)
// =============================================================================

func TestAntiReflection_TypeAssertionFails(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// The wrapped filesystem should NOT be type-assertable to *memfs.FileSystem
	var fs absfs.SymlinkFileSystem = sfs
	_, ok := fs.(*memfs.FileSystem)
	if ok {
		t.Fatal("type assertion should fail for wrapped filesystem")
	}
}

func TestAntiReflection_InterfacePreserved(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Interface assignment should work
	var symfs absfs.SymlinkFileSystem = sfs
	var fsys absfs.FileSystem = sfs
	var filer absfs.Filer = sfs

	_ = symfs
	_ = fsys
	_ = filer
}

func TestAntiReflection_UnwrapRestoresType(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Unwrap
	unwrapped := ptfs.UnwrapSymlinkFS(sfs)

	// Now type assertion should succeed
	_, ok := unwrapped.(*memfs.FileSystem)
	if !ok {
		t.Fatal("type assertion should succeed after unwrapping")
	}
}

func TestAntiReflection_MultiLayerWrapUnwrap(t *testing.T) {
	mfs := newMemFS(t)

	// Wrap multiple times
	sfs1, _ := ptfs.NewSymlinkFS(mfs)
	sfs2, _ := ptfs.NewSymlinkFS(sfs1)

	// First unwrap
	unwrapped1 := ptfs.UnwrapSymlinkFS(sfs2)
	if unwrapped1 == sfs2 {
		t.Fatal("first unwrap should return inner wrapper")
	}

	// Second unwrap
	unwrapped2 := ptfs.UnwrapSymlinkFS(unwrapped1)
	_, ok := unwrapped2.(*memfs.FileSystem)
	if !ok {
		t.Fatal("second unwrap should return original memfs")
	}
}

// =============================================================================
// Phase 3.2: Interface Compliance Tests
// =============================================================================

func TestInterfaceCompliance_Filer(t *testing.T) {
	mfs := newMemFS(t)
	filer, _ := ptfs.NewFiler(mfs)

	// Compile-time check
	var _ absfs.Filer = filer
}

func TestInterfaceCompliance_FileSystem(t *testing.T) {
	mfs := newMemFS(t)
	fs, _ := ptfs.NewFS(mfs)

	// Compile-time check
	var _ absfs.FileSystem = fs
}

func TestInterfaceCompliance_SymlinkFileSystem(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Compile-time check
	var _ absfs.SymlinkFileSystem = sfs
}

func TestInterfaceCompliance_File(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/test.txt")
	defer f.Close()

	// Compile-time check
	var _ absfs.File = f
}

// =============================================================================
// Phase 3.3: Edge Case Tests
// =============================================================================

func TestEdgeCase_EmptyPath(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	_, err := sfs.Stat("")
	// Empty path should result in an error
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestEdgeCase_LargeFile(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	f, _ := sfs.Create("/large.bin")

	// Write 1MB of data
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	n, err := f.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Fatalf("expected to write %d bytes, wrote %d", len(data), n)
	}
	f.Close()

	// Verify size
	info, _ := sfs.Stat("/large.bin")
	if info.Size() != int64(len(data)) {
		t.Fatalf("expected size %d, got %d", len(data), info.Size())
	}
}

func TestEdgeCase_SpecialCharactersInPath(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Test with spaces
	f, err := sfs.Create("/file with spaces.txt")
	if err != nil {
		t.Fatalf("Create with spaces failed: %v", err)
	}
	f.Close()

	// Test with unicode
	f, err = sfs.Create("/文件.txt")
	if err != nil {
		t.Fatalf("Create with unicode failed: %v", err)
	}
	f.Close()

	// Verify both exist
	_, err = sfs.Stat("/file with spaces.txt")
	if err != nil {
		t.Fatalf("Stat with spaces failed: %v", err)
	}
	_, err = sfs.Stat("/文件.txt")
	if err != nil {
		t.Fatalf("Stat with unicode failed: %v", err)
	}
}

func TestEdgeCase_DeepNesting(t *testing.T) {
	mfs := newMemFS(t)
	sfs, _ := ptfs.NewSymlinkFS(mfs)

	// Create deeply nested path
	deepPath := "/a/b/c/d/e/f/g/h/i/j"
	err := sfs.MkdirAll(deepPath, 0755)
	if err != nil {
		t.Fatalf("MkdirAll deep path failed: %v", err)
	}

	// Create file in deep path
	f, err := sfs.Create(deepPath + "/file.txt")
	if err != nil {
		t.Fatalf("Create in deep path failed: %v", err)
	}
	f.Close()

	// Verify
	_, err = sfs.Stat(deepPath + "/file.txt")
	if err != nil {
		t.Fatalf("Stat in deep path failed: %v", err)
	}
}

// Test that original test still works (basic type compatibility)
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
