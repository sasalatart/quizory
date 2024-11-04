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

// RecordSuccessfulGeneration records the count of all LLM calls that successfully complete and
// result in a full set of questions generated (i.e. no questions are discarded).
func (s Service) RecordSuccessfulGeneration(ctx context.Context) {
	s.successfulGenerationCounter.Add(ctx, 1)
}

// RecordFailedValidations records cases where one or more LLM-generated questions have failed their
// validations and therefore need to be discarded.
func (s Service) RecordFailedValidations(ctx context.Context, amount int64) {
	s.failedValidationCounter.Add(ctx, amount)
}

// RecordGenerationDuration records the duration of the whole questions generation process,
// regardless of whether it had validation issues or not.
func (s Service) RecordGenerationDuration(ctx context.Context, d time.Duration) {
	s.generationDuration.Record(ctx, d.Milliseconds())
}
