package ch03

import (
	"net"
	"syscall"
	"testing"
	"time"
)

/*
	타임아웃과 일시적인 에러 이해하기

	완벽한 조건에서는 연결시도는 즉시성공, 최악의 상황대비
	net패키지에는 에러에 대한 기준을 주는 정보가 많음.
*/

/*
	DialTimeout 함수를 이용한 연결 시도에 대한 타임아웃

	Dial 함수를 이용하면 각 연결 시도의 타임아웃을 운영체제의 타임아웃 시간에 의존해야함.
	그래서 코드상에서 타임아웃을 제어해야함.
*/

/*
	TCP 연결 시도 시 타임아웃 기간 설정하기
*/

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(err)
	}
	if !nErr.Timeout() {
		t.Fatal("error is not a timeout")
	}
}
