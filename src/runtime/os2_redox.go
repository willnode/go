// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/abi"
	"internal/runtime/atomic"
	"unsafe"
)

//go:cgo_export_dynamic runtime.end _end
//go:cgo_export_dynamic runtime.etext _etext
//go:cgo_export_dynamic runtime.edata _edata

// TODO: this is a stub
//go:cgo_export_dynamic _init _init
//go:cgo_export_dynamic _fini _fini

//go:cgo_import_dynamic libc_relibc_start_v1 relibc_start_v1 "libc.so.6"
//go:cgo_import_dynamic libc_getrandom getrandom "libc.so.6"
//go:cgo_import_dynamic libc____errno __errno_location "libc.so.6"
//go:cgo_import_dynamic libc_clock_gettime clock_gettime "libc.so.6"
//go:cgo_import_dynamic libc_exit _exit "libc.so.6"
//go:cgo_import_dynamic libc_kill kill "libc.so.6"
//go:cgo_import_dynamic libc_madvise madvise "libc.so.6"
//go:cgo_import_dynamic libc_malloc malloc "libc.so.6"
//go:cgo_import_dynamic libc_mmap mmap "libc.so.6"
//go:cgo_import_dynamic libc_munmap munmap "libc.so.6"
//go:cgo_import_dynamic libc_open open "libc.so.6"
//go:cgo_import_dynamic libc_pthread_attr_destroy pthread_attr_destroy "libc.so.6"
//go:cgo_import_dynamic libc_pthread_attr_getstack pthread_attr_getstack "libc.so.6"
//go:cgo_import_dynamic libc_pthread_attr_init pthread_attr_init "libc.so.6"
//go:cgo_import_dynamic libc_pthread_attr_setdetachstate pthread_attr_setdetachstate "libc.so.6"
//go:cgo_import_dynamic libc_pthread_attr_setstack pthread_attr_setstack "libc.so.6"
//go:cgo_import_dynamic libc_pthread_attr_setstacksize pthread_attr_setstacksize "libc.so.6"
//go:cgo_import_dynamic libc_pthread_create pthread_create "libc.so.6"
//go:cgo_import_dynamic libc_pthread_self pthread_self "libc.so.6"
//go:cgo_import_dynamic libc_pthread_kill pthread_kill "libc.so.6"
//go:cgo_import_dynamic libc_raise raise "libc.so.6"
//go:cgo_import_dynamic libc_read read "libc.so.6"
//go:cgo_import_dynamic libc_select select "libc.so.6"
//go:cgo_import_dynamic libc_sched_yield sched_yield "libc.so.6"
//go:cgo_import_dynamic libc_sem_init sem_init "libc.so.6"
//go:cgo_import_dynamic libc_sem_post sem_post "libc.so.6"
//go:cgo_import_dynamic libc_sem_timedwait sem_timedwait "libc.so.6"
//go:cgo_import_dynamic libc_sem_wait sem_wait "libc.so.6"
//go:cgo_import_dynamic libc_setitimer setitimer "libc.so.6"
//go:cgo_import_dynamic libc_sigaction sigaction "libc.so.6"
//go:cgo_import_dynamic libc_sigaltstack sigaltstack "libc.so.6"
//go:cgo_import_dynamic libc_sigprocmask sigprocmask "libc.so.6"
//go:cgo_import_dynamic libc_sysconf sysconf "libc.so.6"
//go:cgo_import_dynamic libc_usleep usleep "libc.so.6"
//go:cgo_import_dynamic libc_write write "libc.so.6"
//go:cgo_import_dynamic libc_pipe2 pipe2 "libc.so.6"
//go:cgo_import_dynamic libc_redox_futex_wait_v0 redox_futex_wait_v0 "libc.so.6"
//go:cgo_import_dynamic libc_redox_futex_wake_v0 redox_futex_wake_v0 "libc.so.6"

