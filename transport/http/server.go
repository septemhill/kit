package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/septemhill/kit/endpoint"
	"github.com/septemhill/kit/transport"
)

type Server[Request, Response any] struct {
	e          endpoint.Endpoint[Request, Response]
	dec        DecodeRequestFunc[Request]
	enc        EncodeResponseFunc[Response]
	errHandler transport.ErrorHandler
}

func DecodeRequest[Request any](_ context.Context, r *http.Request) (*Request, error) {
	var req Request
	v := mux.Vars(r)

	buff := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buff).Encode(v); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(buff).Decode(&req); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return &req, nil
}

func EncodeJSONResponse[Response any](_ context.Context, w http.ResponseWriter, rsp *Response) error {
	b, err := json.Marshal(rsp)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func (s Server[Request, Response]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := s.dec(ctx, r)
	if err != nil {
		s.errHandler.Handle(ctx, err)
		return
	}

	rsp, err := s.e(ctx, req)
	if err != nil {
		s.errHandler.Handle(ctx, err)
		return
	}

	if err := s.enc(ctx, w, rsp); err != nil {
		s.errHandler.Handle(ctx, err)
		return
	}
}

func NewServer[Request, Response any](
	e endpoint.Endpoint[Request, Response],
	dec DecodeRequestFunc[Request],
	enc EncodeResponseFunc[Response],
) *Server[Request, Response] {
	return &Server[Request, Response]{
		e:   e,
		dec: dec,
		enc: enc,
	}
}

func NewEndpoint[Request, Response any](e endpoint.Endpoint[Request, Response]) *Server[Request, Response] {
	return NewServer(e, DecodeRequest[Request], EncodeJSONResponse[Response])
}
