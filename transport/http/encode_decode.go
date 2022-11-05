package http

import (
	"context"
	"net/http"
)

type DecodeRequestFunc[Request any] func(context.Context, *http.Request) (request *Request, err error)

type EncodeResponseFunc[Response any] func(context.Context, http.ResponseWriter, *Response) error