//go:linkname libc_relibc_start_v1 libc_relibc_start_v1
//go:linkname libc_getrandom libc_getrandom
//go:linkname libc____errno libc____errno
//go:linkname libc_clock_gettime libc_clock_gettime
//go:linkname libc_exit libc_exit
//go:linkname libc_kill libc_kill
//go:linkname libc_madvise libc_madvise
//go:linkname libc_malloc libc_malloc
//go:linkname libc_mmap libc_mmap
//go:linkname libc_munmap libc_munmap
//go:linkname libc_open libc_open
//go:linkname libc_pthread_attr_destroy libc_pthread_attr_destroy
//go:linkname libc_pthread_attr_getstack libc_pthread_attr_getstack
//go:linkname libc_pthread_attr_init libc_pthread_attr_init
//go:linkname libc_pthread_attr_setdetachstate libc_pthread_attr_setdetachstate
//go:linkname libc_pthread_attr_setstack libc_pthread_attr_setstack
//go:linkname libc_pthread_attr_setstacksize libc_pthread_attr_setstacksize
//go:linkname libc_pthread_create libc_pthread_create
//go:linkname libc_pthread_self libc_pthread_self
//go:linkname libc_pthread_kill libc_pthread_kill
//go:linkname libc_raise libc_raise
//go:linkname libc_read libc_read
//go:linkname libc_select libc_select
//go:linkname libc_sched_yield libc_sched_yield
//go:linkname libc_sem_init libc_sem_init
//go:linkname libc_sem_post libc_sem_post
//go:linkname libc_sem_timedwait libc_sem_timedwait
//go:linkname libc_sem_wait libc_sem_wait
//go:linkname libc_setitimer libc_setitimer
//go:linkname libc_sigaction libc_sigaction
//go:linkname libc_sigaltstack libc_sigaltstack
//go:linkname libc_sigprocmask libc_sigprocmask
//go:linkname libc_sysconf libc_sysconf
//go:linkname libc_usleep libc_usleep
//go:linkname libc_write libc_write
//go:linkname libc_pipe2 libc_pipe2
//go:linkname libc_redox_futex_wait_v0 libc_redox_futex_wait_v0
//go:linkname libc_redox_futex_wake_v0 libc_redox_futex_wake_v0

var (
	libc_relibc_start_v1,
	libc_getrandom,
	libc____errno,
	libc_clock_gettime,
	libc_exit,
	libc_kill,
	libc_madvise,
	libc_malloc,
	libc_mmap,
	libc_munmap,
	libc_open,
	libc_pthread_attr_destroy,
	libc_pthread_attr_getstack,
	libc_pthread_attr_init,
	libc_pthread_attr_setdetachstate,
	libc_pthread_attr_setstack,
	libc_pthread_attr_setstacksize,
	libc_pthread_create,
	libc_pthread_self,
	libc_pthread_kill,
	libc_raise,
	libc_read,
	libc_sched_yield,
	libc_select,
	libc_sem_init,
	libc_sem_post,
	libc_sem_timedwait,
	libc_sem_wait,
	libc_setitimer,
	libc_sigaction,
	libc_sigaltstack,
	libc_sigprocmask,
	libc_sysconf,
	libc_usleep,
	libc_write,
	libc_pipe2,
	libc_redox_futex_wait_v0,
	libc_redox_futex_wake_v0 libcFunc
)

var sigset_all = ^sigset(0)

//go:nosplit
func futexsleep(addr *uint32, val uint32, ns int64) {
	var tsPtr uintptr
	if ns >= 0 {
		var ts timespec
		ts.tv_sec = ns / 1000000000
		ts.tv_nsec = int64(ns % 1000000000)
		tsPtr = uintptr(unsafe.Pointer(&ts))
	}

	sysvicall6(&libc_redox_futex_wait_v0, uintptr(unsafe.Pointer(addr)), uintptr(val), tsPtr, 0, 0, 0)
}

//go:nosplit
func futexwakeup(addr *uint32, cnt uint32) {
	sysvicall6(&libc_redox_futex_wake_v0, uintptr(unsafe.Pointer(addr)), uintptr(cnt), 0, 0, 0, 0)
}

func getCPUCount() int32 {
	n := int32(sysconf(_SC_NPROCESSORS_ONLN))
	if n < 1 {
		return 1
	}

	return n
}

func getPageSize() uintptr {
	n := int32(sysconf(_SC_PAGESIZE))
	if n <= 0 {
		return 0
	}
	return uintptr(n)
}


func osinit() {
	asmcgocall(unsafe.Pointer(abi.FuncPCABI0(miniterrno)), unsafe.Pointer(&libc____errno))

	numCPUStartup = getCPUCount()
	if physPageSize == 0 {
		physPageSize = getPageSize()
	}
}

func tstart_sysvicall(newm *m) uint32

