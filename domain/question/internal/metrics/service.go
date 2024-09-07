package metrics

import (
	"context"
	"log"
	"time"

	"github.com/sasalatart/quizory/infra/otel"
	"go.opentelemetry.io/otel/metric"
)

type Service struct {
	successfulGenerationCounter metric.Int64Counter
	failedValidationCounter     metric.Int64Counter
	generationDuration          metric.Int64Histogram
}

func NewService(meter otel.Meter) Service {
	successfulGenerationCounter, err := meter.Int64Counter("questions_generation_success")
	if err != nil {
		log.Fatal("unable to create questions_generation_success counter")
	}

	failedValidationCounter, err := meter.Int64Counter("questions_generation_validation_failure")
	if err != nil {
		log.Fatal("unable to create questions_generation_validation_failure counter")
	}

	generationDuration, err := meter.Int64Histogram(
		"questions_generation_duration",
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatal("unable to create questions_generation_duration histogram")
	}

	return Service{
		successfulGenerationCounter: successfulGenerationCounter,
		failedValidationCounter:     failedValidationCounter,
		generationDuration:          generationDuration,
	}
}

// OnSuccessfulGeneration records the count of all LLM calls that successfully complete and result
// in a full set of questions generated (i.e. no questions are discarded).
func (s Service) OnSuccessfulGeneration(ctx context.Context) {
	s.successfulGenerationCounter.Add(ctx, 1)
}

// OnFailedValidation records cases where an LLM-generated question has failed its validation and
// therefore needs to be discarded.
func (s Service) OnFailedValidation(ctx context.Context) {
	s.failedValidationCounter.Add(ctx, 1)
}

// OnLLMCallFinished records the duration of an LLM call.
func (s Service) OnLLMCallFinished(ctx context.Context, d time.Duration) {
	s.generationDuration.Record(ctx, d.Milliseconds())
}
