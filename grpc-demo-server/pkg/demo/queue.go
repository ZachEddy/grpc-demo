package demo 

import (
	"errors"
	"sync"

	"grpc-demo/grpc-demo-server/pkg/pb"
)

var (
	// ErrEmptyQueue ...
	ErrEmptyQueue = errors.New("queue is empty")
)

// Queue thing
type EventQueue struct {
	sync.Mutex
	Events []*pb.ClockEvent
}

// Dequeue ...
func (q *EventQueue) Dequeue() (*pb.ClockEvent, error) {
	q.Lock()
	defer q.Unlock()
	if len(q.Events) == 0 {
		return nil, ErrEmptyQueue
	}
	event := q.Events[0]
	q.Events = q.Events[1:]
	return event, nil
}

// Enqueue ...
func (q *EventQueue) Enqueue(event *pb.ClockEvent) {
	q.Lock()
	defer q.Unlock()
	q.Events = append(q.Events, event)
}