// May run with m.p==nil, so write barriers are not allowed.
//
//go:nowritebarrier
func newosproc(mp *m) {
	var (
		attr pthreadattr
		oset sigset
		tid  pthread
		ret  int32
		size uint64
	)
	if pthread_attr_init(&attr) != 0 {
		throw("pthread_attr_init")
	}
	const stackSize = 0x200000
	if pthread_attr_setstacksize(&attr, stackSize) != 0 {
		throw("pthread_attr_setstacksize")
	}
	size = stackSize
	_ = size
	if pthread_attr_setdetachstate(&attr, _PTHREAD_CREATE_DETACHED) != 0 {
		throw("pthread_attr_setdetachstate")
	}

	_ = oset
	ret = retryOnEAGAIN(func() int32 {
		return pthread_create(&tid, &attr, abi.FuncPCABI0(tstart_sysvicall), unsafe.Pointer(mp))
	})
	if ret != 0 {
		print("runtime: failed to create new OS thread (have ", mcount(), " already; errno=", ret, ")\n")
		if ret == _EAGAIN {
			println("runtime: may need to increase max user processes (ulimit -u)")
		}
		throw("newosproc")
	}
}

func exitThread(wait *atomic.Uint32) {
	// should be unreachable
	throw("exitThread")
}

//go:nosplit
func readRandom(r []byte) int {
	n := sysvicall3(&libc_getrandom, uintptr(unsafe.Pointer(&r[0])), uintptr(len(r)), 0)
	return int(n)
}

// populated by _rt0_amd64_redox_main
var redoxEnvp **byte

func goenvs() {
	if redoxEnvp == nil {
		envs = nil
		return
	}
	n := 0
	for argv_index(redoxEnvp, int32(n)) != nil {
		n++
	}
	envs = make([]string, n)
	for i := 0; i < n; i++ {
		envs[i] = gostring(argv_index(redoxEnvp, int32(i)))
	}
}

// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
func mpreinit(mp *m) {
	mp.gsignal = malg(32 * 1024)
	mp.gsignal.m = mp
}

func miniterrno()

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, cannot allocate memory.
func minit() {
	asmcgocall(unsafe.Pointer(abi.FuncPCABI0(miniterrno)), unsafe.Pointer(&libc____errno))

	minitSignals()

	getg().m.procid = uint64(pthread_self())
}

// Called from dropm to undo the effect of an minit.
func unminit() {
	unminitSignals()
	getg().m.procid = 0
}

// Called from mexit, but not from dropm, to undo the effect of thread-owned
// resources in minit, semacreate, or elsewhere. Do not take locks after calling this.
//
// This always runs without a P, so //go:nowritebarrierrec is required.
//
//go:nowritebarrierrec
func mdestroy(mp *m) {
}

func sigtramp()

