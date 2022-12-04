package rpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/exfly/vulcan/internel/config"
	userrepo "github.com/exfly/vulcan/internel/user/repo"
	vulcanv1 "github.com/exfly/vulcan/pb"
	"github.com/exfly/vulcan/pkg/token"

	"github.com/sirupsen/logrus"
)

func New(config *config.Config, userQ userrepo.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	return &Server{
		config:     config,
		store:      userQ,
		tokenMaker: tokenMaker,
	}, nil
}

type Server struct {
	vulcanv1.UnimplementedSimpleBankServer

	config     *config.Config
	store      userrepo.Store
	tokenMaker token.Maker
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := logrus.StandardLogger()
		logEntry := logger.WithField("code", rec.StatusCode)
		if rec.StatusCode != http.StatusOK {
			logEntry = logEntry.WithField("body", string(rec.Body))
		}

		logEntry.
			WithField("protocol", "http").
			WithField("method", req.Method).
			WithField("path", req.RequestURI).
			WithField("status_code", rec.StatusCode).
			WithField("status_text", http.StatusText(rec.StatusCode)).
			WithField("duration", duration).
			Info("received a HTTP request")
	})
}
