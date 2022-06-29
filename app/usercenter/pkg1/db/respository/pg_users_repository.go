package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/cache"
	. "github.com/tiptok/gocomm/pkg/orm/pgx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"zero-demo/app/usercenter/pkg1/db/models"
	"zero-demo/app/usercenter/pkg1/db/transaction"
	"zero-demo/app/usercenter/pkg1/domain"
)

type UsersRepository struct {
	DB *pg.DB
	*cache.CachedRepository
}

func (repository *UsersRepository) Save(ctx context.Context, transaction transaction.Conn, dm *domain.Users) (*domain.Users, error) {
	var (
		err error
		m   = &models.Users{}
		tx  = transaction.DB()
	)
	if err = common.GobModelTransform(m, dm); err != nil {
		return nil, err
	}
	if dm.Identify() == nil {
		if _, err = tx.Model(m).Insert(m); err != nil {
			return nil, err
		}
		dm.Id = m.Id
		return dm, nil
	}
	queryFunc := func() (interface{}, error) {
		return tx.Model(m).WherePK().Update(m)
	}
	if _, err = repository.Query(queryFunc, m.CacheKeyFunc()); err != nil {
		return nil, err
	}
	return dm, nil
}

func (repository *UsersRepository) Remove(ctx context.Context, transaction transaction.Conn, User *domain.Users) (*domain.Users, error) {
	var (
		tx        = transaction.DB()
		UserModel = &models.Users{Id: User.Identify().(int64)}
	)
	queryFunc := func() (interface{}, error) {
		return tx.Model(UserModel).Where("id = ?", User.Id).Delete()
	}
	if _, err := repository.Query(queryFunc, UserModel.CacheKeyFunc()); err != nil {
		return User, err
	}
	return User, nil
}

func (repository *UsersRepository) FindOne(ctx context.Context, id int64) (*domain.Users, error) {
	tx := repository.DB
	UserModel := new(models.Users)
	queryFunc := func() (interface{}, error) {
		query := NewQuery(tx.Model(UserModel), nil)
		query.Where("id = ?", id)
		if err := query.First(); err != nil {
			return nil, sqlx.ErrNotFound
		}
		return UserModel, nil
	}
	cacheModel := new(models.Users)
	cacheModel.Id = id
	if err := repository.QueryCache(cacheModel.CacheKeyFunc, UserModel, queryFunc); err != nil {
		return nil, err
	}
	return repository.transformPgModelToDomainModel(UserModel)
}

func (repository *UsersRepository) FindOneByPhone(ctx context.Context, phone string) (*domain.Users, error) {
	tx := repository.DB
	UserModel := new(models.Users)
	queryFunc := func() (interface{}, error) {
		query := NewQuery(tx.Model(UserModel), nil)
		query.Where("phone = ?", phone)
		if err := query.First(); err != nil {
			return nil, fmt.Errorf("query row not found")
		}
		return UserModel, nil
	}
	cacheModel := new(models.Users)
	cacheModel.Phone = phone
	if err := repository.QueryUniqueIndexCache(cacheModel.CachePrimaryKeyFunc, UserModel, func(obj interface{}) string {
		if v, ok := obj.(*models.Users); ok {
			return v.CacheKeyFunc()
		}
		return ""
	}, queryFunc); err != nil {
		return nil, err
	}

	if UserModel.Id == 0 {
		return nil, sqlx.ErrNotFound
	}
	return repository.transformPgModelToDomainModel(UserModel)
}

func (repository *UsersRepository) Find(ctx context.Context, queryOptions map[string]interface{}) (int64, []*domain.Users, error) {
	tx := repository.DB
	var UserModels []*models.Users
	Users := make([]*domain.Users, 0)
	var query *Query
	queryFunc := func() (interface{}, error) {
		query = NewQuery(tx.Model(&UserModels), queryOptions).
			SetOrder("create_time", "sortByCreateTime").
			SetOrder("update_time", "sortByUpdateTime").
			SetOrder("id", "sortById").
			SetLimit()

		if searchByText, ok := queryOptions["searchByText"]; ok && len(searchByText.(string)) > 0 {
			query.Where(fmt.Sprintf(`name like '%%%v%%'`, searchByText))
		}
		var err error
		if query.AffectRow, err = query.SelectAndCount(); err != nil {
			return Users, err
		}
		return Users, err
	}

	if _, err := repository.Query(queryFunc); err != nil {
		return 0, nil, err
	}

	for _, UserModel := range UserModels {
		if User, err := repository.transformPgModelToDomainModel(UserModel); err != nil {
			return 0, Users, err
		} else {
			Users = append(Users, User)
		}
	}
	return int64(query.AffectRow), Users, nil
}

func (repository *UsersRepository) transformPgModelToDomainModel(UserModel *models.Users) (*domain.Users, error) {
	m := &domain.Users{}
	err := common.GobModelTransform(m, UserModel)
	return m, err
}

func NewUserRepository(db *pg.DB, cache *cache.CachedRepository) domain.UsersRepository {
	return &UsersRepository{DB: db, CachedRepository: cache}
}
