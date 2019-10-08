package hubserver

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/ariefdarmawan/rpchub"
	"github.com/eaciit/toolkit"
)

type Server struct {
	hub      *rpchub.Hub
	log      *toolkit.LogEngine
	listener net.Listener

	chStop chan bool
}

func NewServer() *Server {
	s := new(Server)
	s.hub = rpchub.NewHub()
	return s
}

func (s *Server) SetLog(l *toolkit.LogEngine) *Server {
	s.log = l
	return s
}

func (s *Server) Log() *toolkit.LogEngine {
	if s.log == nil {
		s.log = toolkit.NewLogEngine(true, false, "", "", "")
	}
	return s.log
}

func (s *Server) Register(objs ...interface{}) *Server {
	if s.hub == nil {
		return s
	}

	for _, o := range objs {
		rpchub.RegisterToHub(s.hub, o)
	}

	return s
}

func (s *Server) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		err = fmt.Errorf("unable to listen to address %s. %s", addr, err.Error())
		s.Log().Error(err.Error())
		return err
	}
	s.listener = l

	server := rpc.NewServer()
	if err = server.RegisterName("hub", s.hub); err != nil {
		err = fmt.Errorf("unable to create hub for %s. %s", addr, err.Error())
		s.Log().Error(err.Error())
		return err
	}

	go server.Accept(s.listener)

	s.chStop = make(chan bool)
	s.Log().Infof("Server started on %s", addr)
	return nil
}

func (s *Server) RequestToStop() {
	if s.chStop == nil {
		s.Stop()
	}

	s.chStop <- true
}

func (s *Server) WaitForStop() {
	if s.chStop == nil {
		return
	}

	<-s.chStop
}

func (s *Server) Stop() {
	if s.log != nil {
		s.log.Close()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	close(s.chStop)
	s.chStop = nil

	s.Log().Infof("Sever is stopped")
}
