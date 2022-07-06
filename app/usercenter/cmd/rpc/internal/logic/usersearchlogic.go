package logic

import (
	"context"

	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSearchLogic {
	return &UserSearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserSearchLogic) UserSearch(in *user.UserSearchReq) (*user.UserSearchResp, error) {
	// todo: add your logic here and delete this line

	return &user.UserSearchResp{}, nil
}
