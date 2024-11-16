package hyper

const (
	delimNewLine            = '\n'
	pathParamMatchingRegexp = `^/[^/]+/:.+$`
)

type routeStruct struct {
	UserGivenPath string
	Handler       HandlerFunc
}
