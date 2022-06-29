package db

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"zero-demo/app/usercenter/pkg1/db/models"
)

func Migrate(DB *pg.DB) {
	for _, model := range []interface{}{
		(*models.Users)(nil),
	} {
		err := DB.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			panic(err)
		}
	}
}
