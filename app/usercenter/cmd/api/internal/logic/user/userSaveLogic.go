package user

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-demo/app/usercenter/cmd/rpc/userservice"

	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSaveLogic {
	return &UserSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSaveLogic) UserSave(req *types.UserSaveReq) (resp *types.UserSaveResp, err error) {
	var user userservice.UserItem
	copier.Copy(&user, req.User)
	_, err = l.svcCtx.UserServiceRpc.UserSave(l.ctx, &userservice.UserSaveReq{User: &user})
	if err != nil {
		return nil, err
	}
	return &types.UserSaveResp{}, nil
}
