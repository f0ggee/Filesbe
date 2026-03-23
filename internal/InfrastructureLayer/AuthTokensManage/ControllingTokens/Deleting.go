package ControllingTokens

import (
	"Kaban/internal/Dto"
	"time"
)

func (c ManageTokens) DeleteRefreshToken(Rf string) {
	delete(Dto.AllowList, Rf)
	Dto.DenyList[Rf] = time.Now()
}
