package testutil

import (
	"github.com/sasalatart.com/quizory/answer"
	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db"
	"github.com/sasalatart.com/quizory/http/server"
	"github.com/sasalatart.com/quizory/llm"
	"github.com/sasalatart.com/quizory/question"
	"go.uber.org/fx"
)

// Module defines a reusable module so that we do not need to manually provide all the dependencies
// in every test suite. It also provides test-specific defaults.
var Module = fx.Module(
	"testutil",

	fx.Provide(config.NewTestConfig),
	db.TestModule,
	llm.TestModule,
	server.TestModule,

	answer.Module,
	question.Module,

	// Repositories are injected privately in the modules above, so we provide them here to make
	// them available for tests (e.g. for seeding the database with test data).
	fx.Provide(answer.NewRepository),
	fx.Provide(question.NewRepository),
)
