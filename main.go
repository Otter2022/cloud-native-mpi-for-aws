// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	mpi "github.com/Otter2022/cloud-native-mpi-for-aws/mpi"
)

func main() {
	mpi.MPI_Init()
	defer mpi.MPI_Finalize()

	rank := mpi.MPI_Comm_rank()
	size := mpi.MPI_Comm_size()
	const ROOT = 0

	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	chunkSize := N / size

	var array []int
	if rank == ROOT {
		// Initialize the array
		array = make([]int, N)
		for i := 0; i < N; i++ {
			array[i] = i + 1
		}
	}

	// Broadcast the array to all processes
	err = mpi.MPI_Bcast(&array, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Bcast: %v", rank, err)
	}

	// Each process computes its chunk
	start := rank * chunkSize
	end := start + chunkSize
	chunk := array[start:end]

	// Compute partial sum
	partialSum := 0
	for _, val := range chunk {
		partialSum += val
	}

	// Reduce partial sums to total sum
	var totalSum int
	err = mpi.MPI_Reduce(partialSum, &totalSum, mpi.Sum, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Reduce: %v", rank, err)
	}

	if rank == ROOT {
		fmt.Printf("Total sum is %d\n", totalSum)
	}
}
