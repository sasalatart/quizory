package grpc

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/http/grpc/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server that handles INTERNAL requests.
type Server struct {
	proto.UnimplementedQuizoryServiceServer

	cfg             config.ServerConfig
	grpcServer      *grpc.Server
	questionService *question.Service
}

func NewServer(
	cfg config.ServerConfig,
	questionService *question.Service,
) *Server {
	statsHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
	return &Server{
		cfg:             cfg,
		grpcServer:      grpc.NewServer(grpc.StatsHandler(statsHandler)),
		questionService: questionService,
	}
}

func (s *Server) Start() {
	slog.Info("Running gRPC server", slog.String("address", s.cfg.GRPCAddress()))

	l, err := net.Listen("tcp", s.cfg.GRPCAddress())
	if err != nil {
		slog.Error("Failed to open TCP connection for GRPC", slog.Any("error", err))
		os.Exit(1)
	}

	proto.RegisterQuizoryServiceServer(s.grpcServer, s)

	if err := s.grpcServer.Serve(l); err != nil {
		slog.Error("Failed to start GRPC server", slog.Any("error", err))
		os.Exit(1)
	}
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}

func (s *Server) GetLatestQuestions(
	ctx context.Context,
	req *proto.GetLatestQuestionsRequest,
) (*proto.GetLatestQuestionsResponse, error) {
	parsedTopic, err := enums.TopicString(req.Topic)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid topic: %v", err)
	}

	questions, err := s.questionService.Latest(ctx, parsedTopic, int(req.Amount))
	if err != nil {
		return nil, err
	}
	return &proto.GetLatestQuestionsResponse{
		Questions: toLatestQuestions(questions),
	}, nil
}

func (s *Server) CreateQuestion(
	ctx context.Context,
	req *proto.CreateQuestionRequest,
) (*proto.CreateQuestionResponse, error) {
	topic, err := enums.TopicString(req.Topic)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid topic: %v", err)
	}

	difficulty, err := fromDifficulty(req.Difficulty)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid difficulty: %v", err)
	}

	q := question.
		New(req.Question, req.Hint, req.MoreInfo).
		WithTopic(topic).
		WithDifficulty(*difficulty)
	for _, choice := range req.Choices {
		q.WithChoice(choice.Choice, choice.IsCorrect)
	}
	if err := s.questionService.Insert(ctx, q); err != nil {
		return nil, err
	}
	return &proto.CreateQuestionResponse{
		Id: q.ID.String(),
	}, nil
}
