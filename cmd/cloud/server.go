package main

import (
	"github.com/Persists/fcproto/internal/cloud/config"
	"github.com/Persists/fcproto/internal/cloud/connection"
	"github.com/Persists/fcproto/internal/cloud/database"
	"github.com/Persists/fcproto/internal/shared/sender"
	"log"
)

type Server struct {
	connMgr       *connection.ConnectionManager
	dbm           *database.DatabaseManager
	serverConfig  *config.ServerConfig
	senderManager *sender.SenderManager
}

func NewServer() *Server {
	return &Server{
		connMgr:       &connection.ConnectionManager{},
		dbm:           &database.DatabaseManager{},
		senderManager: sender.NewSenderManager(),
	}
}

func (s *Server) Init() error {
	config, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed to load env config: %v", err)
	}
	s.serverConfig = config

	err = s.dbm.Init(s.serverConfig)
	if err != nil {
		log.Printf("failed to initialize the database manager: %v", err)
		return err
	}

	err = s.connMgr.Init(s.serverConfig, s.dbm.GetDB())
	if err != nil {
		log.Printf("failed to initialize the connection manager: %v", err)
		return err
	}

	return nil
}

func (s *Server) Start() error {
	err := s.dbm.Start()
	if err != nil {
		log.Printf("failed to start the database manager: %v", err)
		return err
	}

	err = s.connMgr.Start()
	if err != nil {
		log.Printf("failed to start the connection manager: %v", err)
		return err
	}

	s.senderManager.Start()

	return nil
}

func (s *Server) Stop() error {
	err := s.connMgr.Stop()
	if err != nil {
		log.Printf("failed to stop the connection manager: %v", err)
		return err
	}

	err = s.dbm.Stop()
	if err != nil {
		log.Printf("failed to stop the database manager: %v", err)
		return err
	}

	return nil
}
