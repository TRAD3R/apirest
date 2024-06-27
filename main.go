package main

import (
	"context"
	"errors"
	"github.com/trad3r/hskills/apirest/handler"
	"github.com/trad3r/hskills/apirest/router"
	"github.com/trad3r/hskills/apirest/storage"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	// Graceful shutdown
	// TODO: обработку контекста завершения программы
	// docker stop app -> SIGTERM (SIGINT) -> 10s timeout -> SIGKILL
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	us := storage.NewUserStorage()
	ps := storage.NewPostStorage()

	r := router.NewRouter(us, ps)
	h := handler.NewHandler(r)

	log.Printf("Listening on port 8080")

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
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	if err := s.Close(); err != nil {
		log.Printf("Error closing server: %v\n", err)
	}
}
