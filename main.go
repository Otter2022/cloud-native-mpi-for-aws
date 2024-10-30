package main

import (
	"fmt"
	"log"

	mpi "github.com/Otter2022/cloud-native-mpi-for-aws/mpi"
)

func main() {
	mpi.MPI_Init()
	defer mpi.MPI_Finalize()

	rank := mpi.MPI_Comm_rank()
	size := mpi.MPI_Comm_size()
	const ROOT = 0
	N := 1000
	chunkSize := N / size

	var chunk []int
	if rank == ROOT {
		// Initialize the array
		array := make([]int, N)
		for i := 0; i < N; i++ {
			array[i] = i + 1
		}
		// Distribute chunks to other processes
		for i := 1; i < size; i++ {
			start := i * chunkSize
			end := start + chunkSize
			data := array[start:end]
			err := mpi.MPI_Send(mpi.Serialize(data), i, 0)
			if err != nil {
				log.Fatalf("Error sending data to rank %d: %v", i, err)
			}
		}
		// Root process handles its own chunk
		chunk = array[0:chunkSize]
	} else {
		// Receive chunk from ROOT
		dataBytes, err := mpi.MPI_Recv(ROOT, 0)
		if err != nil {
			log.Fatalf("Rank %d: Error receiving data: %v", rank, err)
		}
		var data []int
		mpi.Deserialize(dataBytes, &data)
		chunk = data
	}

	// Compute partial sum
	partialSum := 0
	for _, val := range chunk {
		partialSum += val
	}

	if rank != ROOT {
		// Send partial sum to ROOT
		err := mpi.MPI_Send(mpi.Serialize(partialSum), ROOT, 1)
		if err != nil {
			log.Fatalf("Rank %d: Error sending partial sum: %v", rank, err)
		}
	} else {
		totalSum := partialSum
		for i := 1; i < size; i++ {
			dataBytes, err := mpi.MPI_Recv(i, 1)
			if err != nil {
				log.Fatalf("ROOT: Error receiving partial sum from %d: %v", i, err)
			}
			var receivedSum int
			mpi.Deserialize(dataBytes, &receivedSum)
			totalSum += receivedSum
		}
		fmt.Printf("Total Sum = %d\n", totalSum)
	}
}
