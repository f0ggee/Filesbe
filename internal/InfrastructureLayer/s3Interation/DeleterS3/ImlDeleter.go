package DeleterS3

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DeleterS3 struct {
	Conf *s3.Client
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
