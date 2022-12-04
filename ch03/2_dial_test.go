package ch03

import (
	"io"
	"net"
	"testing"
)

/*
서버와 연결 수립

랜덤 포트의 127.0.0.1에서 리스너를 바인딩하고 있는 서버와 TCP 연결을 수립하는 절차를 보여주는 테스트.
*/

func TestDial(t *testing.T) {
	//랜덤 포트에 리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	/*
		1.리스너를 고루틴에서 시작해서 이후의 테스트에서 클라이언트 측에서 연결할 수 있도록 함.
		리스너의 고루틴에는 TCP 수신 연결을 루프에서 받아들이고, 각 연결 처리 로직을 담당하는 고루틴을 시작.
		이러한 고루틴을 '핸들러'라고 부름.
		여기 예제에서는 핸들러가 소켓으로부터 1024바이트를 읽어서 수신한 데이터를 로깅함.
	*/
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log()
				return
			}
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					/*
						2.FIN패킷을 받고나면 Read 메서드는 io.EOF 에러를 반환하는데,
						이는 리스너 측에서는 반대편 연결이 종료되었다는 것을 의미함.
						커넥션 핸들러는 연결 객체의 Close 메서드를 호출하면 종료함.
					*/
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
	/*
		3. net.Dial 함수는 tcp 같은 네트워크의 종류로 IP주소:포트를 매개변수로 받아서 리스너로 연결을 시도함.
		IP주소 대신 호스트명, 포트번호대신 http와 같은 서비스명을 사용 할 수 있음.

	*/
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	//8
	conn.Close()
	<-done
	//9
	listener.Close()
	<-done
}
