# A script to generate Redox syscalls in Go source files, put this script to "src" folder
# You must have Redox build system in $HOME/redox with "prefix" generated
# You must have "$ARCH-unknown-redox-gcc" in your PATH, add it from Redox build system
# mkerrors: You need to patch "$CC" to "gcc" at the bottom of the script, also "$(uname)" to "Redox"
# mkerrors: You need to comment out "#include <netinet/ip6.h>" and "#include <sys/signal.h>" since Redox doesn't have it

ARCH=x86_64 GOARCH=amd64 GOOS=redox # intentionally not exported
export redox_sysroot="$HOME/redox/prefix/$ARCH-unknown-redox/sysroot/$ARCH-unknown-redox"
export CGO_FLAGS="--sysroot=${redox_sysroot}" \
CC=$ARCH-unknown-redox-gcc TARGET=$ARCH-unknown-redox GNU_TARGET=$ARCH-unknown-redox \
CFLAGS="-I${redox_sysroot}/include"


cd syscall
go tool cgo -godefs types_redox.go | go run mkpost.go > ztypes_redox_$GOARCH.go
./mksyscall_libc.pl -redox -tags redox,$GOARCH syscall_redox.go > zsyscall_redox_$GOARCH.go
SYSCALL=$(pwd)
bash ./mkerrors.sh  > zerrors_redox_$GOARCH.go
cd ..

cd cmd/vendor/golang.org/x/sys/unix
bash ./mkerrors.sh  > zerrors_redox_$GOARCH.go
sed -i s/sizeOf/SizeOf/g zerrors_redox_$GOARCH.go
# there are no types file here
go tool cgo -godefs $SYSCALL/types_redox.go | go run $SYSCALL/mkpost.go > ztypes_redox_$GOARCH.go
sed -i s/sizeof/Sizeof/g ztypes_redox_$GOARCH.go
sed -i "s/package syscall/package unix/g" ztypes_redox_$GOARCH.go
sed -i "s/_AT/AT/g" ztypes_redox_$GOARCH.go
sed -i "s/\[65\]int8/[65]byte/g" ztypes_redox_$GOARCH.go
cd ../../../../../..

cd runtime
# prefixed with ~, see https://github.com/golang/go/issues/23341
# we have to make it as close as possible by our hand :(
go tool cgo -godefs defs_redox.go > defs_redox_$GOARCH.go~
cd ..
