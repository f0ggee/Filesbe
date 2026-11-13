package testFunctions

import "context"

type GoodDownload struct {
}

func (g GoodDownload) DownloadApplication(ctx context.Context, s string) error {
	return nil
}
