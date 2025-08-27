// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: AT THIS MOMENT REDOX CAN'T USE ALL OF THIS
// TODO: Should be merged with at_libc when openat is fully implemented

package unix

import (
	"syscall"
)

// Implemented as sysvicall6 in runtime/syscall_redox.go.
func syscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

// Implemented as rawsysvicall6 in runtime/syscall_redox.go.
func rawSyscall6(trap, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err syscall.Errno)

const (
	AT_EACCESS          = 0x4
	AT_FDCWD            = 0xffd19553
	AT_REMOVEDIR        = 0x1
	AT_SYMLINK_NOFOLLOW = 0x1000

	UTIME_OMIT = -0x2
)

func Unlinkat(dirfd int, path string, flags int) error {
	return syscall.ENOSYS
}

func Openat(dirfd int, path string, flags int, perm uint32) (int, error) {
	return -1, syscall.ENOSYS
}

func Fstatat(dirfd int, path string, stat *syscall.Stat_t, flags int) error {
	return syscall.ENOSYS
}

func Readlinkat(dirfd int, path string, buf []byte) (int, error) {
	return -1, syscall.ENOSYS
}

func Mkdirat(dirfd int, path string, mode uint32) error {
	return syscall.ENOSYS
}

func Fchmodat(dirfd int, path string, mode uint32, flags int) error {
	return syscall.ENOSYS
}

func Fchownat(dirfd int, path string, uid, gid int, flags int) error {
	return syscall.ENOSYS
}

func Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) error {
	return syscall.ENOSYS
}

func Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flag int) error {
	return syscall.ENOSYS
}

func Symlinkat(oldpath string, newdirfd int, newpath string) error {
	return syscall.ENOSYS
}

func Eaccess(path string, mode uint32) error {
	return syscall.ENOSYS
}
