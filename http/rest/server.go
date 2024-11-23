package rest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/pagination"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/http/rest/middleware"
	"github.com/sasalatart/quizory/http/rest/oapi"
	"github.com/sasalatart/quizory/infra/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check.
var _ oapi.ServerInterface = (*Server)(nil)

// Server represents the REST server that handles external requests (e.g. from browsers).
type Server struct {
	httpServer      http.Server
	cfg             config.ServerConfig
	db              *sql.DB
	answerService   *answer.Service
	questionService *question.Service
	meter           otel.Meter
}

func NewServer(
	cfg config.ServerConfig,
	db *sql.DB,
	answerService *answer.Service,
	questionService *question.Service,
	meter otel.Meter,
) *Server {
	return &Server{
		cfg:             cfg,
		db:              db,
		answerService:   answerService,
		questionService: questionService,
		meter:           meter,
	}
}

func (s *Server) Start() {
	slog.Info("Running REST server", slog.String("address", s.cfg.RESTAddress()))

	mux := http.NewServeMux()
	handler := oapi.HandlerWithOptions(s, oapi.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []oapi.MiddlewareFunc{
			middleware.WithJWTAuth(s.cfg.JWTSecret, []string{"/openapi", "/health-check"}),
			middleware.WithRecover,
			middleware.WithMonitoring(s.meter),
			otelhttp.NewMiddleware("quizory-server"),
		},
	})
	if err := registerSwaggerHandlers(mux, s.cfg.OAPISchemaDir()); err != nil {
		slog.Error("Failed to register Swagger handlers", slog.Any("error", err))
		os.Exit(1)
	}

	s.httpServer = http.Server{
		Addr:         s.cfg.RESTAddress(),
		Handler:      middleware.WithCORS(handler),
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.Any("error", err))
		return err
	}
	return nil
}

// HealthCheck returns a 204 status code if the server is healthy, and a 503 status code otherwise.
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: add a more meaningful health check. In the meantime, it is just checking if the
	// database connection is healthy.
	if err := s.db.Ping(); err != nil {
		http.Error(w, "Database connection is unhealthy", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetRemainingTopics returns the list of topics with questions still unanswered by the user making
// the request. Each of these topics comes with the actual amount of questions left to answer.
func (s *Server) GetRemainingTopics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	remainingTopics, err := s.questionService.RemainingTopicsFor(ctx, userID)
	if err != nil {
		handleServerError(w, "Failed to get remaining topics", err)
		return
	}
	encodeJSON(w, toRemainingTopics(remainingTopics))
}

// GetNextQuestion returns the next question that a user should answer for the specified topic.
func (s *Server) GetNextQuestion(
	w http.ResponseWriter,
	r *http.Request,
	params oapi.GetNextQuestionParams,
) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)

	topic, err := enums.TopicString(params.Topic)
	if err != nil {
		handleBadRequest(w, "Invalid topic", err)
		return
	}

	q, err := s.questionService.NextFor(ctx, userID, topic)
	if errors.Is(err, question.ErrNoQuestionsLeft) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		handleServerError(w, "Failed to get next question", err)
		return
	}
	encodeJSON(w, toUnansweredQuestion(*q))
}

// SubmitAnswer registers the choice made by a user for a specific question, and returns the correct
// choice for it, plus some more info for the user to know how they did.
func (s *Server) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(oapi.SubmitAnswerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handleBadRequest(w, "Failed to parse request body", err)
		return
	}

	userID := middleware.GetUserID(ctx)
	submissionResponse, err := s.answerService.Submit(ctx, userID, req.ChoiceId)
	if err != nil {
		handleServerError(w, "Failed to submit answer", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encodeJSON(w, toSubmitAnswerResult(*submissionResponse))
}

// GetAnswersLog returns a list of previous attempts at answering questions from the specified user.
func (s *Server) GetAnswersLog(
	w http.ResponseWriter,
	r *http.Request,
	userID uuid.UUID,
	params oapi.GetAnswersLogParams,
) {
	ctx := r.Context()

	p := pagination.New(params.Page, params.PageSize)
	if err := p.Validate(); err != nil {
		handleBadRequest(w, "Invalid pagination", err)
		return
	}

	logItems, err := s.answerService.LogFor(ctx, answer.LogRequest{
		UserID:     userID,
		Pagination: p,
	})
	if err != nil {
		handleServerError(w, "Failed to get answers log", err)
		return
	}

	result := make([]oapi.AnswersLogItem, 0, len(logItems))
	for _, logItem := range logItems {
		result = append(result, toAnswersLogItem(logItem))
	}
	encodeJSON(w, result)
}

// encodeJSON encodes the given data as JSON into the response writer. If an error occurs, it logs
// it and sends a 500 status code.
func encodeJSON(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		handleServerError(w, "Failed to encode response", err)
	}
}

func handleBadRequest(w http.ResponseWriter, msg string, err error) {
	slog.Error(msg, slog.Any("error", err))
	http.Error(w, msg, http.StatusBadRequest)
}

func handleServerError(w http.ResponseWriter, msg string, err error) {
	slog.Error(msg, slog.Any("error", err))
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}
