# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import hello as hello__pb2


class GreetingsStub(object):
    """Greetings service.
    The connection is made from the hostagent (client) to the landscape server (sass or on-prem).
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Connect = channel.stream_stream(
                '/helloworld.Greetings/Connect',
                request_serializer=hello__pb2.Hello.SerializeToString,
                response_deserializer=hello__pb2.World.FromString,
                )


class GreetingsServicer(object):
    """Greetings service.
    The connection is made from the hostagent (client) to the landscape server (sass or on-prem).
    """

    def Connect(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_GreetingsServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Connect': grpc.stream_stream_rpc_method_handler(
                    servicer.Connect,
                    request_deserializer=hello__pb2.Hello.FromString,
                    response_serializer=hello__pb2.World.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'helloworld.Greetings', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Greetings(object):
    """Greetings service.
    The connection is made from the hostagent (client) to the landscape server (sass or on-prem).
    """

    @staticmethod
    def Connect(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/helloworld.Greetings/Connect',
            hello__pb2.Hello.SerializeToString,
            hello__pb2.World.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)