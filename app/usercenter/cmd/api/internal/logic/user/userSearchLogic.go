package user

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-demo/app/usercenter/cmd/rpc/userservice"

	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSearchLogic {
	return &UserSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSearchLogic) UserSearch(req *types.UserSearchReq) (resp *types.UserSearchResp, err error) {
	userSearchResp, err := l.svcCtx.UserServiceRpc.UserSearch(l.ctx, &userservice.UserSearchReq{PageNumber: req.PageNumber, PageSize: req.PageSize})
	if err != nil {
		return nil, err
	}
	resp = &types.UserSearchResp{Total: userSearchResp.Total}
	var items []*types.UserItem
	err = copier.Copy(&items, userSearchResp.List)
	resp.List = items
	return
}
