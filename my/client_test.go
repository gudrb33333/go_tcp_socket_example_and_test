package my

import (
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:50000")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	//1.랜덤 바이트 송신
	// payload := make([]byte, 1<<24)
	// value, err := rand.Read(payload)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	//2.문자열 송신
	const payload = "The bigger the interface, the weaker the abstraction."
	_, err = conn.Write([]byte(payload))
	if err != nil {
		t.Error(err)
	}
}
