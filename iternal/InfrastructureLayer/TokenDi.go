package InfrastructureLayer

import "Kaban/iternal/DomainLevel"

type Generates struct {
	Tokens DomainLevel.ManageTokens
}

func NewGenerateTokenEr(Re DomainLevel.ManageTokens) *Generates {
	return &Generates{Tokens: Re}
}
