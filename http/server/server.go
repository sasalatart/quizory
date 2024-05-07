package server

import (
	"errors"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/sasalatart.com/quizory/answer"
	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/http/oapi"
	"github.com/sasalatart.com/quizory/question"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check.
var _ oapi.ServerInterface = (*Server)(nil)

type Server struct {
	*fiber.App
	cfg             config.ServerConfig
	answerService   *answer.Service
	questionService *question.Service
}

func NewServer(
	cfg config.ServerConfig,
	answerService *answer.Service,
	questionService *question.Service,
) *Server {
	return &Server{
		cfg:             cfg,
		answerService:   answerService,
		questionService: questionService,
	}
}

func (s *Server) Start() {
	s.App = fiber.New()
	applyDefaultMiddleware(s.App, s.cfg)

	oapi.RegisterHandlers(s.App, s)

	addr := s.cfg.Address()
	slog.Info("Server is running", "address", addr)
	if err := s.Listen(addr); err != nil {
		log.Fatal(err)
	}
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
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userID, err := GetUserID(c)
	if err != nil {
		return err
	}

	submissionResponse, err := s.answerService.Submit(ctx, userID, req.ChoiceId)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(toSubmitAnswerResult(*submissionResponse))
}
