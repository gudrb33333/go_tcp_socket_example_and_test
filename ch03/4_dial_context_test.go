package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

/*
	데드라인 컨텍스트를 사용하여 연결을 타임아웃하기

	컨텍스트를 사용하면 더욱 현대적인 방법으로 연결 시도를 타임아웃 시킬 수 있음.
	컨텍스트란 비동기 프로세스에 취소 시그널을 보낼 수 있는 객체.
*/
/*
	연결 시도를 타임아웃하기 위해 데드라인 컨택스트 사용하기
*/

func TestDialContext(t *testing.T) {
	dl := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()

	var d net.Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// 컨텍스트의 데드라인이 지나기 위해 충분히 긴 시간 동안 대기함.
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}

	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}
