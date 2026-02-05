package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/TokenInteraction"
)

type Generates struct {
	Tokens DomainLevel.ManageTokens
}

func NewGenerateTokenEr(Re DomainLevel.ManageTokens) *Generates {
	return &Generates{Tokens: Re}
}

func SetSittingsTokenInteraction() *Generates {

	as := &TokenInteraction.ControlTokens{A: false}
	app := NewGenerateTokenEr(as)
	return app
}
