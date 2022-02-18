package srv

import (
	"fmt"
	"log"
	"net"

	apiv1 "github.com/easyCZ/qfy/gen/v1"
	"github.com/easyCZ/qfy/internal/db"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type CPConfig struct {
	DB db.ConnectionParams
}

func ListenAndServeControlPlane(c CPConfig) error {
	database, err := setupDB(c.DB)
	if err != nil {
		log.Fatalf("Failed to setup db: %v", err)
	}

	syntheticsRepo := db.NewSyntheticsRepository(database)
	agentsRepo := db.NewAgentsRepository(database)

	syntheticsSvc := &SyntheticsService{repo: syntheticsRepo}
	agentsSvc := &AgentService{repo: agentsRepo}

	grpcServer := grpc.NewServer()

	apiv1.RegisterAgentServiceServer(grpcServer, agentsSvc)
	apiv1.RegisterSyntheticsServiceServer(grpcServer, syntheticsSvc)

	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("Failed to listen on port 3001, %v", err)
	}

	log.Printf("Starting gRPC server on localhost:3000")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed to start: %v", err)
	}

	log.Printf("Finished serving gRPC API.")
	return nil
}

func setupDB(params db.ConnectionParams) (*gorm.DB, error) {
	database, err := db.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to setup db: %v", err)
	}

	if err := db.Migrate(database); err != nil {
		return nil, fmt.Errorf("failed to migrate DB: %v", err)
	}

	return database, nil
}
