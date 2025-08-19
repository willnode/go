// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/abi"
	"internal/goarch"
	"internal/runtime/atomic"
	"unsafe"
)


//
// C library function declarations
//

//go:noescape
func setitimer(mode int32, new, old *itimerval)
//go:noescape
func sigaction(sig uint32, new, old *sigactiont)
//go:noescape
func sigaltstack(new, old *stackt)
//go:noescape
func sigprocmask(how int32, new, old *sigset)
//go:noescape
func getcontext(ctxt unsafe.Pointer)
//go:noescape
func osyield()
//go:noescape
func pipe2(flags int32) (r, w int32, errno int32)
//go:noescape
func fcntl(fd, cmd, arg int32) (ret int32, errno int32)
//go:noescape
func sysconf(name int32) int64
//go:noescape
func pthread_attr_init(attr *pthread_attr_t) int32
//go:noescape
func pthread_attr_setdetachstate(attr *pthread_attr_t, state int32) int32
//go:noescape
func pthread_attr_setstack(attr *pthread_attr_t, addr unsafe.Pointer, size uintptr) int32
//go:noescape
func pthread_create(tid *pthread_t, attr *pthread_attr_t, start unsafe.Pointer, arg unsafe.Pointer) int32
//go:noescape
func pthread_self() pthread_t
//go:noescape
func pthread_kill(tid pthread_t, sig int) int32
//go:noescape
func pthread_mutex_init(m *pthread_mutex_t, attr unsafe.Pointer) int32
//go:noescape
func pthread_mutex_lock(m *pthread_mutex_t) int32
//go:noescape
func pthread_mutex_unlock(m *pthread_mutex_t) int32
//go:noescape
func pthread_cond_init(c *pthread_cond_t, attr unsafe.Pointer) int32
//go:noescape
func pthread_cond_wait(c *pthread_cond_t, m *pthread_mutex_t) int32
//go:noescape
func pthread_cond_timedwait(c *pthread_cond_t, m *pthread_mutex_t, ts *timespec) int32
//go:noescape
func pthread_cond_signal(c *pthread_cond_t) int32

// Standard error and clock constants.
const (
	_ESRCH     = 3
	_ETIMEDOUT = 60
	_EAGAIN    = 35

	_CLOCK_REALTIME  = 0
	_CLOCK_MONOTONIC = 3

	_TIMER_RELTIME = 0
	_TIMER_ABSTIME = 1

	_PTHREAD_CREATE_DETACHED = 1

	_SC_NPROCESSORS_ONLN = 58
	_SC_PAGESIZE         = 30
)

var sigset_all = sigset{[4]uint32{^uint32(0), ^uint32(0), ^uint32(0), ^uint32(0)}}

func getCPUCount() int32 {
	n := sysconf(_SC_NPROCESSORS_ONLN)
	if n < 1 {
		return 1
	}
	return int32(n)
}

func getPageSize() uintptr {
	n := sysconf(_SC_PAGESIZE)
	if n <= 0 {
		return 4096 // fallback
	}
	return uintptr(n)
}

//go:nosplit
func osyield_no_g() {
	osyield()
}

//go:nosplit
func semacreate(mp *m) {
	mp.waitsemacount = 0
	pthread_mutex_init(&mp.waitsemamutex, nil)
	pthread_cond_init(&mp.waitsemacond, nil)
}

//go:nosplit
func semasleep(ns int64) int32 {
	gp := getg()
	mp := gp.m

	// Fast path: check count without locking.
	if atomic.Load(&mp.waitsemacount) > 0 {
		if atomic.Cas(&mp.waitsemacount, 1, 0) {
			return 0 // Acquired semaphore.
		}
	}

	pthread_mutex_lock(&mp.waitsemamutex)
	for mp.waitsemacount == 0 {
		if ns < 0 {
			pthread_cond_wait(&mp.waitsemacond, &mp.waitsemamutex)
		} else {
			var ts timespec
			ts.setNsec(nanotime() + ns)
			ret := pthread_cond_timedwait(&mp.waitsemacond, &mp.waitsemamutex, &ts)
			if ret == _ETIMEDOUT {
				pthread_mutex_unlock(&mp.waitsemamutex)
				return -1
			}
		}
	}
	mp.waitsemacount--
	pthread_mutex_unlock(&mp.waitsemamutex)
	return 0
}

//go:nosplit
func semawakeup(mp *m) {
	atomic.Xadd(&mp.waitsemacount, 1)
	pthread_mutex_lock(&mp.waitsemamutex)
	pthread_cond_signal(&mp.waitsemacond)
	pthread_mutex_unlock(&mp.waitsemamutex)
}

// This is the entry point for a new thread created by pthread_create.
// It is written in assembly.
func threadentry()

