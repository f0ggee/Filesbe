package FileKeyInteration

import Dto2 "Kaban/iternal/Dto"

func (*FileInfoController) GetRealNameFile(TemporallyName string) string {
	if Name, ok := Dto2.NamesToConvert[TemporallyName]; ok {
		NameFile := Name
		delete(Dto2.NamesToConvert, TemporallyName)
		return NameFile
	}
	return ""
}
