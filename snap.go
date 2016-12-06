package snap

import (
	// "flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CreateSnapshot(s *session.Session, instance string, device string, mount string, volumeId string, expires int) error {

	// Set up the expiration time
	startingTime := time.Now().UTC()
	expireTag := fmt.Sprintf("%v", startingTime.AddDate(0, 0, expires).Round(time.Second))
	if Verbose() {
		println("Expiration time set to", expireTag)
	}

	// old autosnap uses hostname instead of instance-id
	// maybe we should find that...
	snapDesc := fmt.Sprintf("ebs-snap %s:%s:%s", instance, device, mount)
	snapId, err := EC2Snapshot(s, volumeId, snapDesc)
	if err != nil {
		fmt.Printf("error creating snapshot for volume %s: %s\n", volumeId, err.Error())
		return err
	}
	if Verbose() {
		println("Created snapshot with id: ", snapId)
	} else {
		println(snapId)
	}
	err = TagSnapshot(s, snapId, volumeId, expireTag)
	if err != nil {
		println("error in tagging:", err)
		// delete here on error, if we can...
		return err
	}
	return nil
}

func EC2Snapshot(s *session.Session, volumeId string, desc string) (string, error) {
	svc := ec2.New(s)

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

func TagSnapshot(s *session.Session, snapId string, volumeId string, expires string) error {
	svc := ec2.New(s)

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
