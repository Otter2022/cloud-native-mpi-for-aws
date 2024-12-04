package mpi

import (
	"fmt"
	"reflect"
)

type ReductionOp func(a, b interface{}) interface{}

// Improved Sum reduction for various numeric types
func Sum(a, b interface{}) interface{} {
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)

	switch v1.Kind() {
	case reflect.Float64:
		return v1.Float() + v2.Float()
	case reflect.Int, reflect.Int32, reflect.Int64:
		return v1.Int() + v2.Int()
	case reflect.Float32:
		return float32(v1.Float() + v2.Float())
	default:
		panic(fmt.Sprintf("Unsupported type for Sum reduction: %T", a))
	}
}

// MPI_Bcast broadcasts data from the root process to all other processes
func MPI_Bcast(data interface{}, count int, root int) error {
	// If this is the root process, send to all others
	if rank == root {
		// Serialize the entire data
		serializedData := Serialize(data)

		for i := 0; i < size; i++ {
			if i != root {
				err := MPI_Send(serializedData, i, TagBroadcast)
				if err != nil {
					return fmt.Errorf("error broadcasting to rank %d: %v", i, err)
				}
			}
		}
	} else {
		// Non-root processes receive data
		receivedData, err := MPI_Recv(root, TagBroadcast)
		if err != nil {
			return fmt.Errorf("error receiving broadcast data: %v", err)
		}

		// Deserialize into the provided data interface
		err = Deserialize(receivedData, data)
		if err != nil {
			return fmt.Errorf("error deserializing broadcast data: %v", err)
		}
	}

	return nil
}

// MPI_Reduce reduces values from all processes to the root using the specified operation
func MPI_Reduce(sendData interface{}, recvData interface{}, op ReductionOp, root int) error {
	// Serialize the send data
	serializedData := Serialize(sendData)

	if rank == root {
		// Initialize receive data with the first process's data
		Deserialize(serializedData, recvData)

		// Receive and reduce data from other processes
		for i := 0; i < size; i++ {
			if i == root {
				continue // Skip the root process itself
			}

			// Receive data from each non-root process
			receivedBytes, err := MPI_Recv(i, TagReduce)
			if err != nil {
				return fmt.Errorf("error receiving data from rank %d: %v", i, err)
			}

			// Deserialize the received data
			var receivedValue interface{}
			err = Deserialize(receivedBytes, &receivedValue)
			if err != nil {
				return fmt.Errorf("error deserializing data from rank %d: %v", i, err)
			}

			// Perform the reduction operation
			reducedValue := op(reflect.ValueOf(recvData).Elem().Interface(), receivedValue)

			// Update the receive data
			reflect.ValueOf(recvData).Elem().Set(reflect.ValueOf(reducedValue))
		}
	} else {
		// Non-root processes send their data to the root
		err := MPI_Send(serializedData, root, TagReduce)
		if err != nil {
			return fmt.Errorf("error sending data to root: %v", err)
		}
	}

	return nil
}

// MPI_Scatter distributes data from root to all processes
func MPI_Scatter(sendData interface{}, recvData interface{}, count int, root int) error {
	if rank == root {
		for i := 0; i < size; i++ {
			if i == root {
				// Copy data to root's local buffer
				copy(recvData.([]float64), sendData.([]float64)[i*count:(i+1)*count])
			} else {
				// Send data to other processes
				start := i * count
				end := (i + 1) * count
				serializedData := Serialize(sendData.([]float64)[start:end])
				err := MPI_Send(serializedData, i, TagScatter)
				if err != nil {
					return fmt.Errorf("error scattering to rank %d: %v", i, err)
				}
			}
		}
	} else {
		// Receive data from root process
		receivedData, err := MPI_Recv(root, TagScatter)
		if err != nil {
			return fmt.Errorf("error receiving scattered data: %v", err)
		}
		Deserialize(receivedData, recvData)
	}
	return nil
}

// MPI_Gather collects data from all processes to the root
func MPI_Gather(sendData interface{}, recvData interface{}, count int, root int) error {
	if rank == root {
		for i := 0; i < size; i++ {
			if i == root {
				// Copy data from root's local buffer
				copy(recvData.([]float64)[i*count:(i+1)*count], sendData.([]float64))
			} else {
				// Receive data from other processes
				receivedBytes, err := MPI_Recv(i, TagGather)
				if err != nil {
					return fmt.Errorf("error receiving gathered data from rank %d: %v", i, err)
				}
				var receivedData []float64
				Deserialize(receivedBytes, &receivedData)
				copy(recvData.([]float64)[i*count:(i+1)*count], receivedData)
			}
		}
	} else {
		// Send data to root process
		serializedData := Serialize(sendData)
		err := MPI_Send(serializedData, root, TagGather)
		if err != nil {
			return fmt.Errorf("error sending gathered data: %v", err)
		}
	}
	return nil
}

// Add new tag constants in sendrecv.go
const (
	TagScatter = 2
	TagGather  = 3
)