// May run with m.p==nil, so write barriers are not allowed.
//go:nowritebarrier
func newosproc(mp *m) {
	stk := unsafe.Pointer(mp.g0.stack.hi)
	if false {
		print("newosproc stk=", stk, " m=", mp, " g=", mp.g0, " id=", mp.id, " ostk=", &mp, "\n")
	}

	var attr pthread_attr_t
	if err := pthread_attr_init(&attr); err != 0 {
		throw("pthread_attr_init failed")
	}
	if err := pthread_attr_setstack(&attr, mp.g0.stack.lo, mp.g0.stack.hi-mp.g0.stack.lo); err != 0 {
		throw("pthread_attr_setstack failed")
	}
	if err := pthread_attr_setdetachstate(&attr, _PTHREAD_CREATE_DETACHED); err != 0 {
		throw("pthread_attr_setdetachstate failed")
	}

	var oset sigset
	sigprocmask(_SIG_SETMASK, &sigset_all, &oset)

	var tid pthread_t
	ret := pthread_create(&tid, &attr, abi.FuncPCABI0(threadentry), unsafe.Pointer(mp))

	sigprocmask(_SIG_SETMASK, &oset, nil)

	if ret != 0 {
		print("runtime: failed to create new OS thread (have ", mcount()-1, " already; errno=", ret, ")\n")
		if ret == _EAGAIN {
			println("runtime: may need to increase max user processes (ulimit -u)")
		}
		throw("runtime.newosproc")
	}
}

func osinit() {
	numCPUStartup = getCPUCount()
	if physPageSize == 0 {
		physPageSize = getPageSize()
	}
}

var urandom_dev = []byte("/dev/urandom\x00")

//go:nosplit
func readRandom(r []byte) int {
	fd := open(&urandom_dev[0], 0 /* O_RDONLY */, 0)
	n := read(fd, unsafe.Pointer(&r[0]), int32(len(r)))
	closefd(fd)
	return int(n)
}

func goenvs() {
	goenvs_unix()
}

func mpreinit(mp *m) {
	mp.gsignal = malg(32 * 1024)
	mp.gsignal.m = mp
}

func minit() {
	gp := getg()
	gp.m.procid = uint64(pthread_self())
	signalstack(&gp.m.gsignal.stack)
	gp.m.newSigstack = true
	minitSignalMask()
}

//go:nosplit
func unminit() {
	unminitSignals()
}

//go:nowritebarrierrec
func mdestroy(mp *m) {
}

//
// Signal handling
//

func sigtramp()

type sigactiont struct {
	sa_sigaction uintptr
	sa_mask      sigset
	sa_flags     int32
}

//go:nosplit
//go:nowritebarrierrec
func setsig(i uint32, fn uintptr) {
	var sa sigactiont
	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTART
	sa.sa_mask = sigset_all
	if fn == abi.FuncPCABIInternal(sighandler) {
		fn = abi.FuncPCABI0(sigtramp)
	}
	sa.sa_sigaction = fn
	sigaction(i, &sa, nil)
}

//go:nosplit
//go:nowritebarrierrec
func setsigstack(i uint32) {
	throw("setsigstack")
}

//go:nosplit
//go:nowritebarrierrec
func getsig(i uint32) uintptr {
	var sa sigactiont
	sigaction(i, nil, &sa)
	return sa.sa_sigaction
}

//go:nosplit
func setSignalstackSP(s *stackt, sp uintptr) {
	s.ss_sp = sp
}

//go:nosplit
//go:nowritebarrierrec
func sigaddset(mask *sigset, i int) {
	if i > _NSIG {
		return
	}
	mask.__bits[(i-1)/32] |= 1 << ((uint32(i) - 1) & 31)
}

func sigdelset(mask *sigset, i int) {
	if i > _NSIG {
		return
	}
	mask.__bits[(i-1)/32] &^= 1 << ((uint32(i) - 1) & 31)
}

//go:nosplit
func (c *sigctxt) fixsigcode(sig uint32) {
}

func setProcessCPUProfiler(hz int32) {
	// setProcessCPUProfilerTimer(hz)
}

func setThreadCPUProfiler(hz int32) {
	// setThreadCPUProfilerHz(hz)
}

//go:nosplit
func validSIGPROF(mp *m, c *sigctxt) bool {
	return true
}

// raise sends a signal to the calling thread.
//go:nosplit
func raise(sig uint32) {
	pthread_kill(pthread_self(), int(sig))
}

// signalM sends a signal to a specific Go thread (an M).
func signalM(mp *m, sig int) {
	pthread_kill(pthread_t(mp.procid), sig)
}

// sigPerThreadSyscall is only used on linux.
const sigPerThreadSyscall = 1 << 31

//go:nosplit
func runPerThreadSyscall() {
	throw("runPerThreadSyscall only valid on linux")
}