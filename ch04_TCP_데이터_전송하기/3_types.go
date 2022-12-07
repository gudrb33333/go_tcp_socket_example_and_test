package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

/*
	동적 버퍼 사이즈 할당

	송신자와 수신자가 모두 프로토콜에 동의한 경우,
	네트워크 연결에서 가변 길이의 데이터를 읽을 수 있음.

	TLV(Type-Length-Value) 인코딩 체계는 가변 길이의 데이터를 처리하기 좋은 방법 중 하나.
	TLV 인코딩은
	데이터 유형을 나타내는 정해진 길이의 바이트,
	값의 크기를 나타내는 정해진 길이의 바이트
	값 자체를 나타내는 가변 길이의 바이트로
	표현됨.

	TLV 인코딩을 구현하는데 1바이트의 헤더와 4바이트의 길이, 총 5바이트를 헤더로 사용.
*/

/*
	TLV 인코딩 프로토콜을 구현하는 간단한 메세지 구조체
*/

/*
1.정의할 메세지 타입을 나타내는 상수, BinaryType와 StringType을 생성함.
각 타입의 세부 구현 정보를 요약한 후 필요에 맞게 타입을 생성.
보안상의 문제로 인해 최대 페이로드 크기를 반드시 정의 해주어야 함.
*/
const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 10 << 20 // 10MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

/*
각 타입별 메세지들이 구현해야하는 Payload라는 이름의 인터페이스를 정의함.
Payload 인터페이스를 구현할 각 타입별로 Bytes, String, ReadFrom, WriteTo라는 메서드를
반드시 구현해야함.

io.ReaderFrom 인터페이스와 io.WriterTo 인터페이스는 각 타입별 메시지를 reader로부터
읽을 수 있고 writer에 쓸 수 있게 해 주는 기능의 형태를 제공함.
세부 구현은 필요에 맞게 구현.
*/
type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

/* 바이트 슬라이스 */
type Binary []byte

/* 자기자신을 반환 */
func (m Binary) Bytes() []byte { return m }

/* 자기자신을 문자열로 캐스팅하여 반환 */
func (m Binary) String() string { return string(m) }

/*
io.Writer 인터페이스를 매개변수로 받아서 writer에 쓰인
바이트 수와 에러 인터페이스를 반환
1바이트의 타입을 writer에 씀.
Binary 인스턴스 자체의 값을 사용.
*/
func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1바이트 타입
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) //4바이트 크기
	if err != nil {
		return n, err
	}
	n += 4

	o, err := w.Write(m) //페이로드

	return n + int64(o), err
}

/*
ReadFrom 메서드는 reader로부터 1바이트를 typ 변수에 읽어 들인 후 타입이 BinaryType인지 확인 함.
그리고 size 변수에 다음 4바이트를 읽어 들임.
이 size 변숫값을 Binary 인스턴스의 크기로 새로운 바이트 슬라이스를 할당함.
마지막으로, Binary 인스턴스의 바이트 슬라이스를 읽음.

최대 페이로드 크기를 합리적으로 크기로 관리하면 서비스 거부 등의 악위적인 사용자로 부터
메모리 소비를 방지할 수 있음.
*/
func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	if typ != BinaryType {
		return n, errors.New("invalid Binary")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) //4바이트 크기
	if err != nil {
		return n, err
	}
	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size)
	o, err := r.Read(*m) // 페이로드

	return n + int64(o), err
}

type String string

/* 자기자신의 String 인스턴스 값을 바이트 슬라이스로 형변환 함. */
func (m String) Bytes() []byte { return []byte(m) }

/* 자기자신의 String 인스턴스 값을 베이스 타입인 string으로 형변환 함. */
func (m String) String() string { return string(m) }

/*
io.Writer 인터페이스를 매개변수로 받아서 writer에 쓰인
바이트 수와 에러 인터페이스를 반환
1바이트의 타입을 writer에 씀.
Binary 인스턴스 자체의 값을 사용.

String 인스턴스의 값을 Write하기 전에 바이트 슬라이스로 형변환함.
*/
func (m String) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType) //1바이트 타입
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	err = binary.Write(w, binary.BigEndian, uint(len(m)))
	n += 4

	o, err := w.Write([]byte(m)) //페이로드

	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ) //1바이트
	if err != nil {
		return 0, err
	}
	var n int64 = 1
	if typ != StringType {
		return n, errors.New("invalid String")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) //4바이트
	if err != nil {
		return n, err
	}
	n += 4

	buf := make([]byte, size)
	o, err := r.Read(buf) //페이로드
	if err != nil {
		return n, err
	}
	*m = String(buf)

	return n + int64(o), nil
}

func decode(r io.Reader) (Payload, error) {
	/*
		타입을 추론하기 위해 먼저 reader로부터 1바이트를 읽어 들인 후,
		payload 변수를 생성하여 디코딩된 타입의 값을 저장.
		읽어 들인 타입이 미리 정의한 상수 타입이면 payload 변수에 해당 상수 타입을 할당.
	*/
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload

	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unknown type")
	}

	_, err = payload.ReadFrom(
		io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}

	return payload, nil
}
