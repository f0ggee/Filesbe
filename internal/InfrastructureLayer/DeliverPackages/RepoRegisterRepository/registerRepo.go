package RepoRegisterRepository

import "net/http"

type RegisterController struct{}

var ErrorController = &RegisterController{}

type Register interface {
	ErrorAnswer(RegisterErrorIncomingData)
	GoodAnswer(RegisterGoodIncomingData)
}
type RegisterErrorIncomingData struct {
	W         http.ResponseWriter
	Error     error
	Operation error
}
type RegisterGoodIncomingData struct {
	W         http.ResponseWriter
	Operation string
	Redirect  string
}

func GetNewRegisterController() *RegisterController {
	return &RegisterController{}
}
