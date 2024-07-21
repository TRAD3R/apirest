package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/TRAD3R/tlog"
	"github.com/trad3r/hskills/apirest/internal/config"
	"github.com/trad3r/hskills/apirest/internal/handler"
	"github.com/trad3r/hskills/apirest/internal/migrator"
	"github.com/trad3r/hskills/apirest/internal/router"
	"github.com/trad3r/hskills/apirest/internal/storage"
)

/*
*
Сделать:
Реализовать сервис-блог, выполняющий хранение двух моделей данных: Пользователи, Посты.

Пользователь:
ID, Имя, Телефон, Время создания, Время последнего редактирования

Посты:
ID
Subject
Время создания
ID пользователя, написавшего пост
Содержание

В качестве хранения использовать map. Предоставить интерфейс доступа к данным с помощью HTTP API, использовать net/http. Методы доступа к данным:

GET /users
Query:
fromCreatedAt - фильтр на левую границу времени создания
toCreatedAt - фильтр на правую границу времени создания
name - множественный фильтр на имя пользователя
limit - кол-во записей к выдаче
offset - кол-во записей к пропуску
topPostsAmount - при передачи в этом параметре asc/desc - данные в выборке отсортированы по кол-ву написанных пользователем постов

POST /user - запрос на создание пользователя
PATCH /user/:id - запрос на изменение пользователя позволяющий принять одно из двух полей (имя, телефон) или сразу два эти поля и выполнить изменение пользователя. При изменении пользователя обновляется также время последнего редактирования

DELETE /user/:id - удаление пользователя

Методы для /posts оставляю на твою фантазию))
*/
func main() {
	cfg := config.GetConfig()

	logger := tlog.GetLogger(cfg.IsDebug)

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	runtime.SetMutexProfileFraction(1)

	db, err := storage.NewDB(ctx, cfg.DB.Url)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer storage.Stop()

	if err := migrator.ApplyPostgresMigrations("migrations", cfg.DB.Url); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	r := router.NewRouter(logger, db)
	h := handler.NewHandler(r)

	logger.Info("listening on port 8080")

	// Review:
	s := http.Server{
		Addr:              ":8080",
		Handler:           h.Handlers(),
		ReadTimeout:       time.Second * 3,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    1e6,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	if err := s.Close(); err != nil {
		logger.Error("error closing server", "err", err.Error())
	}
}
