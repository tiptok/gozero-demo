package repository

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/tiptok/gocomm/pkg/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
	"zero-demo/app/usercenter/internal/pkg/db/models"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"
	"zero-demo/app/usercenter/internal/pkg/domain"
)

type UserRepository struct {
	*cache.CachedRepository
}

func (repository *UserRepository) Insert(ctx context.Context, transaction transaction.Conn, dm *domain.User) (*domain.User, error) {
	var (
		err error
		m   = &models.User{}
		tx  = transaction.DB()
	)
	if m, err = repository.DomainModelToModel(dm); err != nil {
		return nil, err
	}
	if tx = tx.Model(m).Save(m); tx.Error != nil {
		return nil, tx.Error
	}
	dm.Id = m.Id
	return dm, nil

}

func (repository *UserRepository) Update(ctx context.Context, transaction transaction.Conn, dm *domain.User) (*domain.User, error) {
	var (
		err error
		m   *models.User
		tx  = transaction.DB()
	)
	if m, err = repository.DomainModelToModel(dm); err != nil {
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

func (repository *UserRepository) UpdateWithVersion(ctx context.Context, transaction transaction.Conn, dm *domain.User) (*domain.User, error) {
	var (
		err error
		m   *models.User
		tx  = transaction.DB()
	)
	if m, err = repository.DomainModelToModel(dm); err != nil {
		return nil, err
	}
	oldVersion := dm.Version
	m.Version += 1
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(m).Where("id = ?", m.Id).Where("version = ?", oldVersion).Updates(m)
		return nil, tx.Error
	}
	if _, err = repository.Query(queryFunc, m.CacheKeyFunc()); err != nil {
		return nil, err
	}
	return dm, nil
}

func (repository *UserRepository) Delete(ctx context.Context, transaction transaction.Conn, dm *domain.User) (*domain.User, error) {
	var (
		tx = transaction.DB()
		m  = &models.User{Id: dm.Identify().(int64)}
	)
	queryFunc := func() (interface{}, error) {
		tx = tx.Where("id = ?", m.Id).Delete(m)
		return m, tx.Error
	}
	if _, err := repository.Query(queryFunc, m.CacheKeyFunc()); err != nil {
		return dm, err
	}
	return dm, nil
}

func (repository *UserRepository) FindOne(ctx context.Context, transaction transaction.Conn, id int64) (*domain.User, error) {
	var (
		err error
		tx  = transaction.DB()
		m   = new(models.User)
	)
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(m).Where("id = ?", id).First(m)
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return m, tx.Error
	}
	cacheModel := new(models.User)
	cacheModel.Id = id
	if err = repository.QueryCache(cacheModel.CacheKeyFunc, m, queryFunc); err != nil {
		return nil, err
	}
	return repository.ModelToDomainModel(m)
}

func (repository *UserRepository) FindOneByPhone(ctx context.Context, transaction transaction.Conn, phone string) (*domain.User, error) {
	tx := transaction.DB()
	UserModel := new(models.User)
	queryFunc := func() (interface{}, error) {
		tx := tx.Model(UserModel)
		tx.Where("mobile = ?", phone).First(UserModel)
		return UserModel, tx.Error
	}
	cacheModel := new(models.User)
	cacheModel.Mobile = phone
	if err := repository.QueryUniqueIndexCache(cacheModel.CachePrimaryKeyFunc, UserModel, cacheModel.CacheKeyFuncByObject, queryFunc); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if UserModel.Id == 0 {
		return nil, sqlx.ErrNotFound
	}
	return repository.ModelToDomainModel(UserModel)
}

func (repository *UserRepository) Find(ctx context.Context, conn transaction.Conn, queryOptions map[string]interface{}) (int64, []*domain.User, error) {
	var (
		tx    = conn.DB()
		ms    []*models.User
		dms   = make([]*domain.User, 0)
		total int64
	)
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(&ms).Order("id desc")
		if total, tx = transaction.PaginationAndCount(ctx, tx, queryOptions, &ms); tx.Error != nil {
			return dms, tx.Error
		}
		return dms, nil
	}

	if _, err := repository.Query(queryFunc); err != nil {
		return 0, nil, err
	}

	for _, item := range ms {
		if dm, err := repository.ModelToDomainModel(item); err != nil {
			return 0, dms, err
		} else {
			dms = append(dms, dm)
		}
	}
	return total, dms, nil
}

func (repository *UserRepository) ModelToDomainModel(from *models.User) (*domain.User, error) {
	to := &domain.User{}
	err := copier.Copy(to, from)
	return to, err
}

func (repository *UserRepository) DomainModelToModel(from *domain.User) (*models.User, error) {
	to := &models.User{}
	err := copier.Copy(to, from)
	return to, err
}

func NewUserRepository(cache *cache.CachedRepository) domain.UserRepository {
	return &UserRepository{CachedRepository: cache}
}
