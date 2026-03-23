package Dto

import "time"

var MapForFile = make(map[string]struct {
	AesKey          string
	TimeSet         time.Time
	IsUsed          bool
	IsStartDownload bool
})
