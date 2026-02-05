package Handlers

import (
	"Kaban/iternal/InfrastructureLayer"
)

func DeleteFile(NameFile string, IsEncrypt bool) {
	S3Interation := *InfrastructureLayer.NewConnectToS3()

	if IsEncrypt {
		redisConnect := *InfrastructureLayer.NewSetRedisConnect()

		err := redisConnect.Ras.DeleteFileInfo(NameFile)
		if err != nil {
			return
		}
	}

	err := S3Interation.Manage.DeleteFileFromS3(NameFile, Bucket)
	if err != nil {

		return
	}

	return

}
