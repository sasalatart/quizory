package testutil

import (
	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db"
	"github.com/sasalatart.com/quizory/llm"
	"go.uber.org/fx"
)

// Module defines the Fx module that provides the necessary dependencies for testing purposes.
// These include the configuration and the test database, plus some sensible defaults.
var Module = fx.Module(
	"testutil",
	fx.Provide(config.NewConfig),
	db.TestModule,
	llm.TestModule,
)
