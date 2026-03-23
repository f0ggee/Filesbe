package Dto

type FileInfoLabels struct {
	InfoAboutFile   []byte `redis:"InfoAboutFile"`
	IsStartDownload bool   `redis:"IsStartDownload"`
}
