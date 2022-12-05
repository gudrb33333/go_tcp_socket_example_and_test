package main

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

/*
	고정된 버퍼에 데이터 읽기

	Go의 TCP연결은 io.Reader 인터페이스를 구현하였기 때문에, 네트워크 연결로 부터 데이터를 읽을 수 있음.
*/

func TestReadIntoBuffer(t *testing.T) {
	/*
		1. 클라이언트가 읽어들일 16MB 페이로드의 랜덤 데이터를 생성.
		클라이언트의 512KB 버퍼에서 읽어 들일 수 있는 데이터 양보다
		더 많기 때문에 for 루프에서 몇 번 반복적으로 순회해서 읽어 들여야함.
	*/
	payload := make([]byte, 1<<24)
	value, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("payload is %d bytes", value)

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1<<19)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("read %d bytes", n)
	}

	conn.Close()
}
