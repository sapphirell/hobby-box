package bootstrap

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "sukitime.com/v2/proto"
)

func InitExportRpcServer() {
	addr := "127.0.0.1:8000"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("监听端口%s异常\n", addr)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterSayHelloServer(grpcServer, &RPCServer{})
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Panicln("GRPC服务启动失败", err)
	}
}

type RPCServer struct {
	proto.UnimplementedSayHelloServer
}

func (s *RPCServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Println("接受到RPC", req)
	return &proto.HelloResponse{ResponseMsg: "Hi~" + req.RequestName}, nil
}
