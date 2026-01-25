package TokenInteraction

import "Kaban/iternal/Dto"

type ControlTokens struct {
	A any
}

func (s *ControlTokens) TokenDenyMapChecker(RT string) bool {

	DenyLists := Dto.DenyList
	if _, ok := DenyLists[RT]; !ok {

		return false
	}
	return true
}
