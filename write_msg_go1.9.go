// +build go1.9

package journald

import (
	"net"
	"syscall"
	"unsafe"
)

func writeMsgUnix(c *net.UnixConn, oob []byte, addr *net.UnixAddr) (oobn int, err error) {
	ptr, salen := sockaddr(addr)

	var msg syscall.Msghdr
	msg.Name = (*byte)(ptr)
	msg.Namelen = uint32(salen)
	msg.Control = (*byte)(unsafe.Pointer(&oob[0]))
	msg.SetControllen(len(oob))

	rawConn, err := c.SyscallConn()
	if err != nil {
		return 0, err
	}

	var n uintptr
	var errno syscall.Errno
	err = rawConn.Write(func(fd uintptr) bool {
		_, n, errno = syscall.Syscall(syscall.SYS_SENDMSG, fd, uintptr(unsafe.Pointer(&msg)), 0)
		return true
	})
	if err == nil {
		if errno != 0 {
			err = errno
		}
	}
	return int(n), err
}
