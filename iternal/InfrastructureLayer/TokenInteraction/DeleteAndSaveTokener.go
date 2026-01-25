package TokenInteraction

import (
	"Kaban/iternal/Dto"
	"time"
)

func (*ControlTokens) DeleteAndSaveToken(olRtKey string, newRtKey string) {
	delete(Dto.AllowList, olRtKey)
	Dto.DenyList[olRtKey] = time.Now()
	Dto.AllowList[newRtKey] = time.Now()

}
