package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Marattttt/coderunner/langrunner/internal/config"
	"github.com/Marattttt/coderunner/langrunner/internal/runner"
)

func main() {
	appCtx := context.TODO()

	conf, err := config.Config(appCtx)
	checkFatal(err, "Reading config from env")

	appLogger := conf.CreateLogger()

	manager := runner.NewRuntimeManager(&conf.RunnerConig, appLogger)
	runner := runner.GoRunner{Conf: *conf.Go, Manager: manager, Logger: appLogger}

	res, err := runner.RunCode(appCtx, []byte(`
	package main
	import "fmt"
	func main() {
		fmt.Println("Hello world")
	}
	`))

	checkFatal(err, "Running code")

	fmt.Println(res)
}

func checkFatal(err error, msg string) {
	if err == nil {
		return
	}

	slog.Error(msg, slog.String("err", err.Error()))
	os.Exit(1)
}
