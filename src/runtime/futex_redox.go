// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package runtime

import "unsafe"

// Redox have direct syscall to futex but no libc equivalent
const (
	_SYS_FUTEX  = 240
	_FUTEX_WAIT = 0
	_FUTEX_WAKE = 1
)

//go:nosplit
func futexsleep(addr *uint32, val uint32, ns int64) {
	var tsPtr uintptr
	if ns >= 0 {
		var ts timespec
		ts.tv_sec = ns / 1000000000
		ts.tv_nsec = int64(ns % 1000000000)
		tsPtr = uintptr(unsafe.Pointer(&ts))
	}
	rawSyscall6(_SYS_FUTEX, uintptr(unsafe.Pointer(addr)),
		_FUTEX_WAIT, uintptr(val), tsPtr, 0, 0)
}

//go:nosplit
func futexwakeup(addr *uint32, cnt uint32) {
	rawSyscall6(_SYS_FUTEX, uintptr(unsafe.Pointer(addr)),
		_FUTEX_WAKE, uintptr(cnt), 0, 0, 0)
}
