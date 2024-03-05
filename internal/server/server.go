package server

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type Server struct {
	addr   string
	server net.Listener
}

func NewServer(address string) (*Server, error) {
	return &Server{addr: address}, nil
}

func (s *Server) StartServer() (err error) {
	fmt.Println("Listening on port:", s.addr)
	server, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.server = server
	defer s.server.Close()

	fmt.Println("Waiting for client...")
	for {
		connection, err := s.server.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return err
		}
		fmt.Println("Client connected:", connection.RemoteAddr())
		go s.handleConnection(connection)
	}
}

// handles an individual client connection.
func (s *Server) handleConnection(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	defer connection.Close()
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	query := string(buffer[:mLen])
	fmt.Println("Received query:", query)

	if !strings.HasPrefix(query, "grep") {
		fmt.Println("Error: query must start with 'grep'")
		_, err = connection.Write([]byte("Error: query must start with 'grep'"))
		return
	}

	cmd := exec.Command("bash", "-c", query)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not run command:", err)
	}
	trimmed := strings.TrimSpace(string(out))
	fmt.Println("Output:", trimmed)
	_, err = connection.Write([]byte(trimmed))
}

func (s *Server) StopServer() (err error) {
	return s.server.Close()
}
