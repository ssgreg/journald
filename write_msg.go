// +build !go1.9

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

	f, err := c.File()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	_, n, errno := syscall.Syscall(syscall.SYS_SENDMSG, f.Fd(), uintptr(unsafe.Pointer(&msg)), 0)
	if errno != 0 {
		return int(n), errno
	}
	return int(n), nil
}
