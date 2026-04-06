package DeleterS3

import (
	"Kaban/internal/InfrastructureLayer/s3Interation"
	"context"
	"errors"
	"time"
)

type DeleterS3 struct {
	S3Info s3Interation.Variables
}

func (d *DeleterS3) DeleterS3Test(s string, Cont context.Context) error {

	time.Sleep(2 * time.Second)
	sa, de := Cont.Value("IsFall").(bool)
	if sa != false {
		if de {
			return errors.New("error by s3")
		}
	}
	return nil

}
