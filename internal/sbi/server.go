package sbi

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/Alonza0314/nf-example/pkg/app"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/gin-gonic/gin"
)

type nfApp interface {
	app.App
	Processor() *processor.Processor
}

type Server struct {
	nfApp

	httpServer *http.Server
	router     *gin.Engine
}

func NewServer(nf nfApp, tlsKeyLogPath string) *Server {
	s := &Server{
		nfApp: nf,
	}

	s.router = newRouter(s)

	server, err := bindRouter(nf, s.router, tlsKeyLogPath)
	s.httpServer = server
	if err != nil {
		logger.SBILog.Errorf("bind Router Error: %+v", err)
		panic("Server initialization failed")
	}

	return s
}

func (s *Server) Run(wg *sync.WaitGroup) {
	logger.SBILog.Info("Starting server...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.SBILog.Infof("Start SBI server (listen on %s)", s.httpServer.Addr)

		err := s.serve()
		if err != http.ErrServerClosed {
			logger.SBILog.Panicf("HTTP server setup failed: %+v", err)
		}
		logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
	}()
}

func (s *Server) unsecureServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) secureServe() error {
	sbiConfig := s.Config().Configuration.Sbi

	pemPath := sbiConfig.Tls.Pem
	if pemPath == "" {
		pemPath = factory.NfDefaultCertPemPath
	}

	keyPath := sbiConfig.Tls.Key
	if keyPath == "" {
		keyPath = factory.NfDefaultPrivateKeyPath
	}

	return s.httpServer.ListenAndServeTLS(pemPath, keyPath)
}

func (s *Server) serve() error {
	sbiConfig := s.Config().Configuration.Sbi

	switch sbiConfig.Scheme {
	case "http":
		return s.unsecureServe()
	case "https":
		return s.secureServe()
	default:
		return fmt.Errorf("invalid SBI scheme: %s", sbiConfig.Scheme)
	}
}

func (s *Server) Shutdown() {
	s.shutdownHttpServer()
}

func (s *Server) shutdownHttpServer() {
	logger.SBILog.Infoln("Shutdown Http Server...")
	const shutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.SBILog.Errorf("HTTP server shutdown failed: %+v", err)
	}
}

func (s *Server) getStudentRoute() []Route {
	return []Route{
		{
			Name:    "GetStudentData",
			Method:  "GET",
			Pattern: "/data",
			APIFunc: GetStudentData,
		},
		{
			Name:    "RegisterStudent",
			Method:  "POST",
			Pattern: "/register",
			APIFunc: RegisterStudent,
		},
	}
}
