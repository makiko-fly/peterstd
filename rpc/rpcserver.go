package rpc

//import (
//	"github.com/micro/go-micro"
//	"github.com/micro/go-micro/client"
//	std "gitlab.wallstcn.com/wscnbackend/ivankastd"
//	"gitlab.wallstcn.com/wscnbackend/ivankastd/service"
//)
//
//type RPCConfig = std.ConfigService
//
//type RPCServer struct {
//	MicroService micro.Service
//}
//
//func (s *RPCServer) ServeForever(errorChan chan<- error) {
//	if err := s.MicroService.Run(); err != nil && errorChan != nil {
//		errorChan <- err
//	}
//}
//
//func NewGrpcServer(rpcConfig RPCConfig, opts ...micro.Option) *RPCServer {
//	svc := service.NewService(rpcConfig, opts...)
//	svc.Init()
//	return &RPCServer{MicroService: svc}
//}
//
//func NewGrpcClient(rpcConfig RPCConfig, opts ...micro.Option) client.Client {
//	return service.NewService(rpcConfig, opts...).Client()
//}
