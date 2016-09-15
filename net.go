package testpack

import "net"

func GetOpenPort() (p int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()
	p = l.Addr().(*net.TCPAddr).Port
	return
}

func GetOpenPorts(n int) (p []int, err error) {
	for i := 0; i < n; i++ {
		var addr *net.TCPAddr
		var listener net.Listener
		if addr, err = net.ResolveTCPAddr("tcp", "localhost:0"); err != nil {
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			return
		}
		defer listener.Close()
		p = append(p, listener.Addr().(*net.TCPAddr).Port)
	}
	return
}
