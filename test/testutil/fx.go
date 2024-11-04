package testutil

import (
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db/dbtest"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/generator"
	"github.com/sasalatart/quizory/http/grpc"
	grpclient "github.com/sasalatart/quizory/http/grpc/client"
	"github.com/sasalatart/quizory/http/rest/resttest"
	"github.com/sasalatart/quizory/infra/otel/oteltest"
	"github.com/sasalatart/quizory/llm/llmtest"
	"go.uber.org/fx"
)

// Module defines a reusable module so that we do not need to manually provide all the dependencies
// in every test suite. It also provides test-specific defaults.
// Module DOES NOT include servertest.Module. Use ModuleWithAPI for that instead.
var Module = fx.Module(
	"testutil",

	fx.Provide(config.NewTestConfig),
	oteltest.Module,
	dbtest.Module,
	llmtest.Module,

	answer.Module,
	question.Module,
	generator.Module,

	// Repositories are injected privately in the modules above, so we provide them here to make
	// them available for tests (e.g. for seeding the database with test data).
	fx.Provide(answer.NewRepository),
	fx.Provide(question.NewRepository),
)

// ModuleWithAPI injects the API module in addition to the dependencies provided by Module.
// It is intended to be used in test suites that require the API module, as it also turns on the API
// server, manages its lifecycle, and waits for it to be ready before running tests.
var ModuleWithAPI = fx.Module(
	"testutil-with-api",

	Module,
	grpc.Module,
	resttest.Module,
	grpclient.Module,
)
