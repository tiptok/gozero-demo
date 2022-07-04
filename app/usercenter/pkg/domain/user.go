package domain

import (
	"context"
	"time"
	"zero-demo/app/usercenter/pkg/db/transaction"
)

type User struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime time.Time
	DelState   int64
	Version    int64 // 版本号
	Mobile     string
	Password   string
	Nickname   string
	Sex        int64 // 性别 0:男 1:女
	Avatar     string
	Info       string
}

type UserRepository interface {
	Insert(ctx context.Context, transaction transaction.Conn, dm *User) (*User, error)
	Update(ctx context.Context, transaction transaction.Conn, dm *User) (*User, error)
	UpdateWithVersion(ctx context.Context, transaction transaction.Conn, dm *User) (*User, error)
	Delete(ctx context.Context, transaction transaction.Conn, dm *User) (*User, error)
	FindOne(ctx context.Context, transaction transaction.Conn, id int64) (*User, error)
	FindOneByPhone(ctx context.Context, transaction transaction.Conn, phone string) (*User, error)
	Find(ctx context.Context, transaction transaction.Conn, queryOptions map[string]interface{}) (int64, []*User, error)
}

func (m *User) Identify() interface{} {
	if m.Id == 0 {
		return nil
	}
	return m.Id
}
