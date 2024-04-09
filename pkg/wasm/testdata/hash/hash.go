package main

// #include <stdlib.h>
import "C"

import (
	"crypto/sha1"
	"encoding/base64"
	"unsafe"
)

func hash(in []byte) []byte {
	h := sha1.New()

	_, _ = h.Write(in)

	return []byte(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

// _hash is a WebAssembly export that accepts a string pointer (linear memory
// offset) and returns a pointer/size pair packed into a uint64.
//
// Note: This uses a uint64 instead of two result values for compatibility with
// WebAssembly 1.0.
//
//export hash
func _hash(ptr, size uint32) (ptrSize uint64) {
	name := ptrToString(ptr, size)
	g := hash([]byte(name))
	ptr, size = stringToLeakedPtr(string(g))
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

// ptrToString returns a string from WebAssembly compatible numeric types
// representing its pointer and length.
func ptrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

// stringToLeakedPtr returns a pointer and size pair for the given string in a way
// compatible with WebAssembly numeric types.
// The pointer is not automatically managed by TinyGo hence it must be freed by the host.
func stringToLeakedPtr(s string) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
