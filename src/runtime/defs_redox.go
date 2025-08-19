package runtime

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

// For sigaltstack.
type stackt struct {
	ss_sp    uintptr
	ss_flags int32
	ss_size  uintptr
}

// For setitimer.
type timeval struct {
	tv_sec  int64
	tv_usec int64
}
type itimerval struct {
	it_interval timeval
	it_value    timeval
}

type sigctxt struct {
	info *siginfo
	ctxt unsafe.Pointer
}


type timespec struct {
	tv_sec  int32
	tv_nsec int32
}

//go:nosplit
func (ts *timespec) setNsec(ns int64) {
	ts.tv_sec = timediv(ns, 1e9, &ts.tv_nsec)
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

// OS-specific state for a machine (m).
type mOS struct {
	waitsemacount uint32
	waitsemamutex pthread_mutex_t
	waitsemacond  pthread_cond_t
}
