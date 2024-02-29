#!/bin/python3

import asyncio
import time
import grpc
import hello_pb2
import hello_pb2_grpc

async def run() -> None:
    port = "8080"
    now = time.time()
    channel = grpc.aio.insecure_channel(f'localhost:{port}')
    
    stub = hello_pb2_grpc.GreetingsStub(channel)

    # Create a stream
    stream = stub.Connect()

    # Send a message
    await stream.write(hello_pb2.Hello(data="Hello from Python client!"))

    # Read the response
    response = await stream.read()
    print("Received:", response.data)

    # Send pings
    while True:
        time.sleep(2)
        await stream.write(hello_pb2.Hello(data="PING from Python client"))
        response = await stream.read()
        print(f"[{int(time.time() - now)}s] [{channel.get_state()}] Received: {response.data}")        
        

if __name__ == '__main__':
    asyncio.get_event_loop().run_until_complete(run())