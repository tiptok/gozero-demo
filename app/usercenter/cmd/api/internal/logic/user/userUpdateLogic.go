package user

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-demo/app/usercenter/cmd/rpc/userservice"

	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (resp *types.UserUpdateResp, err error) {
	var user userservice.UserItem
	copier.Copy(&user, req.User)
	_, err = l.svcCtx.UserServiceRpc.UserUpdate(l.ctx, &userservice.UserUpdateReq{Id: req.Id, User: &user})
	if err != nil {
		return nil, err
	}
	return &types.UserUpdateResp{}, nil
}
