// +build wireinject

package main

import (
	"database/sql"

	"github.com/Promacanthus/Go-000/Week04/pkg/biz"
	"github.com/Promacanthus/Go-000/Week04/pkg/dao"
	"github.com/Promacanthus/Go-000/Week04/pkg/service"
	"github.com/google/wire"
)

// wire.go

func InitHandler(db *sql.DB) service.Handler {
	wire.Build(service.NewHandler, biz.NewStringService, dao.NewRepo)
	return &service.HandlerImp{}
}
