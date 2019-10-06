package hubserver

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/ariefdarmawan/rpchub"
	"github.com/eaciit/toolkit"
)

type server struct {
	hub      *rpchub.Hub
	log      *toolkit.LogEngine
	listener net.Listener

	chStop chan bool
}

func NewServer() *server {
	s := new(server)
	s.hub = rpchub.NewHub()
	return s
}

func (s *server) SetLog(l *toolkit.LogEngine) *server {
	s.log = l
	return s
}

func (s *server) Log() *toolkit.LogEngine {
	if s.log == nil {
		s.log = toolkit.NewLogEngine(true, false, "", "", "")
	}
	return s.log
}

func (s *server) Register(objs ...interface{}) *server {
	if s.hub == nil {
		return s
	}

	for _, o := range objs {
		rpchub.RegisterToHub(s.hub, o)
	}

	return s
}

func (s *server) Start(addr string) error {
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

func (s *server) RequestToStop() {
	if s.chStop == nil {
		s.Stop()
	}

	s.chStop <- true
}

func (s *server) WaitForStop() {
	if s.chStop == nil {
		return
	}

	<-s.chStop
}

func (s *server) Stop() {
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
