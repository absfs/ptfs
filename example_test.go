package ptfs_test

import (
	"fmt"

	"github.com/absfs/absfs"
	"github.com/absfs/memfs"
	"github.com/absfs/ptfs"
)

// Example demonstrating basic pass-through usage.
func Example_basicUsage() {
	// Create a memfs filesystem
	mfs, _ := memfs.NewFS()

	// Wrap it with ptfs
	pfs, _ := ptfs.NewSymlinkFS(mfs)

	// Use the wrapped filesystem - all operations pass through
	f, _ := pfs.Create("/hello.txt")
	f.WriteString("Hello, World!")
	f.Close()

	// Read it back
	f, _ = pfs.Open("/hello.txt")
	buf := make([]byte, 100)
	n, _ := f.Read(buf)
	f.Close()

	fmt.Println(string(buf[:n]))
	// Output: Hello, World!
}

// Example demonstrating the anti-reflection feature.
func Example_antiReflection() {
	// Create a memfs filesystem
	mfs, _ := memfs.NewFS()

	// Without wrapping, type assertion succeeds
	var fs absfs.SymlinkFileSystem = mfs
	_, isMemfs := fs.(*memfs.FileSystem)
	fmt.Printf("Before wrapping: is memfs? %v\n", isMemfs)

	// After wrapping with ptfs, type assertion fails
	wrapped, _ := ptfs.NewSymlinkFS(mfs)
	fs = wrapped
	_, isMemfs = fs.(*memfs.FileSystem)
	fmt.Printf("After wrapping: is memfs? %v\n", isMemfs)

	// Output:
	// Before wrapping: is memfs? true
	// After wrapping: is memfs? false
}

// Example demonstrating the unwrap functionality.
func Example_unwrap() {
	// Create and wrap a memfs
	mfs, _ := memfs.NewFS()
	wrapped, _ := ptfs.NewSymlinkFS(mfs)

	// Type assertion fails on wrapped
	var fs absfs.SymlinkFileSystem = wrapped
	_, ok := fs.(*memfs.FileSystem)
	fmt.Printf("Wrapped: type assertion succeeded? %v\n", ok)

	// Unwrap to get original back
	unwrapped := ptfs.UnwrapSymlinkFS(wrapped)
	_, ok = unwrapped.(*memfs.FileSystem)
	fmt.Printf("Unwrapped: type assertion succeeded? %v\n", ok)

	// Output:
	// Wrapped: type assertion succeeded? false
	// Unwrapped: type assertion succeeded? true
}

// Example demonstrating multi-layer wrapping.
func Example_multiLayerWrap() {
	mfs, _ := memfs.NewFS()

	// Wrap multiple times
	layer1, _ := ptfs.NewSymlinkFS(mfs)
	layer2, _ := ptfs.NewSymlinkFS(layer1)

	// Operations still work through multiple layers
	f, _ := layer2.Create("/test.txt")
	f.WriteString("Multi-layer test")
	f.Close()

	info, _ := layer2.Stat("/test.txt")
	fmt.Printf("File size: %d\n", info.Size())

	// Unwrap layer by layer
	unwrap1 := ptfs.UnwrapSymlinkFS(layer2)
	_, isFirstLayer := unwrap1.(*ptfs.SymlinkFileSystem)
	fmt.Printf("First unwrap is ptfs? %v\n", isFirstLayer)

	unwrap2 := ptfs.UnwrapSymlinkFS(unwrap1)
	_, isMemfs := unwrap2.(*memfs.FileSystem)
	fmt.Printf("Second unwrap is memfs? %v\n", isMemfs)

	// Output:
	// File size: 16
	// First unwrap is ptfs? true
	// Second unwrap is memfs? true
}
