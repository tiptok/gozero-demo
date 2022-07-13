package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zero-demo/app/usercenter/cmd/api/internal/config"
	"zero-demo/app/usercenter/cmd/rpc/usercenter"
	"zero-demo/app/usercenter/cmd/rpc/userservice"
)

type ServiceContext struct {
	Config         config.Config
	UserCenterRpc  usercenter.Usercenter
	UserServiceRpc userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		UserCenterRpc:  usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
		UserServiceRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserCenterRpcConf)),
	}
}
