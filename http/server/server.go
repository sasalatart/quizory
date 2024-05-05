package server

import (
	"errors"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/http/oapi"
	"github.com/sasalatart.com/quizory/question"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check.
var _ oapi.ServerInterface = (*Server)(nil)

type Server struct {
	*fiber.App
	cfg             config.ServerConfig
	questionService *question.Service
}

func NewServer(cfg config.ServerConfig, questionService *question.Service) *Server {
	return &Server{
		cfg:             cfg,
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

// GetQuestionsNext returns the next question for a user to answer.
func (s *Server) GetQuestionsNext(c *fiber.Ctx) error {
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
