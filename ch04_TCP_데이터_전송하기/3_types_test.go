package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestPayloads(t *testing.T) {
	/*
		Binary 타입 두개와 String 타입 하나 생성.
	*/
	b1 := Binary("Clear is better than clever.")
	b2 := Binary("Don't panic.")
	s1 := Binary("Errors are values.")
	payloads := []Payload{&b1, &s1, &b2}

	/*
		리스너 생성.
	*/
	listener, err := net.Listen("tcp", "127.0.0.1:50000")
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

		/*
			conn에 payloads데이터를 write함.
		*/
		for _, p := range payloads {
			_, err = p.WriteTo(conn)
			if err != nil {
				t.Error()
				break
			}
		}
	}()

	/*
		연결 수립.
	*/
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < len(payloads); i++ {
		actual, err := decode(conn)
		if err != nil {
			t.Fatal(err)
		}
		if expected := payloads[i]; !reflect.DeepEqual(expected, actual) {
			t.Errorf("value mismatch: %v != %v", expected, actual)
			continue
		}

		t.Logf("[%T] %[1]q", actual)
	}
}

/*
Binary 타입의 최대 페이로드 크기를 강제.
*/
func TestMaxPayloadSize(t *testing.T) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(BinaryType)
	if err != nil {
		t.Fatal(err)
	}

	/*
		4바이트의 부호없는 정수를 포함한 bytes.Buffer를 생성.
		명시적으로 확인 하는것이 중요!
	*/
	err = binary.Write(buf, binary.BigEndian, uint32(1<<30)) //1GB
	if err != nil {
		t.Fatal(err)
	}

	var b Binary
	_, err = b.ReadFrom(buf)
	if err != ErrMaxPayloadSize {
		t.Fatalf("expected ErrMaxPayloadSize; actual: %v", err)
	}
}
