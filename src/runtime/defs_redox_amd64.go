package runtime

const (
	_si_max_size    = 128
	_sigev_max_size = 64
)

const (
	_O_RDONLY    = 0x00010000
	_O_WRONLY    = 0x00020000
	_O_RDWR      = 0x00030000
	_O_ACCMODE   = 0x00030000
	_O_NONBLOCK  = 0x00040000
	_O_APPEND    = 0x00080000
	_O_SHLOCK    = 0x00100000
	_O_EXLOCK    = 0x00200000
	_O_ASYNC     = 0x00400000
	_O_FSYNC     = 0x00800000
	_O_SYNC      = _O_FSYNC
	_O_CLOEXEC   = 0x01000000
	_O_CREAT     = 0x02000000
	_O_TRUNC     = 0x04000000
	_O_EXCL      = 0x08000000
	_O_DIRECTORY = 0x10000000
	_O_PATH      = 0x20000000
	_O_SYMLINK   = 0x40000000
)

// Standard signal-related constants.
const (
	_SS_DISABLE  = 4
	_SIG_BLOCK   = 1
	_SIG_UNBLOCK = 2
	_SIG_SETMASK = 3
	_NSIG        = 32
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
	CLOCK_REALTIME  = 0
	CLOCK_MONOTONIC = 4

	PTHREAD_CREATE_DETACHED = 0

	_SC_NPROCESSORS_ONLN = 58
	_SC_PAGESIZE         = 30

	_MAXHOSTNAMELEN = 0x40
)

//
// C-like types defined in Go for syscalls.
// NOTE: The size and layout of these structs must match the target OS (Redox).
// The sizes chosen here are common for 64-bit POSIX-compliant systems.
//

type pthread_t uintptr

// Pthread mutex, 12 bytes on relibc.
type pthread_mutex_t struct {
	__align [12]byte
}

// Pthread condition variable, 8 bytes on relibc.
type pthread_cond_t struct {
	__align [8]byte
}

// Pthread attributes, 32 bytes on relibc.
type pthread_attr_t struct {
	__align [32]byte
}

type sigset struct {
	__bits [2]uint32
}

type siginfo struct {
	si_signo int32
	si_errno int32
	si_code  int32
	si_pid   int32  // pid_t is int32
	si_uid   uint32 // uid_t is uint32

	// on amd64, the go compiler will add 4 bytes of padding here
	// to align the next field (si_addr) to an 8-byte boundary.

	si_addr   uintptr // *mut c_void
	si_status int32

	// the compiler will add another 4 bytes of padding here
	// to align the final field (si_value).

	si_value uintptr
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
	_ETIMEDOUT = 110

	_PROT_NONE  = 0x0
	_PROT_READ  = 0x4
	_PROT_WRITE = 0x2
	_PROT_EXEC  = 0x1

	_MAP_ANON    = 0x20
	_MAP_PRIVATE = 0x2
	_MAP_FIXED   = 0x4

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

type stackt struct {
	ss_sp     uintptr
	ss_flags  int32
	pad_cgo_0 [4]byte
	ss_size   uintptr
}

type mcontext struct {
	ymmupper [16][2]uint64
	fxsave   [29][2]uint64
	r15      uint64
	r14      uint64
	r13      uint64
	r12      uint64
	rbp      uint64
	rbx      uint64
	r11      uint64
	r10      uint64
	r9       uint64
	r8       uint64
	rax      uint64
	rcx      uint64
	rdx      uint64
	rsi      uint64
	rdi      uint64
	rflags   uint64
	rip      uint64
	rsp      uint64
}

type ucontext struct {
	__pad       [8]byte
	uc_link     *ucontext
	uc_stack    stackt
	uc_sigmask  sigset
	_sival      uintptr
	_sigcode    uint32
	_signum     uint32
	uc_mcontext mcontext
}

type sigcontext struct {
	// this mirrors the ucontext struct
	pad_cgo_0 [8]byte

	uclink    *ucontext // or *sigcontext, depending on api
	ucstack   stackt
	ucsigmask sigset

	// internal fields from relibc
	sival   uintptr
	sigcode uint32
	signum  uint32

	// the full machine state is embedded here
	ucmcontext mcontext
}

type sockaddr_un struct {
	family uint16
	path   [108]byte
}
