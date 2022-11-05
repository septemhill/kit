package endpoint

import "context"

type Endpoint[Request, Response any] func(context.Context, *Request) (*Response, error)
