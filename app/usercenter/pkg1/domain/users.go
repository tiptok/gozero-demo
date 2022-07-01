package domain

import (
	"context"
	"time"
	"zero-demo/app/usercenter/pkg1/db/transaction"
)

const (
	UserAdmin = iota + 1
	UserNormal
)

type User struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime time.Time
	DelState   int64
	Version    int64
	Mobile     string
	Password   string
	Nickname   string
	Sex        int64
	Avatar     string
	Info       string
}

type UsersRepository interface {
	Update(ctx context.Context, transaction transaction.Conn, dm *User) (*User, error)
	Insert(ctx context.Context, transaction transaction.Conn, dm *User) (*User, error)
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
