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

	var A, B, C [][]int // Matrices A, B, and result C

	if rank == ROOT {
		// Initialize matrices A and B with random values
		A = generateMatrix(N, N)
		B = generateMatrix(N, N)
		C = make([][]int, N)
		for i := range C {
			C[i] = make([]int, N)
		}
	}

	// Broadcast matrices A and B to all processes
	err = mpi.MPI_Bcast(&A, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Bcast for A: %v", rank, err)
	}
	err = mpi.MPI_Bcast(&B, ROOT)
	if err != nil {
		log.Fatalf("Rank %d: Error in MPI_Bcast for B: %v", rank, err)
	}

	// Allocate a chunk of the result matrix for this process
	localC := make([][]int, chunkSize)
	for i := range localC {
		localC[i] = make([]int, N)
	}

	// Compute the assigned rows of the result matrix
	startRow := rank * chunkSize
	for i := 0; i < chunkSize; i++ {
		for j := 0; j < N; j++ {
			localC[i][j] = 0
			for k := 0; k < N; k++ {
				localC[i][j] += A[startRow+i][k] * B[k][j]
			}
		}
	}

	// Reduce local result matrices to the final result matrix at ROOT
	for i := 0; i < chunkSize; i++ {
		err = mpi.MPI_Reduce(localC[i], &C[startRow+i], mpi.Sum, ROOT)
		if err != nil {
			log.Fatalf("Rank %d: Error in MPI_Reduce for row %d: %v", rank, startRow+i, err)
		}
	}

	// Print the result matrix at ROOT
	if rank == ROOT {
		fmt.Println("Resultant Matrix C:")
		printMatrix(C)
	}
}

// generateMatrix creates a matrix with random integers
func generateMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(10) // Random values between 0 and 9
		}
	}
	return matrix
}

// printMatrix prints a 2D matrix
func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		for _, val := range row {
			fmt.Printf("%4d", val)
		}
		fmt.Println()
	}
}
