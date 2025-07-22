package processor

import (
	"net/http"
)

func (p *Processor) GetOamIndex() *HandlerResponse {
	return &HandlerResponse{http.StatusOK, nil, nil}
}
