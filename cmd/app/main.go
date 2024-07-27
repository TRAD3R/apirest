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
	"github.com/trad3r/hskills/apirest/internal/service"
	"github.com/trad3r/hskills/apirest/internal/storage"
)

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
	defer db.Close()

	if err := migrator.ApplyPostgresMigrations("migrations", cfg.DB.Url); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	u := service.NewUserService(logger, db)
	p := service.NewPostService(logger, db)
	up := service.NewUserPostService(u, p)
	h := handler.NewHandler(u, p, up)

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
