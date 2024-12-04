package main

import (
	"fmt"
	"log"
	"math/rand"
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

	// Read matrix size from command-line argument
	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing matrix size: %v", err)
	}

	// Ensure the size can be evenly divided among processes
	if N%size != 0 {
		log.Fatalf("Matrix size %d must be divisible by the number of processes %d", N, size)
	}
	chunkSize := N / size

	// Matrices to be shared across processes
	var A, B, C []float64
	var localA, localB []float64

	if rank == ROOT {
		// Initialize matrices A and B with random values
		A = make([]float64, N*N)
		B = make([]float64, N*N)
		C = make([]float64, N*N)

		for i := range A {
			A[i] = float64(rand.Intn(10))
		}
		for i := range B {
			B[i] = float64(rand.Intn(10))
		}
	}

	// Prepare local matrices for each process
	localA = make([]float64, chunkSize*N)
	localB = make([]float64, N*N)
	localC := make([]float64, chunkSize*N)

	// Scatter matrix A rows and broadcast matrix B fully
	err = mpi.MPI_Scatter(A, localA, N*chunkSize, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Scatter for A: %v", rank, err)
	}

	err = mpi.MPI_Bcast(B, N*N, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Bcast for B: %v", rank, err)
	}

	// Matrix multiplication for local chunk
	for i := 0; i < chunkSize; i++ {
		for j := 0; j < N; j++ {
			sum := 0.0
			for k := 0; k < N; k++ {
				sum += localA[i*N+k] * localB[k*N+j]
			}
			localC[i*N+j] = sum
		}
	}

	// Gather results back to the root process
	err = mpi.MPI_Gather(localC, C, N*chunkSize, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Gather: %v", rank, err)
	}

	// Print the result matrix at ROOT
	if rank == ROOT {
		fmt.Println("Resultant Matrix C:")
		for i := 0; i < N; i++ {
			for j := 0; j < N; j++ {
				fmt.Printf("%6.1f", C[i*N+j])
			}
			fmt.Println()
		}
	}
}
