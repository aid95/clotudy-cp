package main

import (
	"log"

	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

const compileBufferSize = 256

// 현재 요청되어진 전체 서비스 리스트
var services []*Service

// Service 요청된 컴파일 서비스 정보
type Service struct {
	Conn    *websocket.Conn
	Send    chan *ExecuteResponse
	CableID string
}

// ExecuteResponse
type ExecuteResponse struct {
	ExecuteOut string
	ExecuteErr string
	CompileOut string
	CompileErr string
	Err        string
	CPUTime    int
	MemSize    int
	ExitCode   int
	CableID    bson.ObjectId
}

func newService(conn *websocket.Conn, cableID string) {
	// 새로운 서비스 생성
	s := &Service{
		Conn:    conn,
		Send:    make(chan *ExecuteResponse, compileBufferSize),
		CableID: cableID,
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
	close(s.Send)
	// 서비스 연결 종료.
	if err := s.Conn.Close(); err != nil {
		log.Fatal(err)
	}
}

// 수신을 위한 루프
func (s *Service) readLoop() {
	for {
		err := s.read()
		if err != nil {
			break
		}
	}
	s.Close()
}

// 소켓으로 데이터를 받아 반환
func (s *Service) read() error {
	var cr *CompileRequest
	cr.create()

	// 데이터를 CompileRequest 에 저장
	if err := s.Conn.ReadJSON(&cr); err != nil {
		log.Fatal(err)
		return err
	}

	s.Send <- cr.CompileAndRun()
	return nil
}

// 송신을 위한 루프
func (s *Service) writeLoop() {
	for c := range s.Send {
		s.write(c)
	}
}

// 데이터를 Json 타입으로 전송
func (s *Service) write(m *ExecuteResponse) error {
	return s.Conn.WriteJSON(m)
}
