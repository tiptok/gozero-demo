package user

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-demo/app/usercenter/cmd/rpc/userservice"

	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetLogic {
	return &UserGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetLogic) UserGet(req *types.UserGetReq) (resp *types.UserGetResp, err error) {
	userGetResp, err := l.svcCtx.UserServiceRpc.UserGet(l.ctx, &userservice.UserGetReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	l.Logger.Info(userGetResp, err)
	var userItem types.UserItem
	_ = copier.Copy(&userItem, userGetResp.User)
	return &types.UserGetResp{
		User: userItem,
	}, nil
}
