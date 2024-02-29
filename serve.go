package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	helloworld "example.com/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Define your gRPC service by implementing the generated interface
type GreetingsService struct {
	helloworld.UnimplementedGreetingsServer
	start time.Time
}

// Implement the methods defined in the service interface
func (s GreetingsService) Connect(stream helloworld.Greetings_ConnectServer) error {
	l := logger{time.Now()}

	msg, err := stream.Recv()
	if err != nil {
		l.Log("Could not receive: %v", err)
		return err
	}
	l.Log("Received message: %s", msg.GetData())

	if err := stream.Send(&helloworld.World{
		Data: fmt.Sprintf("Hello from Go server! (you sent: %s)", msg.GetData()),
	}); err != nil {
		l.Log("Could not send: %v", err)
		return err
	}

	// Implement your logic here
	for {
		time.Sleep(5 * time.Second)
		response := &helloworld.World{
			Data: "Hello from Go server",
		}
		if err := stream.Send(response); err != nil {
			l.Log("Could not send: %v", err)
			return err
		}
	}
}

type logger struct {
	start time.Time
}

func (l logger) Log(msg string, args ...interface{}) {
	t := int(time.Since(l.start) / time.Second)
	if _, err := fmt.Printf("[%ds] %s\n", t, fmt.Sprintf(msg, args...)); err != nil {
		panic(err)
	}
}

func (l logger) LogFatal(msg string, args ...interface{}) {
	l.Log(msg, args...)
	os.Exit(1)
}

func main() {
	// Create a gRPC server
	server := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(enfPolicy), grpc.KeepaliveParams(serverParams))

	// Register your service with the server
	GreetingsService := &GreetingsService{
		start: time.Now(),
	}
	helloworld.RegisterGreetingsServer(server, GreetingsService)

	// Start listening on a specific port
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server started at %s", listener.Addr())

	// Start serving requests
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

var enfPolicy = keepalive.EnforcementPolicy{
	// MinTime:             5 * time.Second,
	// PermitWithoutStream: true,
}

var serverParams = keepalive.ServerParameters{
	// MaxConnectionIdle is a duration for the amount of time after which an
	// idle connection would be closed by sending a GoAway. Idleness duration is
	// defined since the most recent time the number of outstanding RPCs became
	// zero or the connection establishment.
	MaxConnectionIdle: 5 * time.Second,

	// MaxConnectionAge is a duration for the maximum amount of time a
	// connection may exist before it will be closed by sending a GoAway. A
	// random jitter of +/-10% will be added to MaxConnectionAge to spread out
	// connection storms.
	//MaxConnectionAge: 10 * time.Second, // The current default value is infinity.

	// MaxConnectionAgeGrace is an additive period after MaxConnectionAge after
	// which the connection will be forcibly closed.
	//MaxConnectionAgeGrace: 5 * time.Second, // The current default value is infinity.

	// After a duration of this time if the server doesn't see any activity it
	// pings the client to see if the transport is still alive.
	// If set below 1s, a minimum value of 1s will be used instead.
	// Time time.Duration // The current default value is 2 hours.

	// After having pinged for keepalive check, the server waits for a duration
	// of Timeout and if no activity is seen even after that the connection is
	// closed.
	// Timeout: time.Second,
}
