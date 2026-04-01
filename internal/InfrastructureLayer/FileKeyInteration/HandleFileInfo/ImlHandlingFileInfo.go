package HandleFileInfo

import Dto2 "Kaban/internal/Dto"

type ProcessingFileInfo struct{}

func (p ProcessingFileInfo) GetRealNameFile(TemporallyName string) string {
	if Name, ok := Dto2.NamesToConvert[TemporallyName]; ok {
		NameFile := Name
		delete(Dto2.NamesToConvert, TemporallyName)
		return NameFile
	}
	return ""
}

func (p ProcessingFileInfo) SayHi() string {

	return "Hi"
}
