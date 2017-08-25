package journald

import (
	"net"
	"syscall"
	"unsafe"
)

func sockaddr(addr *net.UnixAddr) (unsafe.Pointer, uint8) {
	sa := syscall.RawSockaddrUnix{Family: syscall.AF_UNIX}
	name := addr.Name
	n := len(name)

	for i := 0; i < n; i++ {
		sa.Path[i] = int8(name[i])
	}
	return unsafe.Pointer(&sa), byte(2 + n + 1) // length is family (uint16), name, NUL.
}

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
