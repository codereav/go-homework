package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to load .env"))
	}

	err = run(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			config := NewConfig(configFile)
			log := logger.New(config.Logger.Level, config.Logger.Path)

			var storage app.Storage
			switch config.Database.Type {
			case "sql":
				storage = sqlstorage.New(config.Database.Dsn)
			case "memory":
				storage = memorystorage.New()
			default:
				log.Error("failed to initialize storage: unknown storage type")
				return
			}
			if err := storage.Connect(ctx); err != nil {
				log.Error(errors.Wrap(err, "failed to connect storage").Error())
				return
			}
			defer func(storage app.Storage, ctx context.Context) {
				err := storage.Close(ctx)
				if err != nil {
					log.Error(errors.Wrap(err, "failed to close storage connection").Error())
				}
			}(storage, ctx)

			calendar := app.New(log, storage)

			httpAddr := net.JoinHostPort(config.Server.HTTP.Host, config.Server.HTTP.Port)
			httpServer := internalhttp.NewServer(log, calendar, httpAddr)

			grpcAddr := net.JoinHostPort(config.Server.Grpc.Host, config.Server.Grpc.Port)
			grpcServer := internalgrpc.NewServer(log, calendar, grpcAddr)

			go func() {
				<-ctx.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				if err := httpServer.Stop(ctx); err != nil {
					log.Error(errors.Wrap(err, "failed to stop http server").Error())
				}
				if err := grpcServer.Stop(ctx); err != nil {
					log.Error(errors.Wrap(err, "failed to stop grpc server").Error())
				}
			}()

			log.Info("calendar is running...")

			wg := &sync.WaitGroup{}
			wg.Add(2)
			go func() {
				defer wg.Done()
				if err := httpServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error(errors.Wrap(err, "failed to start http server").Error())
					return
				}
			}()
			go func() {
				defer wg.Done()
				if err := grpcServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error(errors.Wrap(err, "failed to start http server").Error())
					return
				}
			}()
			wg.Wait()
		},
	}
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/calendar/config.yaml", "Path to config file")
	err := rootCmd.PersistentFlags().Parse(os.Args)
	if err != nil {
		return err
	}

	return errors.Wrap(rootCmd.ExecuteContext(ctx), "run application")
}