//go:nosplit
//go:nowritebarrierrec
func setsig(i uint32, fn uintptr) {
	var sa sigactiont

	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTART
	sa.sa_mask = uint64(sigset_all)
	if fn == abi.FuncPCABIInternal(sighandler) {
		fn = abi.FuncPCABI0(sigtramp)
	}
	*((*uintptr)(unsafe.Pointer(&sa.sa_handler))) = fn
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func setsigstack(i uint32) {
	var sa sigactiont
	sigaction(i, nil, &sa)
	if sa.sa_flags&_SA_ONSTACK != 0 {
		return
	}
	sa.sa_flags |= _SA_ONSTACK
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func getsig(i uint32) uintptr {
	var sa sigactiont
	sigaction(i, nil, &sa)
	return *((*uintptr)(unsafe.Pointer(&sa.sa_handler)))
}

// setSignalstackSP sets the ss_sp field of a stackt.
//
//go:nosplit
func setSignalstackSP(s *stackt, sp uintptr) {
	*(*uintptr)(unsafe.Pointer(&s.ss_sp)) = sp
}

//go:nosplit
//go:nowritebarrierrec
func sigaddset(mask *sigset, i int) {
	*mask |= 1 << uint(i-1)
}

//go:nosplit
//go:nowritebarrierrec
func sigdelset(mask *sigset, i int) {
	*mask &^= 1 << uint(i-1)
}

//go:nosplit
func (c *sigctxt) fixsigcode(sig uint32) {
}

func setProcessCPUProfiler(hz int32) {
	setProcessCPUProfilerTimer(hz)
}

func setThreadCPUProfiler(hz int32) {
	setThreadCPUProfilerHz(hz)
}

//go:nosplit
func validSIGPROF(mp *m, c *sigctxt) bool {
	return true
}

//go:nosplit
func closefd(fd int32) int32 {
	return int32(sysvicall1(&libc_close, uintptr(fd)))
}

//go:nosplit
func exit(r int32) {
	sysvicall1(&libc_exit, uintptr(r))
}


//go:nosplit
func madvise(addr unsafe.Pointer, n uintptr, flags int32) {
	// not implemented
	// sysvicall3(&libc_madvise, uintptr(addr), uintptr(n), uintptr(flags))
}

//go:nosplit
func mmap(addr unsafe.Pointer, n uintptr, prot, flags, fd int32, off uint32) (unsafe.Pointer, int) {
	p, err := doMmap(uintptr(addr), n, uintptr(prot), uintptr(flags), uintptr(fd), uintptr(off))
	if p == ^uintptr(0) {
		return nil, int(err)
	}
	return unsafe.Pointer(p), 0
}

//go:nosplit
//go:cgo_unsafe_args
func doMmap(addr, n, prot, flags, fd, off uintptr) (uintptr, uintptr) {
	var libcall libcall
	libcall.fn = uintptr(unsafe.Pointer(&libc_mmap))
	libcall.n = 6
	libcall.args = uintptr(noescape(unsafe.Pointer(&addr)))
	asmcgocall(unsafe.Pointer(&asmsysvicall6x), unsafe.Pointer(&libcall))
	return libcall.r1, libcall.err
}

//go:nosplit
func munmap(addr unsafe.Pointer, n uintptr) {
	sysvicall2(&libc_munmap, uintptr(addr), uintptr(n))
}

const (
	_CLOCK_REALTIME  = 1
	_CLOCK_MONOTONIC = 4
)

//go:nosplit
func nanotime1() int64 {
	var ts mts
	sysvicall2(&libc_clock_gettime, _CLOCK_MONOTONIC, uintptr(unsafe.Pointer(&ts)))
	return ts.tv_sec*1e9 + ts.tv_nsec
}

//go:nosplit
func open(path *byte, mode, perm int32) int32 {
	return int32(sysvicall3(&libc_open, uintptr(unsafe.Pointer(path)), uintptr(mode), uintptr(perm)))
}

func pthread_attr_destroy(attr *pthreadattr) int32 {
	return int32(sysvicall1(&libc_pthread_attr_destroy, uintptr(unsafe.Pointer(attr))))
}

func pthread_attr_getstack(attr *pthreadattr, addr unsafe.Pointer, size *uint64) int32 {
	return int32(sysvicall3(&libc_pthread_attr_getstack, uintptr(unsafe.Pointer(attr)), uintptr(addr), uintptr(unsafe.Pointer(size))))
}

func pthread_attr_init(attr *pthreadattr) int32 {
	return int32(sysvicall1(&libc_pthread_attr_init, uintptr(unsafe.Pointer(attr))))
}

func pthread_attr_setdetachstate(attr *pthreadattr, state int32) int32 {
	return int32(sysvicall2(&libc_pthread_attr_setdetachstate, uintptr(unsafe.Pointer(attr)), uintptr(state)))
}

func pthread_attr_setstack(attr *pthreadattr, addr uintptr, size uint64) int32 {
	return int32(sysvicall3(&libc_pthread_attr_setstack, uintptr(unsafe.Pointer(attr)), uintptr(addr), uintptr(size)))
}

func pthread_attr_setstacksize(attr *pthreadattr, size uint64) int32 {
	return int32(sysvicall2(&libc_pthread_attr_setstacksize, uintptr(unsafe.Pointer(attr)), uintptr(size)))
}

func pthread_create(thread *pthread, attr *pthreadattr, fn uintptr, arg unsafe.Pointer) int32 {
	return int32(sysvicall4(&libc_pthread_create, uintptr(unsafe.Pointer(thread)), uintptr(unsafe.Pointer(attr)), uintptr(fn), uintptr(arg)))
}

func pthread_self() pthread {
	return pthread(sysvicall0(&libc_pthread_self))
}

func signalM(mp *m, sig int) {
	sysvicall2(&libc_pthread_kill, uintptr(pthread(mp.procid)), uintptr(sig))
}

//go:nosplit
//go:nowritebarrierrec
func raise(sig uint32) /* int32 */ {
	sysvicall1(&libc_raise, uintptr(sig))
}

func raiseproc(sig uint32) /* int32 */ {
	pid := sysvicall0(&libc_getpid)
	sysvicall2(&libc_kill, pid, uintptr(sig))
}

//go:nosplit
func read(fd int32, buf unsafe.Pointer, nbyte int32) int32 {
	r1, err := sysvicall3Err(&libc_read, uintptr(fd), uintptr(buf), uintptr(nbyte))
	if c := int32(r1); c >= 0 {
		return c
	}
	return -int32(err)
}

//go:nosplit
func sem_init(sem *semt, pshared int32, value uint32) int32 {
	return int32(sysvicall3(&libc_sem_init, uintptr(unsafe.Pointer(sem)), uintptr(pshared), uintptr(value)))
}

//go:nosplit
func sem_post(sem *semt) int32 {
	return int32(sysvicall1(&libc_sem_post, uintptr(unsafe.Pointer(sem))))
}

//go:nosplit
func sem_timedwait(sem *semt, abstime *timespec) int32 {
	return int32(sysvicall2(&libc_sem_timedwait, uintptr(unsafe.Pointer(sem)), uintptr(unsafe.Pointer(abstime))))
}

//go:nosplit
func sem_wait(sem *semt) int32 {
	return int32(sysvicall1(&libc_sem_wait, uintptr(unsafe.Pointer(sem))))
}

//go:nosplit
func setitimer(which int32, value *itimerval, ovalue *itimerval) /* int32 */ {
	sysvicall3(&libc_setitimer, uintptr(which), uintptr(unsafe.Pointer(value)), uintptr(unsafe.Pointer(ovalue)))
}

//go:nosplit
//go:nowritebarrierrec
func sigaction(sig uint32, act *sigactiont, oact *sigactiont) /* int32 */ {
	sysvicall3(&libc_sigaction, uintptr(sig), uintptr(unsafe.Pointer(act)), uintptr(unsafe.Pointer(oact)))
}

//go:nosplit
//go:nowritebarrierrec
func sigaltstack(ss *stackt, oss *stackt) /* int32 */ {
	sysvicall2(&libc_sigaltstack, uintptr(unsafe.Pointer(ss)), uintptr(unsafe.Pointer(oss)))
}

//go:nosplit
//go:nowritebarrierrec
func sigprocmask(how int32, set *sigset, oset *sigset) /* int32 */ {
	sysvicall3(&libc_sigprocmask, uintptr(how), uintptr(unsafe.Pointer(set)), uintptr(unsafe.Pointer(oset)))
}

func sysconf(name int32) int64 {
	return int64(sysvicall1(&libc_sysconf, uintptr(name)))
}

func usleep1(usec uint32)

//go:nosplit
func usleep_no_g(µs uint32) {
	usleep1(µs)
}

//go:nosplit
func usleep(µs uint32) {
	usleep1(µs)
}

func walltime() (sec int64, nsec int32) {
	var ts mts
	sysvicall2(&libc_clock_gettime, _CLOCK_REALTIME, uintptr(unsafe.Pointer(&ts)))
	return ts.tv_sec, int32(ts.tv_nsec)
}

//go:nosplit
func write1(fd uintptr, buf unsafe.Pointer, nbyte int32) int32 {
	r1, err := sysvicall3Err(&libc_write, fd, uintptr(buf), uintptr(nbyte))
	if c := int32(r1); c >= 0 {
		return c
	}
	return -int32(err)
}

//go:nosplit
func pipe2(flags int32) (r, w int32, errno int32) {
	var p [2]int32
	_, e := sysvicall2Err(&libc_pipe2, uintptr(noescape(unsafe.Pointer(&p))), uintptr(flags))
	return p[0], p[1], int32(e)
}

//go:nosplit
func fcntl(fd, cmd, arg int32) (ret int32, errno int32) {
	r1, err := sysvicall3Err(&libc_fcntl, uintptr(fd), uintptr(cmd), uintptr(arg))
	return int32(r1), int32(err)
}

func osyield1()

//go:nosplit
func osyield_no_g() {
	osyield1()
}

//go:nosplit
func osyield() {
	sysvicall0(&libc_sched_yield)
}

//go:linkname executablePath os.executablePath
var executablePath string

func sysargs(argc int32, argv **byte) {
	// there's no auxv
	auxv = nil
}

// sigPerThreadSyscall is only used on linux, so we assign a bogus signal
// number.
const sigPerThreadSyscall = 1 << 31

//go:nosplit
func runPerThreadSyscall() {
	throw("runPerThreadSyscall only valid on linux")
}

//go:nowritebarrierrec
//go:nosplit
func libpreinit() {}

//go:nowritebarrierrec
//go:nosplit
func newosproc0(stacksize uintptr, fn unsafe.Pointer) {
	throw("bad newosproc0")
}
