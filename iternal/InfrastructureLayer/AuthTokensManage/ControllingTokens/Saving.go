package ControllingTokens

import (
	"Kaban/iternal/Dto"
	"time"
)

func (c ControllingTokens) SaveToken(s string) {
	Dto.AllowList[s] = time.Now()

}
