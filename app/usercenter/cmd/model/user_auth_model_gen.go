package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"zero-demo/common/xerr"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userAuthFieldNames          = builder.RawFieldNames(&UserAuth{})
	userAuthRows                = strings.Join(userAuthFieldNames, ",")
	userAuthRowsExpectAutoSet   = strings.Join(stringx.Remove(userAuthFieldNames, "`id`", "`create_time`", "`update_time`", "`delete_time`"), ",")
	userAuthRowsWithPlaceHolder = strings.Join(stringx.Remove(userAuthFieldNames, "`id`", "`create_time`", "`update_time`", "`delete_time`"), "=?,") + "=?"

	cacheUserAuthIdPrefix              = "cache:userAuth:id:"
	cacheUserAuthAuthTypeAuthKeyPrefix = "cache:userAuth:authType:authKey:"
	cacheUserAuthUserIdAuthTypePrefix  = "cache:userAuth:userId:authType:"
)

type (
	UserAuthModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *UserAuth) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*UserAuth, error)
		FindOneByAuthTypeAuthKey(ctx context.Context, authType string, authKey string) (*UserAuth, error)
		FindOneByUserIdAuthType(ctx context.Context, userId int64, authType string) (*UserAuth, error)
		Update(ctx context.Context, session sqlx.Session, data *UserAuth) error
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *UserAuth) error
		Delete(ctx context.Context, id int64) error
	}

	defaultUserAuthModel struct {
		sqlc.CachedConn
		table string
	}

	UserAuth struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		Version    int64     `db:"version"` // 版本号
		UserId     int64     `db:"user_id"`
		AuthKey    string    `db:"auth_key"`  // 平台唯一id
		AuthType   string    `db:"auth_type"` // 平台类型
	}
)

func NewUserAuthModel(conn sqlx.SqlConn, c cache.CacheConf) UserAuthModel {
	return &defaultUserAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_auth`",
	}
}

func (m *defaultUserAuthModel) Insert(ctx context.Context, session sqlx.Session, data *UserAuth) (sql.Result, error) {
	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, data.Id)
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userAuthRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType)
		}
		return conn.ExecCtx(ctx, query, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType)
	}, userAuthIdKey, userAuthAuthTypeAuthKeyKey, userAuthUserIdAuthTypeKey)
	return ret, err
}

func (m *defaultUserAuthModel) FindOne(ctx context.Context, id int64) (*UserAuth, error) {
	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, id)
	var resp UserAuth
	err := m.QueryRowCtx(ctx, &resp, userAuthIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByAuthTypeAuthKey(ctx context.Context, authType string, authKey string) (*UserAuth, error) {
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, authType, authKey)
	var resp UserAuth
	err := m.QueryRowIndexCtx(ctx, &resp, userAuthAuthTypeAuthKeyKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_type` = ? and `auth_key` = ? limit 1", userAuthRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, authType, authKey); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByUserIdAuthType(ctx context.Context, userId int64, authType string) (*UserAuth, error) {
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, userId, authType)
	var resp UserAuth
	err := m.QueryRowIndexCtx(ctx, &resp, userAuthUserIdAuthTypeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `auth_type` = ? limit 1", userAuthRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, userId, authType); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) Update(ctx context.Context, session sqlx.Session, data *UserAuth) error {
	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, data.Id)
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id)
	}, userAuthIdKey, userAuthAuthTypeAuthKeyKey, userAuthUserIdAuthTypeKey)
	return err
}

func (m *defaultUserAuthModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *UserAuth) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	looklookUsercenterUserAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, data.Id)
	looklookUsercenterUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookUsercenterUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, userAuthRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id, oldVersion)
	}, looklookUsercenterUserAuthUserIdAuthTypeKey, looklookUsercenterUserAuthIdKey, looklookUsercenterUserAuthAuthTypeAuthKeyKey)
	if err != nil {
		return err
	}
	updateCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if updateCount == 0 {
		return xerr.NewErrCode(xerr.DB_UPDATE_AFFECTED_ZERO_ERROR)
	}

	return nil
}

func (m *defaultUserAuthModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, id)
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, userAuthAuthTypeAuthKeyKey, userAuthUserIdAuthTypeKey, userAuthIdKey)
	return err
}

func (m *defaultUserAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, primary)
}

func (m *defaultUserAuthModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}
