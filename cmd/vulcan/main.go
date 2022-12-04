package main

import (
	"context"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	mgr "github.com/exfly/vulcan/db/migration"
	"github.com/exfly/vulcan/internel/config"
	userrpc "github.com/exfly/vulcan/internel/user/rpc"
	vulcanv1 "github.com/exfly/vulcan/pb"

	// for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	app, err := NewApp(".")
	if err != nil {
		panic(err)
	}

	runDBMigration(app.cfg)

	go func() {
		log.WithField("addr", app.httpListener.Addr().String()).Info("start http server")
		err := http.Serve(app.httpListener, app.httpHandler)
		if err != nil {
			log.WithError(err).Error("cannot start HTTP gateway server")
		}
	}()

	log.WithField("addr", app.grpcListener.Addr().String()).Info("start gRPC server")
	err = app.grpcServer.Serve(app.grpcListener)
	if err != nil {
		log.WithError(err).Fatal("cannot start gRPC server")
	}
}

type App struct {
	ctx          context.Context
	grpcListener net.Listener
	grpcServer   *grpc.Server
	httpListener net.Listener
	httpHandler  http.Handler

	cfg *config.Config
}

func newApp(
	cfg *config.Config,
	server vulcanv1.SimpleBankServer,
) (*App, error) {
	ctx := context.Background()
	ret := &App{
		ctx: ctx,
		cfg: cfg,
	}
	var (
		err error
	)

	logEntry := log.NewEntry(log.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logEntry),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logEntry),
		)),
	)
	vulcanv1.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	ret.grpcServer = grpcServer

	grpcListener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create listener")
	}
	ret.grpcListener = grpcListener

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	err = vulcanv1.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		return nil, errors.Wrap(err, "cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	httpListener, err := net.Listen("tcp", cfg.HTTPServerAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create listener")
	}
	ret.httpListener = httpListener

	handler := userrpc.HTTPLogger(mux)
	ret.httpHandler = handler

	return ret, nil
}

func runDBMigration(cfg *config.Config) {
	migration, err := mgr.NewMigration(cfg.DBSource)
	if err != nil {
		log.WithError(err).Fatal("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.WithError(err).Fatal("failed to run migrate up")
	}

	log.Info("db migrated successfully")
}
