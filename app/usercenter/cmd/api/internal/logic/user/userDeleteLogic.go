package user

import (
	"context"
	"zero-demo/app/usercenter/cmd/rpc/userservice"

	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteLogic) UserDelete(req *types.UserDeleteReq) (resp *types.UserDeleteResp, err error) {
	_, err = l.svcCtx.UserServiceRpc.UserDelete(l.ctx, &userservice.UserDeleteReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.UserDeleteResp{}, nil
}
