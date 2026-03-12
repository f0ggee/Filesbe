package InfrastructureLayer

import (
	"Kaban/iternal/DomainLevel"
	"Kaban/iternal/InfrastructureLayer/s3Interation"
)

type S3Connection struct {
	Manage DomainLevel.S3Handle
}

func NewS3Connection(manage DomainLevel.S3Handle) *S3Connection {
	return &S3Connection{Manage: manage}
}

func NewConnectToS3() *S3Connection {

	apps := &s3Interation.ConntrolerForS3{}

	s := NewS3Connection(apps)

	return s
}
