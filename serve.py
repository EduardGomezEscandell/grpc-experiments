#!/bin/python3

import time
from concurrent import futures

import grpc

# Generated classes from the .proto file
import hello_pb2
import hello_pb2_grpc


# Implementation of the LandscapeHostAgent service
class GreetingsServer(hello_pb2_grpc.GreetingsServicer):

    def Connect(self, request_iterator, context):
        # Process the stream of Hello messages
        for hello in request_iterator:
            print("Received:", hello.data)
            yield hello_pb2.World(data=f"Hello back from the Python server! (you said: {hello.data})")

            # # Send a response for each message
            # while True:
            #     time.sleep(5)
                # yield hello_pb2.World(data=f"PING from Python server")

def serve_options():
    # return [
    #      ("grpc.keepalive_time_ms", 5000),
    #      ("grpc.keepalive_timeout_ms", 3000),
    # ]
    return  [
        ("grpc.max_connection_age_ms", 10_000),                                          

        # ("grpc.keepalive_time_ms", 30_000),
        # ("grpc.keepalive_timeout_ms", 10_000),
        # ("grpc.http2.min_ping_interval_without_data_ms", 5_000),
        # ("grpc.max_connection_idle_ms", 10_000),
        
        ("grpc.max_connection_age_grace_ms", 5_000),
        # ("grpc.http2.max_pings_without_data", 5),
           
        # ("grpc.keepalive_permit_without_calls", 1),
        # ("grpc.client_idle_timeout_ms", 15_000),
    ]

def serve():
    # Create a gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), options=serve_options())
    # Register the servicer with the server
    hello_pb2_grpc.add_GreetingsServicer_to_server(
        GreetingsServer(), server
    )
    # Bind the server to the port
    port = "8080"
    server.add_insecure_port(f"localhost:{port}")  # You can modify the port as needed
    server.start()
    print(f"Server started, listening on port :{port}")
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        # server.stop(0)
        pass


if __name__ == "__main__":
    serve()
