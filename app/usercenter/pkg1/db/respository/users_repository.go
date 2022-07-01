package repository

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
	"zero-demo/app/usercenter/pkg1/db/models"
	"zero-demo/app/usercenter/pkg1/db/transaction"
	"zero-demo/app/usercenter/pkg1/domain"
)

type UsersRepository struct {
	*cache.CachedRepository
}

func (repository *UsersRepository) Insert(ctx context.Context, transaction transaction.Conn, dm *domain.User) (*domain.User, error) {
	var (
		err error
		m   = &models.User{}
		tx  = transaction.DB()
	)
	if err = common.GobModelTransform(m, dm); err != nil {
		return nil, err
	}
	if tx = tx.Model(m).Create(m); tx.Error != nil {
		return nil, tx.Error
	}
	dm.Id = m.Id
	return dm, nil
}

func (repository *UsersRepository) Update(ctx context.Context, transaction transaction.Conn, dm *domain.User) (*domain.User, error) {
	var (
		err error
		m   = &models.User{}
		tx  = transaction.DB()
	)
	if err = common.GobModelTransform(m, dm); err != nil {
		return nil, err
	}
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(m).Updates(m)
		return nil, tx.Error
	}
	if _, err = repository.Query(queryFunc, m.CacheKeyFunc()); err != nil {
		return nil, err
	}
	return dm, nil
}

func (repository *UsersRepository) Delete(ctx context.Context, transaction transaction.Conn, User *domain.User) (*domain.User, error) {
	var (
		tx        = transaction.DB()
		UserModel = &models.User{Id: User.Identify().(int64)}
	)
	queryFunc := func() (interface{}, error) {
		tx = tx.Where("id = ?", UserModel.Id).Delete(UserModel)
		return UserModel, tx.Error
	}
	if _, err := repository.Query(queryFunc, UserModel.CacheKeyFunc()); err != nil {
		return User, err
	}
	return User, nil
}

func (repository *UsersRepository) FindOne(ctx context.Context, transaction transaction.Conn, id int64) (*domain.User, error) {
	tx := transaction.DB()
	user := new(models.User)
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(user).Where("id = ?", id).Find(user)
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, sqlx.ErrNotFound
		}
		return user, tx.Error
	}
	cacheModel := new(models.User)
	cacheModel.Id = id
	if err := repository.QueryCache(cacheModel.CacheKeyFunc, user, queryFunc); err != nil {
		return nil, err
	}
	return repository.ToDomainModel(user)
}

func (repository *UsersRepository) FindOneByPhone(ctx context.Context, transaction transaction.Conn, phone string) (*domain.User, error) {
	tx := transaction.DB()
	UserModel := new(models.User)
	queryFunc := func() (interface{}, error) {
		tx := tx.Model(UserModel)
		tx.Where("mobile = ?", phone).First(UserModel)
		return UserModel, tx.Error
	}
	cacheModel := new(models.User)
	cacheModel.Mobile = phone
	if err := repository.QueryUniqueIndexCache(cacheModel.CachePrimaryKeyFunc, UserModel, func(obj interface{}) string {
		if v, ok := obj.(*models.User); ok {
			return v.CacheKeyFunc()
		}
		return ""
	}, queryFunc); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if UserModel.Id == 0 {
		return nil, sqlx.ErrNotFound
	}
	return repository.ToDomainModel(UserModel)
}

func (repository *UsersRepository) Find(ctx context.Context, transaction transaction.Conn, queryOptions map[string]interface{}) (int64, []*domain.User, error) {
	tx := transaction.DB()
	var UserModels []*models.User
	Users := make([]*domain.User, 0)
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(&UserModels).Order("id desc")

		if searchByText, ok := queryOptions["searchByText"]; ok && len(searchByText.(string)) > 0 {
			tx.Where(fmt.Sprintf(`name like '%%%v%%'`, searchByText))
		}
		if tx = tx.Find(&UserModels); tx.Error != nil {
			return Users, tx.Error
		}
		return Users, nil
	}

	if _, err := repository.Query(queryFunc); err != nil {
		return 0, nil, err
	}

	for _, UserModel := range UserModels {
		if User, err := repository.ToDomainModel(UserModel); err != nil {
			return 0, Users, err
		} else {
			Users = append(Users, User)
		}
	}
	return int64(tx.RowsAffected), Users, nil
}

func (repository *UsersRepository) ToDomainModel(m *models.User) (*domain.User, error) {
	dm := &domain.User{}
	err := common.GobModelTransform(dm, m)
	return dm, err
}

func NewUserRepository(cache *cache.CachedRepository) domain.UsersRepository {
	return &UsersRepository{CachedRepository: cache}
}
