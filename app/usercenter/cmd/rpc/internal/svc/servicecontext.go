package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zero-demo/app/usercenter/cmd/model"
	"zero-demo/app/usercenter/cmd/rpc/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	RedisClient   *redis.Redis
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		UserModel:     model.NewUserModel(sqlConn, c.Cache),
		UserAuthModel: model.NewUserAuthModel(sqlConn, c.Cache),
	}
}
