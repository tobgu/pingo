package probe

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

///////////////////////
/// TCP test server ///
///////////////////////

type testServer struct {
	port         int
	synchronizer chan struct{}
	t            *testing.T
}

var poisonPill byte = 0xFF

func startTestServer(t *testing.T) *testServer {
	s := &testServer{port: 9050, synchronizer: make(chan struct{}), t: t}
	go s.run()
	s.awaitSync()
	return s
}

func (s *testServer) failOnErr(err error) {
	if err != nil {
		s.t.Fatal(err)
	}
}

func (s *testServer) communicate(client net.Conn) (killed bool) {
	defer client.Close()
	buf, err := ioutil.ReadAll(client)
	s.failOnErr(err)

	if len(buf) > 0 && buf[0] == poisonPill {
		return true
	}

	_, err = client.Write(buf)
	s.failOnErr(err)
	return false
}

func (s *testServer) run() {
	s.t.Log("Starting to listen")
	server, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
	s.failOnErr(err)

	if l, ok := server.(*net.TCPListener); ok {
		err = l.SetDeadline(time.Now().Add(3 * time.Second))
		s.failOnErr(err)
	}
	s.sendSync()

	for {
		client, err := server.Accept()
		s.t.Log("Accepted")
		s.failOnErr(err)

		if killed := s.communicate(client); killed {
			s.t.Log("Killed")
			s.sendSync()
			return
		}
		s.t.Log("Not killed")
	}
}

func (s *testServer) sendSync() {
	s.synchronizer <- struct{}{}
}

func (s *testServer) awaitSync() {
	select {
	case <-time.After(2 * time.Second):
		s.t.Fatal("No server sync")
	case <-s.synchronizer:
	}
}

func (s *testServer) terminate() {
	con, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", s.port), time.Second)
	s.failOnErr(err)

	defer con.Close()

	err = con.SetDeadline(time.Now().Add(time.Second))
	s.failOnErr(err)

	_, err = con.Write([]byte{poisonPill})
	s.failOnErr(err)

	if tcpcon, ok := con.(*net.TCPConn); ok {
		tcpcon.CloseWrite()
	} else {
		s.t.Fatal("Not a TCP connection")
	}

	s.awaitSync()
	close(s.synchronizer)
}

///////////
// Tests //
///////////

func TestBasicPing(t *testing.T) {
	server := startTestServer(t)
	server.terminate()
}
