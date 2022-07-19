package logic

import (
	"context"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"

	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateLogic) UserUpdate(in *user.UserUpdateReq) (*user.UserUpdateResp, error) {
	if err := transaction.UseTrans(l.ctx, l.svcCtx.DB, func(ctx context.Context, conn transaction.Conn) error {
		user, err := l.svcCtx.UserRepository.FindOne(l.ctx, conn, in.Id)
		if err != nil {
			return err
		}
		user.Nickname = in.User.Nickname
		user.Avatar = in.User.Avatar
		user.Sex = in.User.Sex
		user.Info = in.User.Info
		_, err = l.svcCtx.UserRepository.Update(l.ctx, conn, user)
		return err
	}, true); err != nil {
		return nil, err
	}
	return &user.UserUpdateResp{}, nil
}
