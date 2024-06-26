//go:build go1.22

// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package oapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for Difficulty.
const (
	DifficultyAvidHistorian   Difficulty = "DifficultyAvidHistorian"
	DifficultyHistoryScholar  Difficulty = "DifficultyHistoryScholar"
	DifficultyNoviceHistorian Difficulty = "DifficultyNoviceHistorian"
)

// AnswersLogItem A past attempt from a user to answer a question.
type AnswersLogItem struct {
	ChoiceId UUID     `json:"choiceId"`
	Id       UUID     `json:"id"`
	Question Question `json:"question"`
}

// Choice defines model for Choice.
type Choice struct {
	Choice    string `json:"choice"`
	Id        UUID   `json:"id"`
	IsCorrect bool   `json:"isCorrect"`
}

// Difficulty defines model for Difficulty.
type Difficulty string

// Question defines model for Question.
type Question struct {
	Choices    []Choice   `json:"choices"`
	Difficulty Difficulty `json:"difficulty"`
	Hint       string     `json:"hint"`
	Id         UUID       `json:"id"`
	MoreInfo   string     `json:"moreInfo"`
	Question   string     `json:"question"`
	Topic      string     `json:"topic"`
}

// RemainingTopic defines model for RemainingTopic.
type RemainingTopic struct {
	AmountOfQuestions int    `json:"amountOfQuestions"`
	Topic             string `json:"topic"`
}

// SubmitAnswerRequest defines model for SubmitAnswerRequest.
type SubmitAnswerRequest struct {
	ChoiceId UUID `json:"choiceId"`
}

// SubmitAnswerResult defines model for SubmitAnswerResult.
type SubmitAnswerResult struct {
	CorrectChoiceId UUID   `json:"correctChoiceId"`
	Id              UUID   `json:"id"`
	MoreInfo        string `json:"moreInfo"`
}

// UUID defines model for UUID.
type UUID = openapi_types.UUID

// UnansweredChoice defines model for UnansweredChoice.
type UnansweredChoice struct {
	Choice string `json:"choice"`
	Id     UUID   `json:"id"`
}

// UnansweredQuestion defines model for UnansweredQuestion.
type UnansweredQuestion struct {
	Choices    []UnansweredChoice `json:"choices"`
	Difficulty Difficulty         `json:"difficulty"`
	Hint       string             `json:"hint"`
	Id         UUID               `json:"id"`
	Question   string             `json:"question"`
	Topic      string             `json:"topic"`
}

// GetNextQuestionParams defines parameters for GetNextQuestion.
type GetNextQuestionParams struct {
	// Topic The topic for which the next question should be retrieved.
	Topic string `form:"topic" json:"topic"`
}

