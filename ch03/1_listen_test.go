package ch03

import (
	"net"
	"testing"
)

/*
소켓 바인딩, 연결 대기, 연결 수락
net.Listner 함수를 사용하면 수신 연결 요청 처리가 가능한 TCP서버를 작성 할 수 있다.
이러한 서버를 리스너하고 한다.

127.0.0.1 주소에서 랜덤 포트에 수신 대기 중인 리스너 생성하기
*/

func TestListener(t *testing.T) {

	/*
		1.net.Listen 함수는 네트워크 종류, ip주소:포트를 매개변수로 받음.
		반환 값으로 net.Listner 인터페이스와 에러 인터페이스를 반환한다.
		함수가 성공적으로 반환되면 리스너는 특정한 IP주소와 포트에 바인딩 된다.
		바인딩이란 운영체제가 지정한 IP와 포트를 해당 리스너에게 단독으로 할당했다는 의미.
		이미 바인딩된 포트에 리스너가 바인딩을 시도할 경우 에러를 반환.
	*/
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	/*
		2.close 메서드를 사용하여 항상 리스너를 우아하게 종료한다.
		defer를 사용하면 명시적으로 종료해야함.
	*/
	defer func() { _ = listener.Close() }()

	t.Logf("bound to %q", listener.Addr())

	/*
		3.for 루프를 이용하여 서버가 계속 수신 연결 요청을 수락하고, 고루틴에서 해당 연결을 처리한다.
		순차적으로 연결 요청을 수락하는 것은 아주 효율적인 방법이지만,
		그 이후에는 반드시 고루틴을 사용하여 각 연결을 처리해야한다.
	*/
	for {
		/*
			4.리스너의 Accept 메서드를 사용하여 수신 연결을 감지하고
			클라이언트와 서버간의 TCP 핸드쉐이크 절차가 완료될 때까지 블로킹됨.
			TCP 핸드 쉐이크가 실패하거나 리스너가 닫힌경우 에러 값으로
			nil 이외의 값이 리턴됨.

			conn의 실제 타입은 TCP 수신 연결을 수락했기 때문에 net.TCPConn객체의 포인터가 됨.
			연결 인터페이스는 서버 측면에서의	TCP연결을 나타낸다.

			그리고 for 무한루프지만 연결 요청이 없을 경우  listener.Accept()에서 블락킹됨.
		*/
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		/*
			5.TCP 연결이 수락된 후 고루틴을 실행하여 각 연결을 동시에 처리
			이후 연결 객체의 Close 메서드를 고루틴이 종료되기 전에 호출하여
			서버로 FIN 패킷을 보내 연결을 우아하게 종료 될 수 있도록 한다.
		*/
		go func(c net.Conn) {
			defer func() {
				c.Close()
				t.Logf("client connection closed")
			}()
			//로직 작성
			t.Logf("client connected")
		}(conn)
	}
}
