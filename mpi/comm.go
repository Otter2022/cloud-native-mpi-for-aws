package mpi

// MPI_Comm_rank returns the rank of the calling process
func MPI_Comm_rank() int {
	return rank
}

// MPI_Comm_size returns the total number of processes
func MPI_Comm_size() int {
	return size
}
