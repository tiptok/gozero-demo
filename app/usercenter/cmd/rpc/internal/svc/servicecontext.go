package svc

import (
	cache "github.com/tiptok/gocomm/pkg/cache"
	"github.com/tiptok/gocomm/pkg/cache/gzcache"
	log "github.com/tiptok/gocomm/pkg/log"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	rawlog "log"
	"os"
	"time"
	"zero-demo/app/usercenter/cmd/model"
	"zero-demo/app/usercenter/cmd/rpc/internal/config"
	"zero-demo/app/usercenter/pkg/db/repository"
	"zero-demo/app/usercenter/pkg/db/transaction"
	"zero-demo/app/usercenter/pkg/domain"
)

type ServiceContext struct {
	Config             config.Config
	RedisClient        *redis.Redis
	UserModel          model.UserModel
	UserAuthModel      model.UserAuthModel
	UserRepository     domain.UserRepository
	UserAuthRepository domain.UserAuthRepository
	DB                 *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	//options, err := pg.ParseURL(c.DB.DataSource)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//postgresql := pg.Connect(options)

	newLogger := logger.New(
		rawlog.New(os.Stdout, "\r\n", rawlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	mlCache := cache.NewMultiLevelCacheNew(cache.WithDebugLog(true, func() log.Log {
		return log.DefaultLog
	}))
	mlCache.RegisterCache(gzcache.NewClusterCache([]string{c.Redis.Host}, c.Redis.Pass))

	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		UserModel:          model.NewUserModel(sqlConn, c.Cache),
		UserAuthModel:      model.NewUserAuthModel(sqlConn, c.Cache),
		UserRepository:     repository.NewUserRepository(cache.NewCachedRepository(mlCache)),
		UserAuthRepository: repository.NewUserAuthRepository(cache.NewCachedRepository(mlCache)),
		DB:                 db,
	}
}

func (svc *ServiceContext) DefaultDBConn() transaction.Conn {
	return transaction.NewTransactionContext(svc.DB)
}
