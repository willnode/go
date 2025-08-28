package runtime

import (
	"unsafe"
)

const (
	_si_max_size    = 128
	_sigev_max_size = 64
)

const (
	_O_RDONLY   = 0x0
	_O_WRONLY   = 0x1
	_O_CREAT    = 0x40
	_O_TRUNC    = 0x200
	_O_NONBLOCK = 0x800
	_O_CLOEXEC  = 0x80000
)

// Standard signal-related constants.
const (
	_SS_DISABLE  = 4
	_SIG_BLOCK   = 1
	_SIG_UNBLOCK = 2
	_SIG_SETMASK = 3
	_NSIG        = 128 // Increased to match sigset size.
	_SI_USER     = 0
	_UC_SIGMASK  = 0x01
	_UC_CPU      = 0x04

	_MADV_DONTNEED   = 0x4
	_MADV_FREE       = 0x8
	_MADV_HUGEPAGE   = 0xe
	_MADV_NOHUGEPAGE = 0xf
	_MADV_COLLAPSE   = 0x19

	_FPE_INTDIV = 0x1
	_FPE_INTOVF = 0x2
	_FPE_FLTDIV = 0x3
	_FPE_FLTOVF = 0x4
	_FPE_FLTUND = 0x5
	_FPE_FLTRES = 0x6
	_FPE_FLTINV = 0x7
	_FPE_FLTSUB = 0x8

	_SIGHUP    = 0x1
	_SIGINT    = 0x2
	_SIGQUIT   = 0x3
	_SIGILL    = 0x4
	_SIGTRAP   = 0x5
	_SIGABRT   = 0x6
	_SIGEMT    = 0x7
	_SIGFPE    = 0x8
	_SIGKILL   = 0x9
	_SIGBUS    = 0xa
	_SIGSEGV   = 0xb
	_SIGSYS    = 0xc
	_SIGPIPE   = 0xd
	_SIGALRM   = 0xe
	_SIGTERM   = 0xf
	_SIGURG    = 0x10
	_SIGSTOP   = 0x11
	_SIGTSTP   = 0x12
	_SIGCONT   = 0x13
	_SIGCHLD   = 0x14
	_SIGTTIN   = 0x15
	_SIGTTOU   = 0x16
	_SIGIO     = 0x17
	_SIGXCPU   = 0x18
	_SIGXFSZ   = 0x19
	_SIGVTALRM = 0x1a
	_SIGPROF   = 0x1b
	_SIGWINCH  = 0x1c
	_SIGINFO   = 0x1d
	_SIGUSR1   = 0x1e
	_SIGUSR2   = 0x1f

	_SA_RESTART  = 0x10000000
	_SA_ONSTACK  = 0x8000000
	_SA_RESTORER = 0x4000000
	_SA_SIGINFO  = 0x4
)

// Standard error and clock constants.
const (
	_CLOCK_REALTIME  = 0
	_CLOCK_MONOTONIC = 3

	_TIMER_RELTIME = 0
	_TIMER_ABSTIME = 1

	_PTHREAD_CREATE_DETACHED = 1

	_SC_NPROCESSORS_ONLN = 58
	_SC_PAGESIZE         = 30

	_MAXHOSTNAMELEN = 0x100
)

//
// C-like types defined in Go for syscalls.
// NOTE: The size and layout of these structs must match the target OS (Redox).
// The sizes chosen here are common for 64-bit POSIX-compliant systems.
//

type pthread_t uintptr

// Pthread mutex, 40 bytes on 64-bit Linux.
type pthread_mutex_t struct {
	__align [40]byte
}

// Pthread condition variable, 48 bytes on 64-bit Linux.
type pthread_cond_t struct {
	__align [48]byte
}

// Pthread attributes, 56 bytes on 64-bit Linux.
type pthread_attr_t struct {
	__align [56]byte
}

// For sigaction. A set of signals.
type sigset struct {
	__bits [4]uint32 // Supports up to 128 signals.
}

type siginfoFields struct {
	si_signo int32
	si_errno int32
	si_code  int32
	// below here is a union; si_addr is the only field we use
	si_addr uint64
}

type siginfo struct {
	siginfoFields

	// Pad struct to the max size in the kernel.
	_ [_si_max_size - unsafe.Sizeof(siginfoFields{})]byte
}

type sem_t struct {
	sem_size  uint32
	sem_align uint64 // 32 bit uint32
}

type mscratch struct {
	v [6]uintptr
}

