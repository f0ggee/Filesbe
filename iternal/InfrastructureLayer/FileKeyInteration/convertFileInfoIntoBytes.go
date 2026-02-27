package FileKeyInteration

import (
	"Kaban/iternal/Dto"
	"encoding/hex"
	"encoding/json"
	"log/slog"
)

type FileInfoController struct{}

func (d *FileInfoController) SayHi() string {
	return "Hi"
}

func (d *FileInfoController) ConvertToBytesFileInfo(NameFile string, AesKey []byte) ([]byte, error) {

	FileAboutStruct := Dto.FileDescription{
		FileName: NameFile,
		AesKey:   hex.EncodeToString(AesKey),
	}

	FileInfoJson, err := json.Marshal(FileAboutStruct)
	if err != nil {
		slog.Error("Error while marshalling FileAboutSturct")
		return nil, err
	}

	return FileInfoJson, nil

}
