package RegisterRepo

type RegisterController struct{}

var ErrorController = &RegisterController{}

func NewRegisterController() *RegisterController {
	return &RegisterController{}
}