// GetAnswersLogParams defines parameters for GetAnswersLog.
type GetAnswersLogParams struct {
	// Page The page number to retrieve (index starts at 0).
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// PageSize The number of items per page.
	PageSize *int `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}

// SubmitAnswerJSONRequestBody defines body for SubmitAnswer for application/json ContentType.
type SubmitAnswerJSONRequestBody = SubmitAnswerRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// SubmitAnswerWithBody request with any body
	SubmitAnswerWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SubmitAnswer(ctx context.Context, body SubmitAnswerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// HealthCheck request
	HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetNextQuestion request
	GetNextQuestion(ctx context.Context, params *GetNextQuestionParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetRemainingTopics request
	GetRemainingTopics(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetAnswersLog request
	GetAnswersLog(ctx context.Context, userId UUID, params *GetAnswersLogParams, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) SubmitAnswerWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSubmitAnswerRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SubmitAnswer(ctx context.Context, body SubmitAnswerJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSubmitAnswerRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) HealthCheck(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthCheckRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetNextQuestion(ctx context.Context, params *GetNextQuestionParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetNextQuestionRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetRemainingTopics(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetRemainingTopicsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetAnswersLog(ctx context.Context, userId UUID, params *GetAnswersLogParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetAnswersLogRequest(c.Server, userId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewSubmitAnswerRequest calls the generic SubmitAnswer builder with application/json body
func NewSubmitAnswerRequest(server string, body SubmitAnswerJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSubmitAnswerRequestWithBody(server, "application/json", bodyReader)
}

// NewSubmitAnswerRequestWithBody generates requests for SubmitAnswer with any type of body
func NewSubmitAnswerRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/answers")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewHealthCheckRequest generates requests for HealthCheck
func NewHealthCheckRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/health-check")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetNextQuestionRequest generates requests for GetNextQuestion
func NewGetNextQuestionRequest(server string, params *GetNextQuestionParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/questions/next")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "topic", runtime.ParamLocationQuery, params.Topic); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetRemainingTopicsRequest generates requests for GetRemainingTopics
func NewGetRemainingTopicsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/questions/remaining-topics")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetAnswersLogRequest generates requests for GetAnswersLog
func NewGetAnswersLogRequest(server string, userId UUID, params *GetAnswersLogParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "userId", runtime.ParamLocationPath, userId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/users/%s/answers", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Page != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page", runtime.ParamLocationQuery, *params.Page); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageSize != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "pageSize", runtime.ParamLocationQuery, *params.PageSize); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// SubmitAnswerWithBodyWithResponse request with any body
	SubmitAnswerWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SubmitAnswerResponse, error)

	SubmitAnswerWithResponse(ctx context.Context, body SubmitAnswerJSONRequestBody, reqEditors ...RequestEditorFn) (*SubmitAnswerResponse, error)

	// HealthCheckWithResponse request
	HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error)

	// GetNextQuestionWithResponse request
	GetNextQuestionWithResponse(ctx context.Context, params *GetNextQuestionParams, reqEditors ...RequestEditorFn) (*GetNextQuestionResponse, error)

	// GetRemainingTopicsWithResponse request
	GetRemainingTopicsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetRemainingTopicsResponse, error)

	// GetAnswersLogWithResponse request
	GetAnswersLogWithResponse(ctx context.Context, userId UUID, params *GetAnswersLogParams, reqEditors ...RequestEditorFn) (*GetAnswersLogResponse, error)
}

type SubmitAnswerResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *SubmitAnswerResult
}

// Status returns HTTPResponse.Status
func (r SubmitAnswerResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SubmitAnswerResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type HealthCheckResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r HealthCheckResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthCheckResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetNextQuestionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *UnansweredQuestion
}

// Status returns HTTPResponse.Status
func (r GetNextQuestionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetNextQuestionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetRemainingTopicsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]RemainingTopic
}

// Status returns HTTPResponse.Status
func (r GetRemainingTopicsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetRemainingTopicsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetAnswersLogResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]AnswersLogItem
}

// Status returns HTTPResponse.Status
func (r GetAnswersLogResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetAnswersLogResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// SubmitAnswerWithBodyWithResponse request with arbitrary body returning *SubmitAnswerResponse
func (c *ClientWithResponses) SubmitAnswerWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SubmitAnswerResponse, error) {
	rsp, err := c.SubmitAnswerWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSubmitAnswerResponse(rsp)
}

func (c *ClientWithResponses) SubmitAnswerWithResponse(ctx context.Context, body SubmitAnswerJSONRequestBody, reqEditors ...RequestEditorFn) (*SubmitAnswerResponse, error) {
	rsp, err := c.SubmitAnswer(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSubmitAnswerResponse(rsp)
}

// HealthCheckWithResponse request returning *HealthCheckResponse
func (c *ClientWithResponses) HealthCheckWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthCheckResponse, error) {
	rsp, err := c.HealthCheck(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthCheckResponse(rsp)
}

// GetNextQuestionWithResponse request returning *GetNextQuestionResponse
func (c *ClientWithResponses) GetNextQuestionWithResponse(ctx context.Context, params *GetNextQuestionParams, reqEditors ...RequestEditorFn) (*GetNextQuestionResponse, error) {
	rsp, err := c.GetNextQuestion(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetNextQuestionResponse(rsp)
}

// GetRemainingTopicsWithResponse request returning *GetRemainingTopicsResponse
func (c *ClientWithResponses) GetRemainingTopicsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetRemainingTopicsResponse, error) {
	rsp, err := c.GetRemainingTopics(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetRemainingTopicsResponse(rsp)
}

// GetAnswersLogWithResponse request returning *GetAnswersLogResponse
func (c *ClientWithResponses) GetAnswersLogWithResponse(ctx context.Context, userId UUID, params *GetAnswersLogParams, reqEditors ...RequestEditorFn) (*GetAnswersLogResponse, error) {
	rsp, err := c.GetAnswersLog(ctx, userId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetAnswersLogResponse(rsp)
}

// ParseSubmitAnswerResponse parses an HTTP response from a SubmitAnswerWithResponse call
func ParseSubmitAnswerResponse(rsp *http.Response) (*SubmitAnswerResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SubmitAnswerResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest SubmitAnswerResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseHealthCheckResponse parses an HTTP response from a HealthCheckWithResponse call
func ParseHealthCheckResponse(rsp *http.Response) (*HealthCheckResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthCheckResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetNextQuestionResponse parses an HTTP response from a GetNextQuestionWithResponse call
func ParseGetNextQuestionResponse(rsp *http.Response) (*GetNextQuestionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetNextQuestionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest UnansweredQuestion
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetRemainingTopicsResponse parses an HTTP response from a GetRemainingTopicsWithResponse call
func ParseGetRemainingTopicsResponse(rsp *http.Response) (*GetRemainingTopicsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetRemainingTopicsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []RemainingTopic
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetAnswersLogResponse parses an HTTP response from a GetAnswersLogWithResponse call
func ParseGetAnswersLogResponse(rsp *http.Response) (*GetAnswersLogResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetAnswersLogResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []AnswersLogItem
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /answers)
	SubmitAnswer(w http.ResponseWriter, r *http.Request)

	// (GET /health-check)
	HealthCheck(w http.ResponseWriter, r *http.Request)

	// (GET /questions/next)
	GetNextQuestion(w http.ResponseWriter, r *http.Request, params GetNextQuestionParams)

	// (GET /questions/remaining-topics)
	GetRemainingTopics(w http.ResponseWriter, r *http.Request)

	// (GET /users/{userId}/answers)
	GetAnswersLog(w http.ResponseWriter, r *http.Request, userId UUID, params GetAnswersLogParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// SubmitAnswer operation middleware
func (siw *ServerInterfaceWrapper) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SubmitAnswer(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// HealthCheck operation middleware
func (siw *ServerInterfaceWrapper) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.HealthCheck(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetNextQuestion operation middleware
func (siw *ServerInterfaceWrapper) GetNextQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetNextQuestionParams

	// ------------- Required query parameter "topic" -------------

	if paramValue := r.URL.Query().Get("topic"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "topic"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "topic", r.URL.Query(), &params.Topic)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "topic", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetNextQuestion(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetRemainingTopics operation middleware
func (siw *ServerInterfaceWrapper) GetRemainingTopics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRemainingTopics(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetAnswersLog operation middleware
func (siw *ServerInterfaceWrapper) GetAnswersLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "userId" -------------
	var userId UUID

	err = runtime.BindStyledParameterWithOptions("simple", "userId", r.PathValue("userId"), &userId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "userId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAnswersLogParams

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pageSize", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAnswersLog(w, r, userId, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/answers", wrapper.SubmitAnswer)
	m.HandleFunc("GET "+options.BaseURL+"/health-check", wrapper.HealthCheck)
	m.HandleFunc("GET "+options.BaseURL+"/questions/next", wrapper.GetNextQuestion)
	m.HandleFunc("GET "+options.BaseURL+"/questions/remaining-topics", wrapper.GetRemainingTopics)
	m.HandleFunc("GET "+options.BaseURL+"/users/{userId}/answers", wrapper.GetAnswersLog)

	return m
}
