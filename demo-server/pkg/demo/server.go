package demo

import (
	"context"
	"log"
	"time"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/technosophos/moniker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewDemoServer returns an instance of the DemoServer with sensible
// default values
func NewDemoServer() *Server {
	namer := moniker.New()
	return &Server{
		namer:      namer,
		clockset:   make(map[string]Clock),
		EventQueue: &EventQueue{},
	}
}

// Server implements the demo server gRPC interface
type Server struct {
	namer      moniker.Namer
	clockset   map[string]Clock
	EventQueue *EventQueue
}

// CreateClock handles the CreateClock remote procedure call
func (s *Server) CreateClock(ctx context.Context, req *pb.CreateTickerRequest) (*pb.CreateTickerResponse, error) {
	clock := Clock{
		ticker: time.NewTicker(time.Second),
		done:   make(chan (interface{})),
		name:   s.namer.NameSep("-"),
		queue:  s.EventQueue,
	}
	s.clockset[clock.name] = clock
	clock.Start()
	return &pb.CreateTickerResponse{Name: clock.name}, nil
}

// StopClock handles the StopClock remote procedure call
func (s *Server) StopClock(ctx context.Context, req *pb.StopClockRequest) (*pb.StopClockResponse, error) {
	s.clockset[req.GetName()].Stop()
	delete(s.clockset, req.GetName())
	return &pb.StopClockResponse{}, nil
}

// ListClocks responds with all the clocks that are managed by the demo server
func (s *Server) ListClocks(ctx context.Context, req *pb.ListClocksRequest) (*pb.ListClocksResponse, error) {
	names := []string{}
	for name := range s.clockset {
		names = append(names, name)
	}
	return &pb.ListClocksResponse{Names: names}, nil
}

// GetClockEvents responds with a stream of all clock-related activity
func (s *Server) GetClockEvents(req *pb.GetClockEventsRequest, stream pb.Demo_GetClockEventsServer) error {
	for {
		// time.Sleep(time.Millisecond * 20)
		event, err := s.EventQueue.Dequeue()
		if err != nil {
			continue
		}
		if err := stream.Send(event); err != nil {
			status, ok := status.FromError(err)
			if ok && (status.Code() == codes.Canceled || status.Code() == codes.Unavailable) {
				log.Println("stream closed by client")
				return nil
			}
			log.Printf("failed to send data: %v", err)
			return err
		}
	}
}

// StopAllClocks will cancel all the clocks on the demo server
func (s *Server) StopAllClocks(ctx context.Context, req *pb.StopAllClocksRequest) (*pb.StopAllClocksResponse, error) {
	names := []string{}
	for name, clock := range s.clockset {
		names = append(names, name)
		clock.Stop()
	}
	for _, name := range names {
		delete(s.clockset, name)
	}
	return &pb.StopAllClocksResponse{Names: names}, nil
}

// Clock ticks every second and sends its events to a queue
type Clock struct {
	ticker *time.Ticker
	done   chan (interface{})
	name   string
	queue  *EventQueue
	tick   bool
}

// Start launches the event loop for a clock
func (c Clock) Start() {
	c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_START})
	go func() {
		for {
			select {
			case <-c.ticker.C:
				if c.tick {
					c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_TICK})
				} else {
					c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_TOCK})
				}
				c.tick = !c.tick
				log.Printf("clock '%s' has changed\b", c.name)
			case <-c.done:
				c.ticker.Stop()
				return
			}
		}
	}()
}

// Stop cancels the event loop for a given clock
func (c Clock) Stop() {
	c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_STOP})
	c.done <- struct{}{}
}
