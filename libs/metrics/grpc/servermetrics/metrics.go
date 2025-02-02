package servermetrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	RequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "urler",
		Subsystem: "gRPC",
		Name:      "requests_total",
	},
		[]string{"handle"},
	)
	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "urler",
		Subsystem: "gRPC",
		Name:      "responses_total",
	},
		[]string{"handle", "status"},
	)
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "urler",
		Subsystem: "gRPC",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"handle", "status"},
	)
)

// Interceptor подсчитывает метрики запроса к gRPC серверу
func Interceptor() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		RequestsCounter.WithLabelValues(info.FullMethod).Inc()

		timeStart := time.Now()

		res, err := handler(ctx, req)

		elapsed := time.Since(timeStart)

		st, _ := status.FromError(err)

		HistogramResponseTime.WithLabelValues(info.FullMethod, st.Code().String()).Observe(elapsed.Seconds())
		ResponseCounter.WithLabelValues(info.FullMethod, st.Code().String()).Inc()

		return res, err
	}
}
