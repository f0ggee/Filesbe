package Dto

type FileInfo struct {
	InfoAboutFile   []byte `redis:"InfoAboutFile"`
	IsStartDownload bool   `redis:"IsStartDownload"`
}
