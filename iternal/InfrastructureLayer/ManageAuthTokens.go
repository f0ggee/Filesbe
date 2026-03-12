package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/TokenInteraction"
)

type Generates struct {
	Tokens DomainLevel.ManageAuthTokens
}

func NewGenerateTokenEr(Re DomainLevel.ManageAuthTokens) *Generates {
	return &Generates{Tokens: Re}
}

func SetSittingsTokenInteraction() *Generates {

	as := &TokenInteraction.ControlTokens{A: false}
	app := NewGenerateTokenEr(as)
	return app
}
