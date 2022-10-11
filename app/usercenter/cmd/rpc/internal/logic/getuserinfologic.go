package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zero-demo/app/usercenter/cmd/model"
	"zero-demo/app/usercenter/cmd/rpc/usercenter"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"
	"zero-demo/app/usercenter/internal/pkg/domain"
	"zero-demo/common/xerr"

	"zero-demo/app/usercenter/cmd/rpc/internal/svc"
	"zero-demo/app/usercenter/cmd/rpc/pb"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNoExistsError = xerr.NewErrMsg("用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfoBak(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GetUserInfo find user db err , id:%d , err:%v", in.Id, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, "id:%d", in.Id)
	}

	var respUser usercenter.User
	_ = copier.Copy(&respUser, user)
	return &pb.GetUserInfoResp{
		User: &respUser,
	}, nil
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	var respUser usercenter.User
	if err := transaction.UseTrans(l.ctx, l.svcCtx.DB, func(ctx context.Context, conn transaction.Conn) error {
		var (
			user *domain.User
			err  error
		)
		user, err = l.svcCtx.UserRepository.FindOne(ctx, conn, in.Id)
		if err != nil {
			return err
		}
		_ = copier.Copy(&respUser, user)
		return nil
	}, false); err != nil {
		// example: 替换 gRPC error
		st := status.New(codes.DeadlineExceeded, "some error")
		//st := status.Error(codes.DeadlineExceeded, "some error")
		/*
		  header grpc-status-details-bin，用來補足 status 表現能力不夠的問題。為了統一模型，這個資訊格式也是採用 protobuf，我們可以把它想像成 error 專用欄位，內容經過 protobuf message 編碼後，會放在這個標頭中
		*/
		st, _ = st.WithDetails(&usercenter.GetUserInfoResp{})
		// 解析错误
		stFrom, _ := status.FromError(st.Err())
		if st.Code() == codes.DeadlineExceeded {
			for _, d := range stFrom.Details() {
				switch info := d.(type) {
				case *usercenter.GetUserInfoResp:
					fmt.Println(info)
				}
			}
		}
		return nil, st.Err()
	}
	return &pb.GetUserInfoResp{
		User: &respUser,
	}, nil
}
