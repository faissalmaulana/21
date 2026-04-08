package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/faissalmaulana/21/api/cmd/handler"
	"github.com/faissalmaulana/21/api/internal/db"
	"github.com/faissalmaulana/21/api/internal/repository"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	dbport, err := strconv.Atoi(os.Getenv("PSQL_DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	fx.New(
		fx.Provide(
			NewEchoMux,
			NewHttpServer,
			fx.Annotate(handler.NewPingHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"pingHandler"`)),
			fx.Annotate(handler.NewGetProjectsHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"getProjectsHandler"`)),
			zap.NewDevelopment,
			fx.Annotate(repository.NewProject, fx.As(new(repository.ProjectRepository))),
			func(lc fx.Lifecycle, log *zap.Logger) (*sql.DB, error) {
				pg := db.NewPostgresqlDB(
					dbport,
					os.Getenv("PSQL_DB_HOST"),
					os.Getenv("PSQL_DB_USER"),
					os.Getenv("PSQL_DB_CREDENTIAL"),
					os.Getenv("PSQL_DB_NAME"),
				)
				return pg.Connect(lc, log)
			},
		),
		fx.Invoke(func(_ *sql.DB) {}),
		fx.Invoke(func(*http.Server) {}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
