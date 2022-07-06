package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-demo/app/usercenter/cmd/rpc/userservice"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"
	"zero-demo/app/usercenter/internal/pkg/domain"

	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetLogic {
	return &UserGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserGetLogic) UserGet(in *user.UserGetReq) (*user.UserGetResp, error) {
	// todo: add your logic here and delete this line
	var respUser userservice.UserItem
	if err := transaction.UseTrans(l.ctx, l.svcCtx.DB, func(ctx context.Context, conn transaction.Conn) error {
		var (
			user *domain.User
			err  error
		)
		user, err = l.svcCtx.UserRepository.FindOne(ctx, conn, in.Id)
		if err != nil {
			return err
		}
		err = copier.Copy(&respUser, user)
		return err
	}, false); err != nil {
		return nil, err
	}
	return &user.UserGetResp{
		User: &respUser,
	}, nil
}
