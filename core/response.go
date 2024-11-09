package core

type Response struct {
	statusCode int
}

func (r *Response) WriteStatus(code int) {
	r.statusCode = code
}
