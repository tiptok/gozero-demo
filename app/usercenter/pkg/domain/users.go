package domain

import (
	"context"
	"time"
	"zero-demo/app/usercenter/pkg/db/transaction"
)

const (
	UserAdmin = iota + 1
	UserNormal
)

type Users struct {
	// 唯一标识
	Id int64 `json:"id"`
	// 名称
	Name string `json:"name"`
	// 手机号
	Phone string `json:"phone"`
	// 密码
	Passwd string `json:"-"`
	// 用户角色
	Roles []int64 `json:"roles"`
	// 1启用  2禁用
	Status int `json:"status"`
	// 管理员类型 1:管理员  2：普通员工
	AdminType int `json:"adminType"`
	// 创建时间
	CreateTime time.Time `json:"createTime"`
	// 更新时间
	UpdateTime time.Time `json:"updateTime"`
}

type UsersRepository interface {
	Save(ctx context.Context, transaction transaction.Trans, dm *Users) (*Users, error)
	Remove(ctx context.Context, transaction transaction.Trans, dm *Users) (*Users, error)
	FindOne(ctx context.Context, id int64) (*Users, error)
	FindOneByPhone(ctx context.Context, phone string) (*Users, error)
	Find(ctx context.Context, queryOptions map[string]interface{}) (int64, []*Users, error)
}

func (m *Users) Identify() interface{} {
	if m.Id == 0 {
		return nil
	}
	return m.Id
}
