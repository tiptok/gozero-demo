package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// User
type User struct {
	Id         int64 `gorm:"primaryKey"`
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
	return fmt.Sprintf("%v:cache:user:id:%v", "project", m.Id)
}

func (m *User) CachePrimaryKeyFunc() string {
	if len(m.Mobile) == 0 {
		return ""
	}
	return fmt.Sprintf("%v:cache:user:phone:%v", "project", m.Mobile)
}
