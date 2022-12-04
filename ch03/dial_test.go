package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	//랜덤 포트에 리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	//1
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			//2
			conn, err := listener.Accept()
			if err != nil {
				t.Log()
				return
			}
			//3
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					//4
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()
	//5, 6, 7
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(listener.Addr().String())

	//8
	conn.Close()
	<-done
	//9
	listener.Close()
	<-done
}
