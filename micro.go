package peterstd

//import (
//	"context"
//
//	"github.com/micro/go-micro"
//	"github.com/micro/go-micro/client"
//	"github.com/micro/go-micro/metadata"
//	"github.com/micro/go-micro/server"
//	std "gitlab.wallstcn.com/wscnbackend/ivankastd"
//	"gitlab.wallstcn.com/wscnbackend/ivankastd/service"
//)
//
//type MicroConfig = std.ConfigService
//
//type RPCServer struct {
//	service micro.Service
//}
//
//func (s *RPCServer) GetService() micro.Service {
//	return s.service
//}
//
//// 注册handler到service
//func (s *RPCServer) Register(handlers ...interface{}) error {
//	for _, h := range handlers {
//		handler := s.service.Server().NewHandler(h)
//		if err := s.service.Server().Handle(handler); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//// 运行RPCServer
//func (s *RPCServer) ServeForever(errorChan chan<- error) {
//	if cap(errorChan) == 0 {
//		panic("Capacity of error channel should > 0")
//	}
//	go func(errorChan chan<- error) {
//		if err := s.service.Run(); err != nil {
//			errorChan <- err
//		}
//	}(errorChan)
//}
//
//// ---------------------------------------------------------------------------------------------------------------------
//// go-micro servers
//// NewMicroService 调用 ivankastd 来生成Micro的Service，并加上自定义的中间件
//func NewMicroService(svcCfg std.ConfigService, opts ...micro.Option) micro.Service {
//	//customOptions := []micro.Option{
//	//micro.Name(svcCfg.SvcName),
//	//micro.WrapHandler(ContextHandlerWrapper),
//	//}
//	//opts = append(opts, customOptions...)
//	return service.NewService(svcCfg, opts...)
//}
//
//// ContextHandlerWrapper 在rpc函数调用前将，rpc自带的matedata转化到flash通用的matedata
//func ContextHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
//	return func(ctx context.Context, req server.Request, rsp interface{}) error {
//		mmd, _ := metadata.FromContext(ctx)
//		md := make(map[string]string)
//		for k, v := range mmd {
//			md[k] = v
//		}
//		ctx = context.WithValue(ctx, ContextMetadata, md)
//		return fn(ctx, req, rsp)
//	}
//}
//
//// ---------------------------------------------------------------------------------------------------------------------
//// go-micro client
//
//func NewMicroClient(svcCfg std.ConfigService, opts ...micro.Option) client.Client {
//	return service.NewService(svcCfg, opts...).Client()
//}
//
//// 将Client调用时ctx中的request-id放入Request的Header中
//func ContextCallWrapper(c client.CallFunc) client.CallFunc {
//	return func(ctx context.Context, address string, req client.Request, rsp interface{}, opts client.CallOptions) error {
//		rid, ok := ctx.Value(ContextRequestID).(string)
//		if !ok {
//			return c(ctx, address, req, rsp, opts)
//		}
//
//		md, ok := metadata.FromContext(ctx)
//		if !ok {
//			md = make(map[string]string)
//		}
//		if rid != "" {
//			md[ContextRequestID] = rid
//		}
//		ctx = metadata.NewContext(ctx, md)
//		return c(ctx, address, req, rsp, opts)
//	}
//}
//
//func NewRPCServer(microConfig MicroConfig, handlers ...interface{}) (*RPCServer, error) {
//	svc := NewMicroService(microConfig)
//	svc.Init()
//	s := &RPCServer{
//		service: svc,
//	}
//	if err := s.Register(handlers...); err != nil {
//		return nil, err
//	}
//	return s, nil
//}
//
//// NewGrpcServer initializes the gRPC server.
//func NewGrpcServer(microConfig MicroConfig, handlers ...interface{}) (*RPCServer, error) {
//	svc := NewMicroService(microConfig)
//	svc.Init()
//	s := &RPCServer{
//		service: svc,
//	}
//	return s, nil
//}
