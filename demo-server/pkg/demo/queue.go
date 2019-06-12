package demo

import (
	"errors"
	"sync"

	"grpc-demo/demo-server/pkg/pb"
)

var (
	// ErrEmptyQueue is returned whenever the queue is empty
	ErrEmptyQueue = errors.New("queue is empty")
)

// EventQueue is an atomic queue that stores clock events
type EventQueue struct {
	sync.Mutex
	Events []*pb.ClockEvent
}

// Dequeue removers an event from the queue
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

// Enqueue adds an event to the queue
func (q *EventQueue) Enqueue(event *pb.ClockEvent) {
	q.Lock()
	defer q.Unlock()
	q.Events = append(q.Events, event)
}
