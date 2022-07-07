package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"

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
	resp := &user.UserSearchResp{}
	if err := transaction.UseTrans(l.ctx, l.svcCtx.DB, func(ctx context.Context, conn transaction.Conn) error {
		total, users, err := l.svcCtx.UserRepository.Find(l.ctx, conn, map[string]interface{}{"limit": int(in.PageSize), "offset": int((in.PageNumber - 1) * in.PageSize)})
		if err != nil {
			return err
		}
		resp.Total = total
		err = copier.Copy(&resp.List, users)
		if err != nil {
			return err
		}
		return nil
	}, false); err != nil {
		return nil, err
	}
	return resp, nil
}
