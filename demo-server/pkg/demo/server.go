package demo

import (
	"context"
	"fmt"
	"log"
	"time"

	"grpc-demo/demo-server/pkg/pb"

	"github.com/technosophos/moniker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewDemoServer() *DemoServer {
	namer := moniker.New()
	return &DemoServer{
		namer:      namer,
		clockset:   make(map[string]Clock),
		EventQueue: &EventQueue{},
	}
}

type DemoServer struct {
	namer      moniker.Namer
	clockset   map[string]Clock
	EventQueue *EventQueue
}

func (s *DemoServer) CreateClock(ctx context.Context, req *pb.CreateTickerRequest) (*pb.CreateTickerResponse, error) {
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

func (s *DemoServer) StopClock(ctx context.Context, req *pb.StopClockRequest) (*pb.StopClockResponse, error) {
	s.clockset[req.GetName()].Stop()
	delete(s.clockset, req.GetName())
	return &pb.StopClockResponse{}, nil
}

func (s *DemoServer) ListClocks(ctx context.Context, req *pb.ListClocksRequest) (*pb.ListClocksResponse, error) {
	names := []string{}
	for name, _ := range s.clockset {
		names = append(names, name)
	}
	return &pb.ListClocksResponse{Names: names}, nil
}

func (s *DemoServer) GetClockEvents(req *pb.GetClockEventsRequest, stream pb.Demo_GetClockEventsServer) error {
	for {
		time.Sleep(time.Millisecond * 20)
		event, err := s.EventQueue.Dequeue()
		if err != nil {
			continue
		}
		if err := stream.Send(event); err != nil {
			status, ok := status.FromError(err)
			if ok && (status.Code() == codes.Canceled || status.Code() == codes.Unavailable) {
				log.Println("Stream closed by client")
				return nil
			}
			log.Printf("Failed to send data: %v", err)
			return err
		}
	}
	return nil
}

func (s *DemoServer) StopAllClocks(ctx context.Context, req *pb.StopAllClocksRequest) (*pb.StopAllClocksResponse, error) {
	names := []string{}
	for name, clock := range s.clockset {
		names = append(names, name)
		clock.Stop()
	}
	fmt.Println(names)
	for _, name := range names {
		delete(s.clockset, name)
	}
	return &pb.StopAllClocksResponse{Names: names}, nil
}

type Clock struct {
	ticker *time.Ticker
	done   chan (interface{})
	name   string
	queue  *EventQueue
}

func (c Clock) Start() {
	c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_START})
	go func() {
		for {
			select {
			case <-c.ticker.C:
				fmt.Printf("%s ticked!\n", c.name)
				c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_RESET})
			case <-c.done:
				c.ticker.Stop()
				return
			}
		}
	}()
}

func (c Clock) Stop() {
	c.queue.Enqueue(&pb.ClockEvent{Name: c.name, Event: pb.ClockEvent_STOP})
	c.done <- struct{}{}
}
