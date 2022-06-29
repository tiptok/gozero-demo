package svc

import (
	"github.com/go-pg/pg/v10"
	cache "github.com/tiptok/gocomm/pkg/cache"
	"github.com/tiptok/gocomm/pkg/cache/gzcache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"log"
	"zero-demo/app/usercenter/cmd/model"
	"zero-demo/app/usercenter/cmd/rpc/internal/config"
	repository "zero-demo/app/usercenter/pkg/db/respository"
	"zero-demo/app/usercenter/pkg/domain"
)

type ServiceContext struct {
	Config         config.Config
	RedisClient    *redis.Redis
	UserModel      model.UserModel
	UserAuthModel  model.UserAuthModel
	UserRepository domain.UsersRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	options, err := pg.ParseURL(c.DB.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	postgresql := pg.Connect(options)
	cache.RegisterCache(gzcache.NewClusterCache([]string{c.Redis.Host}, c.Redis.Pass))
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		UserModel:      model.NewUserModel(sqlConn, c.Cache),
		UserAuthModel:  model.NewUserAuthModel(sqlConn, c.Cache),
		UserRepository: repository.NewUserRepository(postgresql, cache.NewDefaultCachedRepository()),
	}
}
