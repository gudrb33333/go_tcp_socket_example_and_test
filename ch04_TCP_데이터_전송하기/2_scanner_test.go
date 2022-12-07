package main

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

/*
	Scanner 메서드를 이용하여 구분자로 구분된 데이터 읽기

	TCP는 스트림 지향적인 프로토콜이기 때문에 클라이언트는 수많은 패킷들로부터 바이트 스트림을 수신할 수 있음.
	문장이나 구문과 같은 텍스트 기반 데이터와 달리 바이너리 데이터에는 하나의 메시지가 어디서 시작하고
	끝나는 지 알 수 없음.

	만약 코드상에서 서버로 부터 일련의 이메일 메시지를 읽을 때, 바이트 스트림으로 부터 각 이메일 메세지를
	구분 할 수 있는 구분자를 찾는 코드가 존재해야함.

	또는 대신 클라이언트와 서버 간 미리 정의된 프로토콜이 존재하고 서버가 다음에 전송할 페이로드의 크기를
	나타내는 고정된 길이의 바이트를 전송 해야할 것. 이 크기를 이용하여 적절한 크기의 버퍼를 만들면 됨.

	구분자로 구분하면 문제가 생길수 있음. 그래서 bufio.Scanner 라이브러리를 쓰면 됨.

	bufio.Scanner는 간편하게 구분자로 구분된 데이터를 읽어 들일 수 있음.
	Scanner는 매개변수로 io.Reader 변수를 받음.
*/

/*
	상수값을 페이로드로 제공하는 테스트 생성하기.
*/

const payload = "The bigger the interface, the weaker the abstraction."

func TestScanner(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	/*
		1. 서버에서 문자열을 읽고 있으므로 먼저 네트워크를
		연결해서 데이터를 읽어 들일 bufio.Scanner를 생성함.

		기본적으로 스캐너는 데이터 스트림으로 부터 개행 문자(\n)를 만나면
		네트워크 연결로부터 읽은 데이터를 분할 함.

		공백이나 마침표 등의 단어 경계를 구분하는 구분자를 만날 때마다
		데이터를 분할해 주는 함수인 bufio.ScanWords를 사용하여 스캐너가
		입력 데이터를 분할하도록 함.
	*/
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	var words []string

	/*
		2. 네트워크 연결에서 읽을 데이터가 있는 한 스캐너는 계속 데이터를 읽음.
		스캐너의 Text 메서드는 네트워크 연결로부터 읽어 들인 데이터 청크를 문자열로 반환
	*/
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}

	expected := []string{"The", "bigger", "the", "interface,", "the", "weaker", "the", "abstraction."}
	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
