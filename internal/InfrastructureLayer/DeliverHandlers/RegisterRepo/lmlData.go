package RegisterRepo

import (
	"github.com/go-playground/validator/v10"
)

type RegisterController struct {
	v *validator.Validate
}
