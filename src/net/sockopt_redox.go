// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"os"
	"syscall"
)

func setDefaultSockopts(s, family, sotype int, ipv6only bool) error {
	if family == syscall.AF_INET6 && sotype != syscall.SOCK_RAW {
		// not implemented
		// syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, boolint(ipv6only))
	}
	if (sotype == syscall.SOCK_DGRAM || sotype == syscall.SOCK_RAW) && family != syscall.AF_UNIX {
		// also not implemented
		// return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1))
	}
	return nil
}

func setDefaultListenerSockopts(s int) error {
	err := syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err == syscall.EINVAL || err == syscall.ENOPROTOOPT {
		// udp in redox is incomplete
		return nil
	}
	return os.NewSyscallError("setsockopt", err)
}

func setDefaultMulticastSockopts(s int) error {
	err := syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err == syscall.EINVAL || err == syscall.ENOPROTOOPT {
		// udp in redox is incomplete
		return nil
	}
	return os.NewSyscallError("setsockopt", err)
}
