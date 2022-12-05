package my

import (
	"bufio"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:50000")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = listener.Close() }()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		go func(c net.Conn) {
			defer func() {
				c.Close()
				t.Logf("client connection closed")
			}()
			t.Logf("client connected")

			//1.랜덤 바이트 수신
			// buf := make([]byte, 1<<19)

			// for {
			// 	n, err := conn.Read(buf)
			// 	if err != nil {
			// 		if err != io.EOF {
			// 			t.Error(err)
			// 		}
			// 		break
			// 	}

			// 	t.Logf("read %d bytes", n)
			// }

			//2.문자열 바이트 수신 후 바이트 파싱
			scanner := bufio.NewScanner(conn)
			scanner.Split(bufio.ScanWords)

			var words []string

			for scanner.Scan() {
				words = append(words, scanner.Text())
			}

			err = scanner.Err()
			if err != nil {
				t.Error(err)
			}

			t.Logf("Scanned words: %#v", words)
		}(conn)
	}
}
