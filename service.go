package main

import (
	"github.com/gorilla/websocket"
)

const compileBufferSize = 256

// 현재 요청되어진 전체 서비스 리스트
var services []*Service

// Service 요청된 컴파일 서비스 정보
type Service struct {
	conn    *websocket.Conn
	send    chan *CompileRequest
	cableID string
type ExecuteResponse struct {
	ExecuteOut string
	CompileOut string
	Err        string
	CpuTime    int
	MemSize    int
	ExitCode   int
	CableID    bson.ObjectId
}

func newService(conn *websocket.Conn, requestID string) {
	// 새로운 서비스 생성
	s := &Service{
		conn:    conn,
		send:    make(chan *CompileRequest, compileBufferSize),
		cableID: requestID,
	}
	// 서비스 목록에 추가
	services = append(services, s)

	// 고루틴 실행
	go s.readLoop()
	go s.writeLoop()
}

// Close 종료된 서비스를 서비스 목록에서 제거
func (s *Service) Close() {
	// 서비스 목록을 순회하며 대상 서비스를 제거
	for i, service := range services {
		if service == s {
			services = append(services[:i], services[i+1:]...)
			break
		}
	}

	// 서비스의 channel 닫기
	close(s.send)

	// 서비스 연결 종료.
	s.conn.Close()
}

// 소켓으로 데이터를 받아 반환
func (s *Service) read() (*CompileRequest, error) {
	var compile *CompileRequest
	// Json 데이터를 CompileRequest 에 저장
	if err := s.conn.ReadJSON(&compile); err != nil {
		return nil, err
	}
	return compile, nil
}

// 데이터를 Json 타입으로 전송
func (s *Service) write(m *CompileRequest) error {
	return s.conn.WriteJSON(m)
}

// 수신을 위한 루프
func (s *Service) readLoop() {
	for {
		c, err := s.read()
		if err != nil {
			break
		}
		c.create()
		broadcast(c)
	}
	s.Close()
}

// 송신을 위한 루프
func (s *Service) writeLoop() {
	for c := range s.send {
		// channel 로 부터 받은 데이터를 적절한 requestID에 전송.
		if s.cableID == c.CableID.Hex() {
			s.write(c)
		}
	}
}

// 1:N 통신을 위한 broadcast, 이후 services 를 map 으로 변경해 1:1 형태로 변형 해야함.
func broadcast(c *CompileRequest) {
	for _, service := range services {
		service.send <- c
	}
}
