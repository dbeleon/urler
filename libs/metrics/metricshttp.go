package metrics

import (
	"context"
	"net/http"
	"sync"

	"github.com/dbeleon/urler/libs/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// MetricsHttpServer сервер для получения метрик приложения прометеем
type MetricsHttpServer struct {
	server *http.Server
	wg     *sync.WaitGroup
}

// New создает новый сервер
func New(addr string) *MetricsHttpServer {
	srv := &http.Server{Addr: addr}
	http.Handle("/metrics", promhttp.Handler())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		log.Info("metrics listening", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("failed to listen and serve http server for prometheus", zap.Error(err), zap.String("func", "ListenAndServe"), zap.String("address", addr))
		}
	}()

	return &MetricsHttpServer{
		server: srv,
		wg:     wg,
	}
}

// Close останавливает сервер
func (m *MetricsHttpServer) Close(ctx context.Context) {
	if err := m.server.Shutdown(ctx); err != nil {
		log.Fatal("http server for prometheus",
			zap.Error(err),
			zap.String("func", "Shutdown"),
			zap.String("address", m.server.Addr),
		)
	}

	m.wg.Wait()
}
