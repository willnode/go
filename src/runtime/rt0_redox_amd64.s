// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// _rt0_amd64_redox is the ELF entry point. The interpreter
// already handle relocation for us, the rest is initialize
// like how crt0 in relibc is done
TEXT _rt0_amd64_redox(SB),NOSPLIT|NOFRAME,$0
	MOVQ	SP, DI		// arg 0: stack pointer (&argc)
	ANDQ	$~15, SP	// align stack to 16 bytes
	SUBQ	$8, SP
	MOVL	$0x00001F80, 0(SP)
	LDMXCSR	0(SP)
	MOVW	$0x037F, 0(SP)
	FLDCW	0(SP)
	ADDQ	$8, SP
	LEAQ	_rt0_amd64_redox_main(SB), SI	// rt0_go
	CALL	libc_relibc_start_v1(SB)
	INT	$3	// unreachable

// _rt0_amd64_redox_main is called by relibc_start_v1, with
// DI = argc, SI = argv, DX = envp. Stash envp then tail-jump 
// to rt0_go.
TEXT _rt0_amd64_redox_main(SB),NOSPLIT,$0
	MOVQ	DX, runtime·redoxEnvp(SB)
	JMP	runtime·rt0_go(SB)

// This is mainly untested
TEXT _rt0_amd64_redox_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)

// TODO _init is called by relibc_start_v1
TEXT _init(SB),NOSPLIT|NOFRAME,$0
	RET

// TODO _fini is called by relibc_start_v1
TEXT _fini(SB),NOSPLIT|NOFRAME,$0
	RET
