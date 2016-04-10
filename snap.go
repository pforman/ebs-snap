package main

import (
	// "flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func snapVolume(session *session.Session, volumeId string, desc string) (string, error) {
	svc := ec2.New(session)

	params := &ec2.CreateSnapshotInput{
		VolumeId:    aws.String(volumeId),
		Description: aws.String(desc),
	}

	resp, err := svc.CreateSnapshot(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return "", err
	}

	// Pretty-print the response data.
	return *resp.SnapshotId, nil

}

func tagSnapshot(session *session.Session, snapId string, expires string) error {
	svc := ec2.New(session)

	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(snapId),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Expires"),
				Value: aws.String(expires),
			},
			// More values...
		},
	}
	_, err := svc.CreateTags(params)

	return err
}
