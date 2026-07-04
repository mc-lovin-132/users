package app

import (
	"context"
	"net"

	"github.com/mc-lovin-132/users/config"
	"github.com/mc-lovin-132/users/internal/infrastructure/delivery/handlers"
	"github.com/mc-lovin-132/users/internal/infrastructure/delivery/interceptors"
	"github.com/mc-lovin-132/users/internal/infrastructure/repository"
	"github.com/mc-lovin-132/users/internal/service"
	"github.com/mc-lovin-132/users/pb"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Start(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", a.cfg.DSN())
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			a.logger.Error("err while closing db connection", zap.Error(err))
			return
		}
		a.logger.Info("db connection successfuly closed")
	}()
	a.logger.Info("successfuly init db connection")

	err = repository.RunMigrations(db.DB)
	if err != nil {
		return err
	}

	repo := repository.New(db)
	srvc := service.New(repo)
	handler := handlers.New(srvc)

	healthService := service.NewHealthService(db)
	healthHandler := handlers.NewHealthHandler(healthService)

	lis, err := net.Listen("tcp", a.cfg.Addr())
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.NewLoggingInterceptor(a.logger)),
	)

	pb.RegisterHealthServiceServer(grpcServer, healthHandler)
	pb.RegisterUserServiceServer(grpcServer, handler)

	errG, gCtx := errgroup.WithContext(ctx)

	errG.Go(func() error {
		<-gCtx.Done()
		grpcServer.GracefulStop()
		return gCtx.Err()
	})

	errG.Go(func() error {
		a.logger.Info("gRPC server listening", zap.String("port", a.cfg.Port))
		if err := grpcServer.Serve(lis); err != nil {
			a.logger.Error("grpc server run failed", zap.Error(err))
			return err
		}
		return nil
	})

	return errG.Wait()
}
