package main

import (
	"github.com/gsasso/go-backend/src/server/internal/controller"
	server "github.com/gsasso/go-backend/src/server/internal/server"
	"github.com/gsasso/go-backend/src/server/internal/ticker"
)

func main() {
	logisticCtlr := controller.NewLogisticController(&ticker.SummaryService{})
	logisticServer := server.RunGRPCServer(logisticCtlr)
	logisticServer.Start()

	// TODO General recommendations: 1. Try keep code clean and self documented (proper names for variables, functions, structures)
	// TODO General recommendations: 2. Keep same code style even if it's not production code - it's your habit and general attitude for work
	// TODO General recommendations: 3. Periodically ask this question when reading own code: "Is it doing 1 thing?"
	// TODO General recommendations: 4. Try to notice how often you need return to code while following context (Abstraction complexity and how easy to follow flow?)
	// TODO General recommendations: 5. I was expecting to see latest 1 unit test, manually tested code is nice but core logic may change how it's protected?
}
