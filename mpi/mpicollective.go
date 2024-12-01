package mpi

import (
	"fmt"
)

type ReductionOp func(a, b interface{}) interface{}

func Sum(a, b interface{}) interface{} {
	return a.(int) + b.(int)
}

// MPI_Reduce reduces values from all processes to the root using the specified operation
func MPI_Reduce(sendData interface{}, recvData interface{}, op ReductionOp, root int) error {
	serializedData := Serialize(sendData)
	if rank == root {
		Deserialize(serializedData, recvData)
		for i := 0; i < size; i++ {
			if i != root {
				receivedBytes, err := MPI_Recv(i, 1)
				if err != nil {
					return fmt.Errorf("error receiving data from rank %d: %v", i, err)
				}
				var receivedValue interface{}
				Deserialize(receivedBytes, &receivedValue)
				*recvData.(*int) = op(*recvData.(*int), receivedValue).(int)
			}
		}
	} else {
		err := MPI_Send(serializedData, root, 1)
		if err != nil {
			return fmt.Errorf("error sending data to root: %v", err)
		}
	}
	return nil
}

// MPI_Bcast broadcasts data from the root process to all other processes
func MPI_Bcast(data interface{}, root int) error {
	serializedData := Serialize(data)
	if rank == root {
		for i := 0; i < size; i++ {
			if i != root {
				err := MPI_Send(serializedData, i, 0)
				if err != nil {
					return fmt.Errorf("error broadcasting data to rank %d: %v", i, err)
				}
			}
		}
	} else {
		receivedData, err := MPI_Recv(root, 0)
		if err != nil {
			return fmt.Errorf("error receiving broadcast data: %v", err)
		}
		Deserialize(receivedData, data)
	}
	return nil
}
