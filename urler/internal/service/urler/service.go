package urler

import (
	"context"
	"net/http"

	"github.com/dbeleon/urler/urler/internal/domain"
	api "github.com/dbeleon/urler/urler/pkg/urler/v1"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	api.UnimplementedUrlerServiceServer

	businessLogic *domain.Model
}

func New(model *domain.Model) *Server {
	return &Server{businessLogic: model}
}

func ResponseHeaderMatcher(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}
