package mpi

import (
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

var (
	rank              int
	size              int
	mpiServer         *grpc.Server
	clients           map[int]MPIServerClient
	addresses         map[int]string
	mpiServerInstance *server
)

// MPI_Init initializes the MPI environment
func MPI_Init() {
	var err error
	rank, err = strconv.Atoi(os.Getenv("MPI_RANK"))
	if err != nil {
		log.Fatalf("MPI_RANK not set or invalid: %v", err)
	}
	size, err = strconv.Atoi(os.Getenv("MPI_SIZE"))
	if err != nil {
		log.Fatalf("MPI_SIZE not set or invalid: %v", err)
	}

	// Initialize clients and addresses
	clients = make(map[int]MPIServerClient)
	addresses = make(map[int]string)
	for i := 0; i < size; i++ {
		addr := os.Getenv("MPI_ADDRESS_" + strconv.Itoa(i))
		if addr == "" {
			log.Fatalf("MPI_ADDRESS_%d not set", i)
		}
		addresses[i] = addr
	}

	// Start the gRPC server
	mpiServerInstance = &server{
		messages: make(map[int32][]*Message),
	}
	go startServer()
}

func MPI_Finalize() {
	if mpiServer != nil {
		mpiServer.GracefulStop()
	}
}

func startServer() {
	lis, err := net.Listen("tcp", addresses[rank])
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// Increase the maximum message size
	var maxMsgSize = 1024 * 1024 * 50 // 50 MiB, adjust as needed

	mpiServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	RegisterMPIServerServer(mpiServer, mpiServerInstance)
	if err := mpiServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
