package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

type User struct {
	Id         int64 `gorm:"primaryKey"`
	CreateTime time.Time
	UpdateTime time.Time
	DeleteTime time.Time
	DelState   soft_delete.DeletedAt `gorm:"softDelete:flag"`
	Version    int64                 // 版本号
	Mobile     string
	Password   string
	Nickname   string
	Sex        int64 // 性别 0:男 1:女
	Avatar     string
	Info       string
}

func (m *User) TableName() string {
	return "user"
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	m.DeleteTime = time.Now()
	return
}

func (m *User) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdateTime = time.Now()
	return
}

func (m *User) CacheKeyFunc() string {
	if m.Id == 0 {
		return ""
	}
	return fmt.Sprintf("%v:cache:%v:id:%v", "project", m.TableName(), m.Id)
}

func (m *User) CacheKeyFuncByObject(obj interface{}) string {
	if v, ok := obj.(*User); ok {
		return v.CacheKeyFunc()
	}
	return ""
}

func (m *User) CachePrimaryKeyFunc() string {
	if len(m.Mobile) == 0 {
		return ""
	}
	return fmt.Sprintf("%v:cache:%v:primarykey:%v", "project", m.TableName(), m.Mobile)
}
