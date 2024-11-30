# Cloud-Native MPI for AWS EC2 Instances

This project enables AWS EC2 instances to communicate with each other similarly to how processes communicate in MPI (Message Passing Interface).

## Features

- **MPI Initialization**
  - `MPI_Init()`: Initialize the MPI environment.
  - `MPI_Finalize()`: Clean up the MPI environment.
  - `MPI_Comm_rank()`: Get the rank of the calling process.
  - `MPI_Comm_size()`: Get the total number of processes.

- **Point-to-Point Communication**
  - `MPI_Send(data []byte, dest int, tag int)`: Send data to a destination process.
  - `MPI_Recv(source int, tag int) ([]byte, error)`: Receive data from a source process.

- **Collective Communication**
  - `MPI_Bcast(data interface{}, root int)`: Broadcast data from the root process to all other processes.
  - `MPI_Reduce(sendData interface{}, recvData interface{}, op ReductionOp, root int)`: Reduce data from all processes to a single value at the root process.

## Getting Started

1. **Set Up Environment Variables**

   - `MPI_RANK`: The rank (ID) of the current process (e.g., `0`, `1`, `2`, ...).
   - `MPI_SIZE`: Total number of processes participating.
   - `MPI_ADDRESS_0`, `MPI_ADDRESS_1`, ..., `MPI_ADDRESS_N`: Network addresses (`host:port`) for each process.

   Example:

   ```bash
   export MPI_RANK=0
   export MPI_SIZE=2
   export MPI_ADDRESS_0=localhost:5000
   export MPI_ADDRESS_1=localhost:5001
