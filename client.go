package main

import (
	"context"
	"fmt"
	"os"
	"time"

	helloworld "example.com/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var start time.Time

func main() {
	port := 8080

	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		LogFatal("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := helloworld.NewGreetingsClient(conn)
	start = time.Now()

	// Make a gRPC request
	stream, err := client.Connect(context.Background())
	if err != nil {
		LogFatal("Failed to call Connect: %v", err)
	}

	// Send a message
	stream.Send(&helloworld.Hello{Data: "Hello from GO!"})
	msg, err := stream.Recv()
	if err != nil {
		LogFatal("Failed to receive a response: %v", err)
	}
	Log("Message: %v", msg.GetData())

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		receiveResponses(stream, conn)
	}()

	go pingServerLoop(ctx, stream)

	<-ctx.Done()
}
func receiveResponses(stream helloworld.Greetings_ConnectClient, conn *grpc.ClientConn) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			Log("Failed to receive a response: %v", err)
			return
		}
		state := conn.GetState()
		Log("[%v] Response: %v", state, msg.GetData())
	}
}

func pingServerLoop(ctx context.Context, stream helloworld.Greetings_ConnectClient) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}

		if err := stream.Send(&helloworld.Hello{Data: "PING from Go client"}); err != nil {
			Log("Could not send a ping: %v", err)
		}
	}
}

func printStateEveryNowAndThen(conn *grpc.ClientConn) {
	for {
		time.Sleep(time.Second)
		state := conn.GetState()
		Log("[POLL] Connection state: %v", state)
	}
}

func waitUntilStateIsShutdown(conn *grpc.ClientConn) {
	state := connectivity.Ready
	for {
		conn.WaitForStateChange(context.Background(), state)
		state = conn.GetState()
		if state == connectivity.Shutdown {
			break
		}
		Log("Connection state: %v", state)
	}
	Log("Connection state: %v", state)
}

func Log(msg string, args ...interface{}) {
	t := int(time.Since(start) / time.Second)
	if _, err := fmt.Printf("[%ds] %s\n", t, fmt.Sprintf(msg, args...)); err != nil {
		panic(err)
	}
}

func LogFatal(msg string, args ...interface{}) {
	Log(msg, args...)
	os.Exit(1)
}
