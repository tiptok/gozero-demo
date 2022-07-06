package repository

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/tiptok/gocomm/pkg/cache"
	"gorm.io/gorm"
	"zero-demo/app/usercenter/internal/pkg/db/models"
	"zero-demo/app/usercenter/internal/pkg/db/transaction"
	"zero-demo/app/usercenter/internal/pkg/domain"
)

type UserAuthRepository struct {
	*cache.CachedRepository
}

func (repository *UserAuthRepository) Insert(ctx context.Context, transaction transaction.Conn, dm *domain.UserAuth) (*domain.UserAuth, error) {
	var (
		err error
		m   = &models.UserAuth{}
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

func (repository *UserAuthRepository) Update(ctx context.Context, transaction transaction.Conn, dm *domain.UserAuth) (*domain.UserAuth, error) {
	var (
		err error
		m   *models.UserAuth
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

func (repository *UserAuthRepository) Delete(ctx context.Context, transaction transaction.Conn, dm *domain.UserAuth) (*domain.UserAuth, error) {
	var (
		tx = transaction.DB()
		m  = &models.UserAuth{Id: dm.Identify().(int64)}
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

func (repository *UserAuthRepository) FindOne(ctx context.Context, transaction transaction.Conn, id int64) (*domain.UserAuth, error) {
	var (
		err error
		tx  = transaction.DB()
		m   = new(models.UserAuth)
	)
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(m).Where("id = ?", id).Find(m)
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return m, tx.Error
	}
	cacheModel := new(models.UserAuth)
	cacheModel.Id = id
	if err = repository.QueryCache(cacheModel.CacheKeyFunc, m, queryFunc); err != nil {
		return nil, err
	}
	return repository.ModelToDomainModel(m)
}

func (repository *UserAuthRepository) Find(ctx context.Context, transaction transaction.Conn, queryOptions map[string]interface{}) (int64, []*domain.UserAuth, error) {
	var (
		tx    = transaction.DB()
		ms    []*models.UserAuth
		dms   = make([]*domain.UserAuth, 0)
		total int64
	)
	queryFunc := func() (interface{}, error) {
		tx = tx.Model(&ms).Order("id desc")
		if v, ok := queryOptions["offset"]; ok {
			tx.Offset(v.(int))
		}
		if v, ok := queryOptions["limit"]; ok {
			tx.Limit(v.(int))
		}
		if tx = tx.Find(&ms); tx.Error != nil {
			return dms, tx.Error
		}
		tx.Count(&total)
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

func (repository *UserAuthRepository) ModelToDomainModel(from *models.UserAuth) (*domain.UserAuth, error) {
	to := &domain.UserAuth{}
	err := copier.Copy(to, from)
	return to, err
}

func (repository *UserAuthRepository) DomainModelToModel(from *domain.UserAuth) (*models.UserAuth, error) {
	to := &models.UserAuth{}
	err := copier.Copy(to, from)
	return to, err
}

func NewUserAuthRepository(cache *cache.CachedRepository) domain.UserAuthRepository {
	return &UserAuthRepository{CachedRepository: cache}
}
