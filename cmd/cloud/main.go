package main

import (
	"log"
)

func main() {
	s := NewServer()
	err := s.Init()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	err = s.Start()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	defer func(s *Server) {
		err := s.Stop()
		if err != nil {
			log.Fatalf("failed to stop server: %v", err)
		}
	}(s)
}
