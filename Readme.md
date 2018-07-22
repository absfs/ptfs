# ptfs - Pass Through File System
The `ptfs` package implements the `absfs` filesystem interfaces. `ptfs` wraps an object that implements one of the `absfs` interfaces, and passes all interfaces methods to the underlying type without modification, returning the unmodified results to the caller. 

This is useful as a template for creating new `absfs` types, but it serves another subtle but important purpose.

## Anti-Reflective Coating
Wrapping a `absfs` compatible object in a `psfs` type forces the object to be the `absfs` interface type exclusively.

The pass through filesystem is useful in contexts in which you need to obscure a type from reflection or type assertion. For example if a package produces and consumes `absfs.FileSystem` objects, and recognizes it's own `FileSystem` objects, wrapping those objects in a `psfs.FileSystem` type will prevent recognition. Type assertion is an important and useful feature of Go, and typically shouldn't be circumvented, but in special cases such as testing, it may be useful to use the pass through filesystem to shield against reflection. An anti-reflection coating if you will.

The `ptfs` package implements the `absfs.FileSystem`, `absfs.Filer`, and `absfs.SymlinkFileSystem` interfaces. Additionally `ptfs` has `Unwrap` functions that can remove the pass through filesystem wrapper and return the underlying `absfs` interface type.

## Install

```bash
$ go get github.com/absfs/ptfs
```

## Example

```go
package main

import (
    "fmt"
    "reflect"

    "github.com/absfs/absfs"
    "github.com/absfs/ptfs"
)

type myfs struct{
    //...
}

// myfs implements `absfs.FileSystem`
// Plus...

func (fs *myfs) String() string {
    return "It is a `myfs` type."
}

func FsStringer(fs absfs.FileSystem) string {
    // myfs is an absfs.FileSystem, and can be cast to a Stringer.
    // ptfs.FileSystem is an absfs.FileSystem, but cannot be cast to a Stringer.
    s, ok := fs.(fmt.Stringer)
    if ok {
        return s.String()
    }
    return "It doesn't look like anything to me."
}

func main() {
    var fs absfs.FileSystem
    fs = &myfs{}
    
    fmt.Println(FsStringer(fs)) // output: "It is a `myfs` type."
    fs, _ = ptfs.NewFs(fs)
    fmt.Println(FsStringer(fs)) // output: "It doesn't look like anything to me."
}
```

## absfs
Check out the [`absfs`](https://github.com/absfs/absfs) repo for more information about the abstract FileSystem interface and features like FileSystem composition.

## LICENSE

This project is governed by the MIT License. See [LICENSE](https://github.com/absfs/osfs/blob/master/LICENSE)



