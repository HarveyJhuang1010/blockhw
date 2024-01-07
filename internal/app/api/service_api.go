package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type apiService struct {
	in      digIn
	handler http.Handler
	srv     *http.Server
}

func newApiService(in digIn) bo.Service {
	svc := &apiService{
		in: in,
	}
	svc.createService()
	return svc
}

func (s *apiService) Run(ctx context.Context, stop context.CancelFunc) {
	logger := appcontext.GetLogger(ctx)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", zap.Any("panic", r))
		}
		stop()
	}()

	cfg := config.GetConfig()
	srv := &http.Server{
		ReadHeaderTimeout: time.Second * 10,
		Addr:              fmt.Sprintf("%s:%d", cfg.Rest.ListenAddress, cfg.Rest.ListenPort),
		Handler:           s.handler,
	}
	s.srv = srv

	logger.Info("Server is running", zap.String("address", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server Error", zap.Error(err))
	}
}

func (s *apiService) Shutdown(ctx context.Context, wg *sync.WaitGroup) {
	logger := appcontext.GetLogger(ctx)

	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second /*replace with config*/)
	defer func() {
		cancel()
		wg.Done()
	}()
	logger.Info("start shutdown api server")
	if err := s.srv.Shutdown(ctxT); err != nil {
		logger.Error("shutdown error", zap.Error(err))
	}
	logger.Info("end shutdown api server")
}

func (s *apiService) createService() {
	// Create our HTTP Router
	e := gin.New()
	cfg := config.GetConfig()
	ctx := appcontext.GetContext()
	logger := appcontext.GetLogger(ctx)

	// health check
	e.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "alive")
	})

	// Register Handlers, Middleware, and API Modules
	s.registerMiddleware(cfg, e, logger)
	s.registerPublicRoutes(e)

	s.handler = e
}

func (s *apiService) registerMiddleware(cfg *config.Config, e *gin.Engine, l *zap.Logger) {
	e.Use(
		gin.Recovery(),
	)

	e.Use(ratelimit.RateLimiter(
		// storage
		ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
			Rate:  time.Duration(cfg.Rest.RateLimitIntervalSeconds) * time.Second,
			Limit: cfg.Rest.RateLimitRequestPerSecond,
		}),
		// options
		&ratelimit.Options{
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.String(http.StatusTooManyRequests, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
			},
			KeyFunc: func(c *gin.Context) string {
				key := c.GetHeader("Cf-Connecting-Ip")
				if len(key) > 0 {
					return key + c.Request.URL.Path
				}
				return c.ClientIP()
			},
		},
	))

	e.Use(cors.New(cors.Config{
		AllowOrigins:  cfg.Rest.AllowOrigins,
		AllowMethods:  cfg.Rest.AllowMethods,
		AllowHeaders:  cfg.Rest.AllowHeaders,
		ExposeHeaders: cfg.Rest.ExposeHeaders,
	}))
}

func (s *apiService) registerPublicRoutes(e *gin.Engine) {
	// GET /blocks?limit=n
	// GET /blocks/:id
	// GET /transactions/:txHash
}
