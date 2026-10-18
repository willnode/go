// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"internal/runtime/atomic"
	"unsafe"
)

//go:cgo_import_dynamic libc_epoll_create1 epoll_create1 "libc.so.6"
//go:cgo_import_dynamic libc_epoll_ctl epoll_ctl "libc.so.6"
//go:cgo_import_dynamic libc_epoll_wait epoll_wait "libc.so.6"

//go:linkname libc_epoll_create1 libc_epoll_create1
//go:linkname libc_epoll_ctl libc_epoll_ctl
//go:linkname libc_epoll_wait libc_epoll_wait

var (
	libc_epoll_create1 libcFunc
	libc_epoll_ctl     libcFunc
	libc_epoll_wait    libcFunc

	epfd            int32         = -1
	netpollBreakRd  uintptr       // self-pipe read end
	netpollBreakWr  uintptr       // self-pipe write end
	netpollWakeSig  atomic.Uint32 // dedupes netpollBreak writes
	netpollBreakTag byte          // its address is the epoll_event.data sentinel
)

const (
	_EPOLLIN       = 0x001
	_EPOLLOUT      = 0x004
	_EPOLLERR      = 0x008
	_EPOLLHUP      = 0x010
	_EPOLL_CLOEXEC = 0x0100_0000
	_EPOLL_CTL_ADD = 1
	_EPOLL_CTL_DEL = 2
	_EPOLL_CTL_MOD = 3
)

type redoxEpollEvent struct {
	events uint32
	_pad0  uint32
	data   uint64
	_pad1  uint64
}

//go:nosplit
//go:cgo_unsafe_args
func epollcreate1(flags int32) int32 {
	return int32(sysvicall1(&libc_epoll_create1, uintptr(flags)))
}

//go:nosplit
//go:cgo_unsafe_args
func epollctl(epfd, op, fd int32, ev *redoxEpollEvent) int32 {
	return int32(sysvicall4(&libc_epoll_ctl,
		uintptr(epfd), uintptr(op), uintptr(fd),
		uintptr(unsafe.Pointer(ev))))
}

//go:nosplit
//go:cgo_unsafe_args
func epollwait(epfd int32, ev *redoxEpollEvent, nev int32, timeoutms int32) int32 {
	return int32(sysvicall4(&libc_epoll_wait,
		uintptr(epfd), uintptr(unsafe.Pointer(ev)),
		uintptr(nev), uintptr(timeoutms)))
}

//go:nosplit
func netpollErrno() int32 {
	return *getg().m.perrno
}

func netpollinit() {
	fd := epollcreate1(_EPOLL_CLOEXEC)
	if fd < 0 {
		println("runtime: epoll_create1 failed with errno=", netpollErrno())
		throw("runtime: netpollinit failed")
	}
	epfd = fd

	// Self-pipe for netpollBreak
	rfd, wfd, errno := pipe2(_O_CLOEXEC | _O_NONBLOCK)
	if errno != 0 {
		println("runtime: pipe2 failed with errno=", errno)
		throw("runtime: netpollinit pipe2 failed")
	}
	netpollBreakRd = uintptr(rfd)
	netpollBreakWr = uintptr(wfd)

	ev := redoxEpollEvent{
		events: _EPOLLIN,
		data:   uint64(uintptr(unsafe.Pointer(&netpollBreakTag))),
	}
	if epollctl(epfd, _EPOLL_CTL_ADD, rfd, &ev) != 0 {
		println("runtime: epoll_ctl ADD break-pipe failed errno=", netpollErrno())
		throw("runtime: netpollinit epoll_ctl failed")
	}
}

func netpollIsPollDescriptor(fd uintptr) bool {
	return fd == uintptr(epfd) || fd == netpollBreakRd || fd == netpollBreakWr
}

