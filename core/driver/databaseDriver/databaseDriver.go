package databaseDriver

import (
	"context"
	"fmt"

	"example/core/config"
	"example/core/helper/enviromentHelper"

	"github.com/go-pg/pg/v10"
)

type DatabaseDriver interface {
	GetDB() *pg.DB
}

type databaseDriver struct {
	db *pg.DB
}

type dbLogger struct{}

func New(config config.Config) DatabaseDriver {
	var db *pg.DB
	dbConfig := config.GetDatabase()

	db = pg.Connect(
		&pg.Options{
			Addr:     dbConfig.Addr,
			User:     dbConfig.User,
			Password: dbConfig.Password,
			Database: dbConfig.Dbname,
		},
	)

	if enviromentHelper.IsDevelopment() {
		db.AddQueryHook(dbLogger{})
	}

	// if err := createSchema(db); err != nil {
	// 	log.Fatalln(err)
	// }

	return databaseDriver{
		db,
	}
}

func (d databaseDriver) GetDB() *pg.DB {
	return d.db
}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	bytes, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
