package mpi

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	TagBroadcast = 0
	TagReduce    = 1
	// ... other tags ...
)

type server struct {
	UnimplementedMPIServerServer
	mu       sync.Mutex
	messages map[int32][]*Message // Keyed by tag
}

func (s *server) Send(ctx context.Context, msg *Message) (*Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[msg.Tag] = append(s.messages[msg.Tag], msg)
	return &Empty{}, nil
}

func (s *server) Recv(ctx context.Context, req *RecvRequest) (*Message, error) {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, errors.New("receive timed out")
		case <-ticker.C:
			s.mu.Lock()
			for tag, msgs := range s.messages {
				if req.Tag != -1 && req.Tag != tag {
					continue
				}
				for i, msg := range msgs {
					if req.Source != -1 && req.Source != msg.Source {
						continue
					}
					// Found matching message
					s.messages[tag] = append(msgs[:i], msgs[i+1:]...)
					s.mu.Unlock()
					return msg, nil
				}
			}
			s.mu.Unlock()
		}
	}
}

// MPI_Send sends data to a specified destination with a tag
func MPI_Send(data []byte, dest int, tag int) error {
	client, err := getClient(dest)
	if err != nil {
		return err
	}
	msg := &Message{
		Source: int32(rank),
		Dest:   int32(dest),
		Tag:    int32(tag),
		Data:   data,
	}
	_, err = client.Send(context.Background(), msg)
	return err
}

// MPI_Recv receives data from a specified source with a tag
func MPI_Recv(source int, tag int) ([]byte, error) {
	req := &RecvRequest{
		Source: int32(source),
		Tag:    int32(tag),
	}
	msg, err := mpiServerInstance.Recv(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}
