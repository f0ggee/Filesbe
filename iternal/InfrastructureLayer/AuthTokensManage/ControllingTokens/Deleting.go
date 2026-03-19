package ControllingTokens

import (
	"Kaban/iternal/Dto"
	"time"
)

func (c ControllingTokens) DeleteRefreshToken(Rf string) {
	delete(Dto.AllowList, Rf)
	Dto.DenyList[Rf] = time.Now()
}
