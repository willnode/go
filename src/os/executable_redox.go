// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	_ "unsafe" // for linkname
)

//go:linkname executablePath
var executablePath string

var initCwd, initCwdErr = Getwd()

// similar to /proc/self/exe analog, but readfile instead of readlink
const redoxExePath = "/scheme/sys/exe"

func executable() (string, error) {
	path := ""
	if b, err := ReadFile(redoxExePath); err == nil {
		for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == 0) {
			b = b[:len(b)-1]
		}
		if len(b) > 0 {
			path = string(b)
		}
	}
	if path[0] == '/' {
		return path, nil
	}
	if initCwdErr != nil {
		return path, initCwdErr
	}
	return initCwd + "/" + path, nil
}
