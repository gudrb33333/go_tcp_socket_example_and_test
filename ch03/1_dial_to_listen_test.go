package ch03

import (
	"net"
	"testing"
)

/*
	1_listen_test.go를 루프로 실행한 후
	해당 주소와 포트로 dial연결하는 테스트
*/

func TestDialToListen(t *testing.T) {
	conn, err := net.Dial("tcp", "1_listen_test 돌린 후 bound 되는 주소로 수정")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn2, err := net.Dial("tcp", "1_listen_test 돌린 후 bound 되는 주소로 수정")
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()
}
