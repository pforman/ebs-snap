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

	return *resp.SnapshotId, nil

}

func tagSnapshot(session *session.Session, snapId string, volumeId string, expires string) error {
	svc := ec2.New(session)

	dparams := &ec2.DescribeTagsInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(volumeId),
				},
			},
		},
	}

	res, err := svc.DescribeTags(dparams)
	if err != nil {
		fmt.Printf("error reading tags from volume %s: %v", volumeId, err)
		// forge on, we're just tagging
	}

	tags := []*ec2.Tag{
		{
			Key:   aws.String("Expires"),
			Value: aws.String(expires),
		},
	}

	for _, v := range res.Tags {
		tags = append(tags, &ec2.Tag{Key: aws.String(*v.Key), Value: aws.String(*v.Value)})
	}

	cparams := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(snapId),
		},
		Tags: tags,
	}

	_, err = svc.CreateTags(cparams)

	return err
}
