package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example/core/config"
	"example/core/driver/cacheDriver"
	"example/core/driver/databaseDriver"
	"example/core/helper/enviromentHelper"

	"github.com/go-pg/pg/v10/orm"

	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type App interface {
	GetDB() *pg.DB
	GetCache() *redis.Client
	GetRouter() *mux.Router
	GetCtx() context.Context
	GetConfig() config.Config
	RegisterModule(name string, moduleFunc func(m Module))
	CreateSchema([]interface{})
	SetContainer(container interface{})
	Start()
}

type modules = map[string]Module

type app struct {
	config    config.Config
	db        *pg.DB
	cache     *redis.Client
	router    *mux.Router
	ctx       context.Context
	container interface{}
	modules
}

func NewApp() App {
	c := config.New()
	db := databaseDriver.New(c).GetDB()
	cache := cacheDriver.New(c).GetCache()
	ctx := context.Background()

	return &app{
		config:  c,
		db:      db,
		cache:   cache,
		router:  mux.NewRouter(),
		ctx:     ctx,
		modules: make(modules),
	}
}

func (a app) GetDB() *pg.DB {
	return a.db
}

func (a app) GetCache() *redis.Client {
	return a.cache
}

func (a app) GetRouter() *mux.Router {
	return a.router
}

func (a app) GetCtx() context.Context {
	return a.ctx
}

func (a app) GetConfig() config.Config {
	return a.config
}

func (a *app) SetContainer(container interface{}) {
	a.container = container
}

func (a app) CreateSchema(models []interface{}) {
	for _, model := range models {
		err := a.db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
			Temp:          enviromentHelper.IsTest(),
		})
		if err != nil {
			log.Fatalln("create db schema: ", err)
		}
	}
}

func (a *app) RegisterModule(name string, moduleFunc func(m Module)) {
	m := newModule(a, name)
	moduleFunc(m)
	a.modules[name] = m
}

func (a *app) Start() {
	a.router.Use(corsMiddleware(), secureMiddleware())
	//a.router.Use(corsMiddleware(), secureMiddleware(), rateLimitMiddleware())
	a.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ps := newProvideService(a, w, r)

		defer ps.handleError()

		handler := newHandler(ps, a.container, a.modules)

		if ok := handler.use(
			func() bool {
				return handler.setAction()
			},
			func() bool {
				return handler.protectAny()
			},
			func() bool {
				return handler.protectAdmin()
			},
			func() bool {
				return handler.protectClient()
			},
			func() bool {
				return handler.validateForm()
			},
		); !ok {
			return
		}

		handler.resolve()
	}).Methods("POST")

	a.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", a.config.GetApp().Port), a.router))
}
