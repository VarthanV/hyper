package core

type Response struct {
	statusCode int
}

func (r *Response) WriteStatusCode(code int) {
	r.statusCode = code
}
