package logic

import (
	"context"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"

	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeleteLogic) UserDelete(in *user.UserDeleteReq) (*user.UserDeleteResp, error) {
	// todo: add your logic here and delete this line
	if err := transaction.UseTrans(l.ctx, l.svcCtx.DB, func(ctx context.Context, conn transaction.Conn) error {
		user, err := l.svcCtx.UserRepository.FindOne(l.ctx, conn, in.Id)
		if err != nil {
			return err
		}
		_, err = l.svcCtx.UserRepository.Delete(l.ctx, conn, user)
		return err
	}, true); err != nil {
		return nil, err
	}
	return &user.UserDeleteResp{}, nil
}
