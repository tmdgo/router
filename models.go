package router

type Controller interface {
	GetRoutes() []Route
}

type Route struct {
	Path            string
	Method          string
	HandleFunc      interface{}
	UseVars         bool
	UseOptionalVars bool
	UseRequestModel bool
	RequestModel    interface{}
}

type Vars struct {
	Value map[string]string
}

type OptionalVars struct {
	Value map[string][]string
}

type Result struct {
	StatusCode int
	Model      interface{}
}

type Error struct {
	StatusCode int
	Message    string
	Err        error
}

type jsonError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
