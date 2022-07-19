package logic

import (
	"context"
	"github.com/pkg/errors"
	"zero-demo/app/usercenter/cmd/rpc/userservice"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"
	"zero-demo/app/usercenter/internal/pkg/domain"
	"zero-demo/common/tool"
	"zero-demo/common/xerr"

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
	conn := l.svcCtx.DefaultDBConn()
	user, err := l.svcCtx.UserRepository.FindOneByPhone(l.ctx, conn, in.User.Mobile)
	if err != nil && err != domain.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "mobile:%s,err:%v", in.User.Mobile, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists mobile:%s,err:%v", in.User.Mobile, err)
	}

	//var userId int64
	if err := transaction.UseTrans(l.ctx, l.svcCtx.DB, func(ctx context.Context, conn transaction.Conn) error {
		var err error
		user := &domain.User{
			Mobile:   in.User.Mobile,
			Nickname: in.User.Nickname,
			Sex:      in.User.Sex,
			Avatar:   in.User.Avatar,
			Info:     in.User.Info,
		}
		if len(user.Nickname) == 0 {
			user.Nickname = tool.Krand(8, tool.KC_RAND_KIND_ALL)
		}
		//if len(in.Password) > 0 {
		//	user.Password = tool.Md5ByString(in.Password)
		//}
		user, err = l.svcCtx.UserRepository.Insert(ctx, conn, user)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
		}
		//userId = user.Id
		//userAuth := &domain2.UserAuth{
		//	UserId:   userId,
		//	AuthKey:  in.AuthKey,
		//	AuthType: in.AuthType,
		//}
		//if _, err := l.svcCtx.UserAuthRepository.Insert(ctx, conn, userAuth); err != nil {
		//	return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user_auth Insert err:%v,userAuth:%v", err, userAuth)
		//}
		return err
	}, true); err != nil {
		return nil, err
	}

	return &userservice.UserSaveResp{}, nil
}
