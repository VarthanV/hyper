package hyper

type HandlerFunc func(w ResponseWriter, request *Request)
