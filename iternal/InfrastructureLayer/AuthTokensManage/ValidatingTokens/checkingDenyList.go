package ValidatingTokens

import "Kaban/iternal/Dto"

func (c Checking) CheckingDenyList(s string) bool {

	DenyLists := Dto.DenyList
	if _, ok := DenyLists[s]; !ok {

		return false
	}
	return true
}
