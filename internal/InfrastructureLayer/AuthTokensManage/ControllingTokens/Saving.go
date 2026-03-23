package ControllingTokens

import (
	"Kaban/internal/Dto"
	"time"
)

func (c ManageTokens) SaveToken(s string) {
	Dto.AllowList[s] = time.Now()

}
