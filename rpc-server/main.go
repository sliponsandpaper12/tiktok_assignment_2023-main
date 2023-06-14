package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "./proto/pingpong" // Import your generated protobuf package

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// Create a global variable for the database connection
var db *sql.DB

func initDB() {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database_name")
	if err != nil {
		log.Fatal(err)
	}

	// Set the maximum number of open connections in the connection pool
	db.SetMaxOpenConns(10)

	// Set the maximum number of idle connections in the connection pool
	db.SetMaxIdleConns(5)

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")
}

type server struct{}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PongResponse, error) {
	if req.Message != "ping" {
		// Invalid request, respond with an error
		return nil, grpc.Errorf(grpc.InvalidArgument, "Invalid request. Only 'ping' messages are accepted.")
	}

	// Store client interaction information (client ID and timestamp) in your desired storage mechanism
	err := storeClientInteraction(req.ClientId, time.Now())
	if err != nil {
		log.Printf("Failed to store client interaction: %v\n", err)
	}


	// Return the response with "pong" message and current timestamp
	return &pb.PongResponse{
		Message:   "pong",
		Timestamp: time.Now().Unix(),
	}, nil
}

func storeClientInteraction(clientID string, timestamp time.Time) error {
	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO client_interactions (client_id, timestamp) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(clientID, timestamp)
	if err != nil {
		return err
	}

	log.Printf("Client interaction stored: ClientID=%s, Timestamp=%v\n", clientID, timestamp)
	return nil
}

func main() {
	// initialise database connection
	initDB()

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register the PingPongService server
	pb.RegisterPingPongServiceServer(grpcServer, &server{})

	// Start the server on a specified port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server started, listening on :50051")
	grpcServer.Serve(listener)

	// Close the database connection when the program exits
	defer db.Close()
}