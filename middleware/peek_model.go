package middleware

import (
	"net/http"

	"github.com/alpe/netutils/common"
)

const XModelHeader = "X-Model"

var _ http.Handler = (*PeekModelMiddleware)(nil)

// PeekModelMiddleware implements an optimistic approach to read the model name from the request payload when it is the first
// attribute of a json object payload. This may not always be the case so the downstream processes need to take care
// of the payload handling themselves.
//
// If the X-Model header is already set, it will not modify the request.
type PeekModelMiddleware struct {
	nextHandler http.Handler
	bufferSize  int
}

// NewPeekModelMiddleware is the constructor for the PeekModelMiddleware type.
// It takes a http.Handler as a parameter, which is the next handler in the middleware chain.
// The bufferSize defines the number of bytes read to fetch the model attribute.
func NewPeekModelMiddleware(nextHandler http.Handler, bufferSize int) *PeekModelMiddleware {
	return &PeekModelMiddleware{nextHandler: nextHandler, bufferSize: bufferSize}
}

// ServeHTTP handle the request.
func (p *PeekModelMiddleware) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Header.Get(XModelHeader) == "" {
		model, stream, err := common.PeekModel(req.Body, p.bufferSize)
		if err != nil {
			http.Error(resp, "invalid content", http.StatusBadRequest)
			return
		}
		if model != "" {
			req.Header.Set(XModelHeader, model)
		}
		req.Body = stream
	}
	p.nextHandler.ServeHTTP(resp, req)
}
