// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

//go:cgo_import_dynamic libc_chdir chdir "libc.so"
//go:cgo_import_dynamic libc_chroot chroot "libc.so"
//go:cgo_import_dynamic libc_close close "libc.so"
//go:cgo_import_dynamic libc_dup2 execve "libc.so"
//go:cgo_import_dynamic libc_execve execve "libc.so"
//go:cgo_import_dynamic libc_fcntl fcntl "libc.so"
//go:cgo_import_dynamic libc_gethostname gethostname "libc.so"
//go:cgo_import_dynamic libc_getpid getpid "libc.so"
//go:cgo_import_dynamic libc_ioctl ioctl "libc.so"
//go:cgo_import_dynamic libc_setgid setgid "libc.so"
//go:cgo_import_dynamic libc_setgroups setgroups "libc.so"
//go:cgo_import_dynamic libc_setrlimit setrlimit "libc.so"
//go:cgo_import_dynamic libc_setsid setsid "libc.so"
//go:cgo_import_dynamic libc_setuid setuid "libc.so"
//go:cgo_import_dynamic libc_setpgid setpgid "libc.so"
//go:cgo_import_dynamic libc_issetugid issetugid "libc.so"

//go:linkname libc_chdir libc_chdir
//go:linkname libc_chroot libc_chroot
//go:linkname libc_close libc_close
//go:linkname libc_dup2 libc_dup2
//go:linkname libc_execve libc_execve
//go:linkname libc_fcntl libc_fcntl
//go:linkname libc_gethostname libc_gethostname
//go:linkname libc_getpid libc_getpid
//go:linkname libc_ioctl libc_ioctl
//go:linkname libc_setgid libc_setgid
//go:linkname libc_setgroups libc_setgroups
//go:linkname libc_setrlimit libc_setrlimit
//go:linkname libc_setsid libc_setsid
//go:linkname libc_setuid libc_setuid
//go:linkname libc_setpgid libc_setpgid
//go:linkname libc_issetugid libc_issetugid

var (
	libc_chdir,
	libc_chroot,
	libc_close,
	libc_dup2,
	libc_execve,
	libc_fcntl,
	libc_gethostname,
	libc_getpid,
	libc_ioctl,
	libc_setgid,
	libc_setgroups,
	libc_setrlimit,
	libc_setsid,
	libc_setuid,
	libc_setpgid,
	libc_issetugid libcFunc
)

// Many of these are exported via linkname to assembly in the syscall
// package.

//go:nosplit
//go:linkname syscall_sysvicall6
//go:cgo_unsafe_args
func syscall_sysvicall6(fn, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	call := libcall{
		fn:   fn,
		n:    nargs,
		args: uintptr(unsafe.Pointer(&a1)),
	}
	entersyscallblock()
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	exitsyscall()
	return call.r1, call.r2, call.err
}

//go:nosplit
//go:linkname syscall_rawsysvicall6
//go:cgo_unsafe_args
func syscall_rawsysvicall6(fn, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	call := libcall{
		fn:   fn,
		n:    nargs,
		args: uintptr(unsafe.Pointer(&a1)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, call.r2, call.err
}

// TODO(aram): Once we remove all instances of C calling sysvicallN, make
// sysvicallN return errors and replace the body of the following functions
// with calls to sysvicallN.

//go:nosplit
//go:linkname syscall_chdir
func syscall_chdir(path uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_chdir)),
		n:    1,
		args: uintptr(unsafe.Pointer(&path)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

//go:nosplit
//go:linkname syscall_chroot
func syscall_chroot(path uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_chroot)),
		n:    1,
		args: uintptr(unsafe.Pointer(&path)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

// like close, but must not split stack, for forkx.
//
//go:nosplit
//go:linkname syscall_close
func syscall_close(fd int32) int32 {
	return int32(sysvicall1(&libc_close, uintptr(fd)))
}

const _F_DUP2FD = 0x9

//go:nosplit
//go:linkname syscall_dup2
func syscall_dup2(oldfd, newfd uintptr) (val, err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_dup2)),
		n:    3,
		args: uintptr(unsafe.Pointer(&oldfd)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, call.err
}

//go:nosplit
//go:linkname syscall_execve
//go:cgo_unsafe_args
func syscall_execve(path, argv, envp uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_execve)),
		n:    3,
		args: uintptr(unsafe.Pointer(&path)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

// like exit, but must not split stack, for forkx.
//
//go:nosplit
//go:linkname syscall_exit
func syscall_exit(code uintptr) {
	sysvicall1(&libc_exit, code)
}

//go:nosplit
//go:linkname syscall_fcntl
//go:cgo_unsafe_args
func syscall_fcntl(fd, cmd, arg uintptr) (val, err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_fcntl)),
		n:    3,
		args: uintptr(unsafe.Pointer(&fd)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, call.err
}

