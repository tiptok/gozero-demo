package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserAuth struct {
	Id         int64 `gorm:"primaryKey"`
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime time.Time
	DelState   int64 `gorm:"softDelete:flag"`
	Version    int64 // 版本号
	UserId     int64
	AuthKey    string // 平台唯一id
	AuthType   string // 平台类型
}

func (m *UserAuth) TableName() string {
	return "user_auth"
}

func (m *UserAuth) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	m.DeleteTime = time.Now()
	return
}

func (m *UserAuth) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdateTime = time.Now()
	return
}

func (m *UserAuth) BeforeDelete(tx *gorm.DB) (err error) {
	m.DeleteTime = time.Now()
	return
}

func (m *UserAuth) CacheKeyFunc() string {
	if m.Id == 0 {
		return ""
	}
	return fmt.Sprintf("%v:cache:%v:id:%v", "project", m.TableName(), m.Id)
}

func (m *UserAuth) CachePrimaryKeyFunc() string {
	if len("") == 0 {
		return ""
	}
	return fmt.Sprintf("%v:cache:%v:primarykey:%v", "project", m.TableName(), "key")
}
