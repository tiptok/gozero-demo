package logic

import (
	"context"

	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSaveLogic {
	return &UserSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserSaveLogic) UserSave(in *user.UserSaveReq) (*user.UserSaveResp, error) {
	// todo: add your logic here and delete this line

	return &user.UserSaveResp{}, nil
}