const (
	_EINTR     = 0x4
	_EAGAIN    = 0xb
	_ESRCH     = 3
	_ETIMEDOUT = 60

	_PROT_NONE  = 0x0
	_PROT_READ  = 0x1
	_PROT_WRITE = 0x2
	_PROT_EXEC  = 0x4

	_MAP_ANON    = 0x20
	_MAP_PRIVATE = 0x2
	_MAP_FIXED   = 0x10

	_SI_KERNEL = 0x80
	_SI_TIMER  = -0x2

	_SIGRTMIN = 0x20

	_BUS_ADRALN = 0x1
	_BUS_ADRERR = 0x2
	_BUS_OBJERR = 0x3

	_SEGV_MAPERR = 0x1
	_SEGV_ACCERR = 0x2

	_ITIMER_REAL    = 0x0
	_ITIMER_VIRTUAL = 0x1
	_ITIMER_PROF    = 0x2

	_CLOCK_THREAD_CPUTIME_ID = 0x3

	_SIGEV_THREAD_ID = 0x4

	_AF_UNIX    = 0x1
	_SOCK_DGRAM = 0x2
)

type timespec struct {
	tv_sec  int64
	tv_nsec int64 // 32 bit is int32
}

//go:nosplit
func (ts *timespec) setNsec(ns int64) {
	ts.tv_sec = ns / 1e9
	ts.tv_nsec = ns % 1e9
}

type timeval struct {
	tv_sec  int64
	tv_usec int64
}

func (tv *timeval) set_usec(x int32) {
	tv.tv_usec = int64(x)
}

type sigactiont struct {
	sa_handler  uintptr
	sa_flags    uint64
	sa_restorer uintptr
	sa_mask     uint64
}

type itimerspec struct {
	it_interval timespec
	it_value    timespec
}

type itimerval struct {
	it_interval timeval
	it_value    timeval
}

type sigeventFields struct {
	value  uintptr
	signo  int32
	notify int32
	// below here is a union; sigev_notify_thread_id is the only field we use
	sigev_notify_thread_id int32
}

type sigevent struct {
	sigeventFields

	// Pad struct to the max size in the kernel.
	_ [_sigev_max_size - unsafe.Sizeof(sigeventFields{})]byte
}

type usigset struct {
	__val [16]uint64
}

type fpxreg struct {
	significand [4]uint16
	exponent    uint16
	padding     [3]uint16
}

type xmmreg struct {
	element [4]uint32
}

type fpstate struct {
	cwd       uint16
	swd       uint16
	ftw       uint16
	fop       uint16
	rip       uint64
	rdp       uint64
	mxcsr     uint32
	mxcr_mask uint32
	_st       [8]fpxreg
	_xmm      [16]xmmreg
	padding   [24]uint32
}

type fpxreg1 struct {
	significand [4]uint16
	exponent    uint16
	padding     [3]uint16
}

type xmmreg1 struct {
	element [4]uint32
}

type fpstate1 struct {
	cwd       uint16
	swd       uint16
	ftw       uint16
	fop       uint16
	rip       uint64
	rdp       uint64
	mxcsr     uint32
	mxcr_mask uint32
	_st       [8]fpxreg1
	_xmm      [16]xmmreg1
	padding   [24]uint32
}

type fpreg1 struct {
	significand [4]uint16
	exponent    uint16
}

type stackt struct {
	ss_sp     uintptr
	ss_flags  int32
	pad_cgo_0 [4]byte
	ss_size   uintptr
}

type mcontext struct {
	gregs       [23]uint64
	fpregs      *fpstate
	__reserved1 [8]uint64
}

type ucontext struct {
	uc_flags     uint64
	uc_link      *ucontext
	uc_stack     stackt
	uc_mcontext  mcontext
	uc_sigmask   usigset
	__fpregs_mem fpstate
}

type sigcontext struct {
	r8          uint64
	r9          uint64
	r10         uint64
	r11         uint64
	r12         uint64
	r13         uint64
	r14         uint64
	r15         uint64
	rdi         uint64
	rsi         uint64
	rbp         uint64
	rbx         uint64
	rdx         uint64
	rax         uint64
	rcx         uint64
	rsp         uint64
	rip         uint64
	eflags      uint64
	cs          uint16
	gs          uint16
	fs          uint16
	__pad0      uint16
	err         uint64
	trapno      uint64
	oldmask     uint64
	cr2         uint64
	fpstate     *fpstate1
	__reserved1 [8]uint64
}

type sockaddr_un struct {
	family uint16
	path   [108]byte
}
