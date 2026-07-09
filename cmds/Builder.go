package cmds

import (
	"Kaban/internal/Deliver"
	"Kaban/internal/InfrastructureLayer/DeliverPackages/RepoDownloadNoEncrypt"
)

//application builders

func GetControllerDownloadBuilder(DeliverPackagesCollector *deliverPackagesCollector) *Deliver.NewDownloadWithNotEncrypt {
	return Deliver.GetNewDownloadWithNotEncrypt(Deliver.AnswerDownloadNoEncrypt{Answ: DeliverPackagesCollector.DownloadCollector.Answ}, Deliver.UrlBuilderDownloadNoEncrypt{UrlWork: RepoDownloadNoEncrypt.NewRepoDownloadNoEncrypt(DeliverPackagesCollector.UrlBuilderCollector.Answ)}, Deliver.NetworkDownloadNoEncrypt{}, Deliver.NewDownloadWithNotEncryptApplication{})
}
