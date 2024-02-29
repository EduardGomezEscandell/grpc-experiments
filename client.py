#!/bin/python3

import time
import grpc
import hello_pb2
import hello_pb2_grpc


def run():
    port = "8000"
    # Establish an insecure connection to the server
    with grpc.insecure_channel(f'localhost:{port}') as channel:  # Adjust if your server is on another host
        # Create a stub for communication with the LandscapeHostAgent service
        stub = hello_pb2_grpc.GreetingsStub(channel)
        # Generate a stream of Hello messages
        def message_generator():
            yield hello_pb2.Hello(data="Hello from Python client!")
            while True:
                time.sleep(2)
                yield hello_pb2.Hello(data="PING from Python client")


        # Establish a bidirectional stream and process responses
        responses = stub.Connect(message_generator())
        for response in responses:
            print("Received from server:", response.data)
        

if __name__ == '__main__':
    run() 