func netpollopen(fd uintptr, pd *pollDesc) uintptr {
	tp := taggedPointerPack(unsafe.Pointer(pd), pd.fdseq.Load())
	ev := redoxEpollEvent{
		events: _EPOLLIN | _EPOLLOUT,
		data:   uint64(uintptr(tp)),
	}
	if epollctl(epfd, _EPOLL_CTL_ADD, int32(fd), &ev) != 0 {
		return uintptr(netpollErrno())
	}
	return 0
}

func netpollclose(fd uintptr) uintptr {
	var ev redoxEpollEvent
	if epollctl(epfd, _EPOLL_CTL_DEL, int32(fd), &ev) != 0 {
		return uintptr(netpollErrno())
	}
	return 0
}

func netpollarm(pd *pollDesc, mode int) {
	throw("runtime: unused")
}

// netpollBreak wakes up a blocked netpoll by writing one byte into the
// self-pipe's write end. Dedupes via netpollWakeSig so repeated calls
// don't fill the pipe buffer.
//
//go:nosplit
func netpollBreak() {
	if !netpollWakeSig.CompareAndSwap(0, 1) {
		return
	}
	var b byte = 1
	for {
		n := write1(netpollBreakWr, noescape(unsafe.Pointer(&b)), 1)
		if n == 1 {
			return
		}
		if n == -_EINTR {
			continue
		}
		if n == -_EAGAIN {
			return
		}
		println("runtime: netpollBreak write failed with", -n)
		throw("runtime: netpollBreak write failed")
	}
}

// netpoll checks for ready network connections. Returns a list of
// goroutines that became runnable and a delta to add to netpollWaiters.
//
// delay < 0: blocks indefinitely
// delay == 0: polls
// delay > 0: blocks for up to that many nanoseconds
func netpoll(delay int64) (gList, int32) {
	if epfd == -1 {
		return gList{}, 0
	}
	var waitms int32
	switch {
	case delay < 0:
		waitms = -1
	case delay == 0:
		waitms = 0
	case delay < 1e6:
		waitms = 1
	case delay < 1e15:
		waitms = int32(delay / 1e6)
	default:
		// ~11.5 days — arbitrary cap.
		waitms = 1e9
	}

	var events [128]redoxEpollEvent
retry:
	n := epollwait(epfd, &events[0], int32(len(events)), waitms)
	if n < 0 {
		errno := netpollErrno()
		if errno != _EINTR {
			println("runtime: epoll_wait on fd", epfd, "failed with", errno)
			throw("runtime: netpoll failed")
		}
		if waitms > 0 {
			return gList{}, 0
		}
		goto retry
	}

	var toRun gList
	delta := int32(0)
	for i := int32(0); i < n; i++ {
		ev := &events[i]
		if ev.events == 0 {
			continue
		}

		if ev.data == uint64(uintptr(unsafe.Pointer(&netpollBreakTag))) {
			if ev.events != _EPOLLIN {
				println("runtime: netpoll: break-pipe ready for", ev.events)
				throw("runtime: netpoll: break-pipe ready for unexpected event")
			}
			if delay != 0 {
				// Drain the pipe so future netpollBreaks can fire.
				var buf [16]byte
				for {
					r := read(int32(netpollBreakRd), noescape(unsafe.Pointer(&buf[0])), int32(len(buf)))
					if r <= 0 {
						break
					}
				}
				netpollWakeSig.Store(0)
			}
			continue
		}

		var mode int32
		if ev.events&(_EPOLLIN|_EPOLLHUP|_EPOLLERR) != 0 {
			mode += 'r'
		}
		if ev.events&(_EPOLLOUT|_EPOLLHUP|_EPOLLERR) != 0 {
			mode += 'w'
		}
		if mode != 0 {
			tp := taggedPointer(uintptr(ev.data))
			pd := (*pollDesc)(tp.pointer())
			tag := tp.tag()
			if pd.fdseq.Load() == tag {
				pd.setEventErr(ev.events == _EPOLLERR, tag)
				delta += netpollready(&toRun, pd, mode)
			}
		}
	}
	return toRun, delta
}
