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
	"github.com/faissalmaulana/21/api/internal/service"
	"github.com/go-playground/validator/v10"
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
			fx.Annotate(handler.NewPostProjectHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"postProjectHandler"`)),
			fx.Annotate(handler.NewDeleteProjectHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"deleteProjectHandler"`)),
			fx.Annotate(handler.NewUpdateProjectHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"updateProjectHandler"`)),
			// TASKS HANDLERS
			fx.Annotate(handler.NewPostTaskHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"postTaskHandler"`)),
			fx.Annotate(handler.NewUpdateTaskHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"updateTaskHandler"`)),
			fx.Annotate(handler.NewGetTaskByIDHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"getTaskByIDHandler"`)),
			fx.Annotate(handler.NewGetTasksHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"getTasksHandler"`)),
			fx.Annotate(handler.NewDeleteTaskByIDHandler, fx.As(new(handler.Handle)), fx.ResultTags(`name:"deleteTaskByIDHandler"`)),
			zap.NewDevelopment,
			validator.New,
			service.NewSugaredErrorMessageValidator,
			fx.Annotate(repository.NewProject, fx.As(new(repository.ProjectRepository))),
			fx.Annotate(repository.NewTask, fx.As(new(repository.TaskRepository))),
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
