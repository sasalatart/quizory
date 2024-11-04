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

// RecordSuccessfulGeneration records a successful generation of questions, along with the time it
// took for the whole process.
func (s Service) RecordSuccessfulGeneration(ctx context.Context, d time.Duration) {
	s.successfulGenerationCounter.Add(ctx, 1)
	s.generationDuration.Record(ctx, d.Milliseconds())
}

// RecordFailedValidations records cases where one or more LLM-generated questions have failed their
// validations and therefore need to be discarded.
func (s Service) RecordFailedValidations(ctx context.Context, amount int64) {
	s.failedValidationCounter.Add(ctx, amount)
}