//go:nosplit
//go:linkname syscall_forkx
func syscall_forkx(flags uintptr) (pid uintptr, err uintptr) {
	panic("forkx TODO")
}

//go:linkname syscall_gethostname
func syscall_gethostname() (name string, err uintptr) {
	cname := new([_MAXHOSTNAMELEN]byte)
	var args = [2]uintptr{uintptr(unsafe.Pointer(&cname[0])), _MAXHOSTNAMELEN}
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_gethostname)),
		n:    2,
		args: uintptr(unsafe.Pointer(&args[0])),
	}
	entersyscallblock()
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	exitsyscall()
	if call.r1 != 0 {
		return "", call.err
	}
	cname[_MAXHOSTNAMELEN-1] = 0
	return gostringnocopy(&cname[0]), 0
}

//go:nosplit
//go:linkname syscall_getpid
func syscall_getpid() (pid, err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_getpid)),
		n:    0,
		args: uintptr(unsafe.Pointer(&libc_getpid)), // it's unused but must be non-nil, otherwise crashes
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, call.err
}

//go:nosplit
//go:linkname syscall_ioctl
//go:cgo_unsafe_args
func syscall_ioctl(fd, req, arg uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_ioctl)),
		n:    3,
		args: uintptr(unsafe.Pointer(&fd)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

// This is syscall.RawSyscall, it exists to satisfy some build dependency,
// but it doesn't work.
//
//go:linkname syscall_rawsyscall
func syscall_rawsyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	panic("RawSyscall not available on Solaris")
}

// This is syscall.RawSyscall6, it exists to avoid a linker error because
// syscall.RawSyscall6 is already declared. See golang.org/issue/24357
//
//go:linkname syscall_rawsyscall6
func syscall_rawsyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	panic("RawSyscall6 not available on Solaris")
}

//go:nosplit
//go:linkname syscall_setgid
func syscall_setgid(gid uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_setgid)),
		n:    1,
		args: uintptr(unsafe.Pointer(&gid)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

//go:nosplit
//go:linkname syscall_setgroups
//go:cgo_unsafe_args
func syscall_setgroups(ngid, gid uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_setgroups)),
		n:    2,
		args: uintptr(unsafe.Pointer(&ngid)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

//go:nosplit
//go:linkname syscall_setrlimit
//go:cgo_unsafe_args
func syscall_setrlimit(which uintptr, lim unsafe.Pointer) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_setrlimit)),
		n:    2,
		args: uintptr(unsafe.Pointer(&which)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

//go:nosplit
//go:linkname syscall_setsid
func syscall_setsid() (pid, err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_setsid)),
		n:    0,
		args: uintptr(unsafe.Pointer(&libc_setsid)), // it's unused but must be non-nil, otherwise crashes
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, call.err
}

//go:nosplit
//go:linkname syscall_setuid
func syscall_setuid(uid uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_setuid)),
		n:    1,
		args: uintptr(unsafe.Pointer(&uid)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

//go:nosplit
//go:linkname syscall_setpgid
//go:cgo_unsafe_args
func syscall_setpgid(pid, pgid uintptr) (err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_setpgid)),
		n:    2,
		args: uintptr(unsafe.Pointer(&pid)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.err
}

//go:linkname syscall_syscall
//go:cgo_unsafe_args
func syscall_syscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	panic("syscall TODO")
}

//go:linkname syscall_wait4
//go:cgo_unsafe_args
func syscall_wait4(pid uintptr, wstatus *uint32, options uintptr, rusage unsafe.Pointer) (wpid int, err uintptr) {
	panic("wait4 TODO")
}

//go:nosplit
//go:linkname syscall_write
//go:cgo_unsafe_args
func syscall_write(fd, buf, nbyte uintptr) (n, err uintptr) {
	call := libcall{
		fn:   uintptr(unsafe.Pointer(&libc_write)),
		n:    3,
		args: uintptr(unsafe.Pointer(&fd)),
	}
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&call))
	return call.r1, call.err
}
