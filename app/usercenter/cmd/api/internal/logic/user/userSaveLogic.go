package user

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
