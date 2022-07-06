package main

import (
	"flag"
	"fmt"
	"zero-demo/app/usercenter/cmd/rpc/pb/user"
	"zero-demo/common/interceptor/rpcserver"

	"zero-demo/app/usercenter/cmd/rpc/internal/config"
	"zero-demo/app/usercenter/cmd/rpc/internal/server"
	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewUsercenterServer(ctx)
	userSvr := server.NewUserServiceServer(ctx)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUsercenterServer(grpcServer, srv)
		user.RegisterUserServiceServer(grpcServer, userSvr)
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(rpcserver.LoggerInterceptor)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
