// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "syscall"

func hostname() (name string, err error) {
	// Try uname first, as it's only one system call and reading
	var un syscall.Utsname
	err = syscall.Uname(&un)

	var buf [512]byte // Enough for a DNS name.
	for i, b := range un.Nodename[:] {
		buf[i] = uint8(b)
		if b == 0 {
			name = string(buf[:i])
			break
		}
	}
	// If we got a name and it's not potentially truncated
	// (Nodename is 65 bytes), return it.
	if err == nil && len(name) > 0 && len(name) < 64 {
		return name, nil
	}
	if name != "" {
		return name, nil
	}
	return "localhost", nil
}
