package server

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/pagination"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/http/oapi"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check.
var _ oapi.ServerInterface = (*Server)(nil)

type Server struct {
	*fiber.App
	cfg             config.ServerConfig
	db              *sql.DB
	answerService   *answer.Service
	questionService *question.Service
}

func NewServer(
	cfg config.ServerConfig,
	db *sql.DB,
	answerService *answer.Service,
	questionService *question.Service,
) *Server {
	return &Server{
		cfg:             cfg,
		db:              db,
		answerService:   answerService,
		questionService: questionService,
	}
}

func (s *Server) Start() {
	s.App = fiber.New()

	s.registerMiddlewares()
	if err := s.registerAppHandlers(); err != nil {
		log.Fatal(err)
	}

	addr := s.cfg.Address()
	slog.Info("Server is running", "address", addr)
	if err := s.Listen(addr); err != nil {
		log.Fatal(err)
	}
}

// HealthCheck returns a 204 status code if the server is healthy, and a 503 status code otherwise.
func (s *Server) HealthCheck(c *fiber.Ctx) error {
	// TODO: add a more meaningful health check. In the meantime, it is just checking if the
	// database connection is healthy.
	if s.db.Ping() != nil {
		return fiber.ErrServiceUnavailable
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetNextQuestion returns the next question for a user to answer.
func (s *Server) GetNextQuestion(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	q, err := s.questionService.NextFor(ctx, userID)
	if err != nil && errors.Is(err, question.ErrNoQuestionsLeft) {
		return c.SendStatus(fiber.StatusNoContent)
	}
	if err != nil {
		slog.Error("Failed to get next question", "error", err)
		return err
	}
	return c.Status(fiber.StatusOK).JSON(toUnansweredQuestion(*q))
}

// SubmitAnswer registers the choice made by a user for a specific question, and returns the correct
// choice for it, plus some more info for the user to know how they did.
func (s *Server) SubmitAnswer(c *fiber.Ctx) error {
	ctx := c.Context()

	req := new(oapi.SubmitAnswerRequest)
	if err := c.BodyParser(req); err != nil {
		slog.Error("Failed to parse request body", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	submissionResponse, err := s.answerService.Submit(ctx, userID, req.ChoiceId)
	if err != nil {
		slog.Error("Failed to submit answer", "error", err)
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(toSubmitAnswerResult(*submissionResponse))
}

// GetAnswersLog returns a list of previous attempts at answering questions from the specified user.
func (s *Server) GetAnswersLog(
	c *fiber.Ctx,
	userID uuid.UUID,
	params oapi.GetAnswersLogParams,
) error {
	ctx := c.Context()

	p := pagination.New(params.Page, params.PageSize)
	if err := p.Validate(); err != nil {
		slog.Error("Invalid pagination", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid pagination")
	}

	logItems, err := s.answerService.LogFor(ctx, answer.LogRequest{
		UserID:     userID,
		Pagination: p,
	})
	if err != nil {
		slog.Error("Failed to get answers log", "error", err)
		return err
	}

	result := make([]oapi.AnswersLogItem, 0, len(logItems))
	for _, logItem := range logItems {
		result = append(result, toAnswersLogItem(logItem))
	}
	return c.Status(fiber.StatusOK).JSON(result)
}
