package domain

import (
	"context"
	"time"
	"zero-demo/app/usercenter/pkg1/db/transaction"
)

type UserAuth struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime time.Time
	DelState   int64
	Version    int64 // 版本号
	UserId     int64
	AuthKey    string // 平台唯一id
	AuthType   string // 平台类型
}

type UserAuthRepository interface {
	Insert(ctx context.Context, transaction transaction.Conn, dm *UserAuth) (*UserAuth, error)
	Update(ctx context.Context, transaction transaction.Conn, dm *UserAuth) (*UserAuth, error)
	Delete(ctx context.Context, transaction transaction.Conn, dm *UserAuth) (*UserAuth, error)
	FindOne(ctx context.Context, transaction transaction.Conn, id int64) (*UserAuth, error)
	Find(ctx context.Context, transaction transaction.Conn, queryOptions map[string]interface{}) (int64, []*UserAuth, error)
}

func (m *UserAuth) Identify() interface{} {
	if m.Id == 0 {
		return nil
	}
	return m.Id
}
