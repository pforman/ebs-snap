package snap

import (
	// "flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func CleanExpiredSnapshot(s *session.Session, dryrun bool) error {

	now := time.Now().UTC()
	svc := ec2.New(s)

	params := &ec2.DescribeSnapshotsInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag-key"),
				Values: []*string{
					aws.String("Expires"),
				},
			},
			// More values...
		},
		MaxResults: aws.Int64(1000),
		OwnerIds: []*string{
			aws.String("self"), // Required
			// More values...
		},
		RestorableByUserIds: []*string{
			aws.String("self"), // Required
			// More values...
		},
	}


	err := svc.DescribeSnapshotsPages(params,
		func(page *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
			for _,e := range page.Snapshots {
				for _, t := range e.Tags {
					if *t.Key == "Expires" {
						expiry, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", *t.Value)
						if err != nil {
							fmt.Println(err)
							// If we can't parse the time, move on.
							continue
						}
						if now.After(expiry) {
							fmt.Printf("EXPIRE snapshot %s - %s expires %s\n", *e.SnapshotId, *e.Description, *t.Value)
							if !dryrun {
								err = deleteSnapshot(s, *e.SnapshotId)
							}
						}
					}
				}
				// fmt.Printf("snapshot %s expires %s\n", *e.SnapshotId, e.Tags)

			}
			return !lastPage
		})

	return err

}

func deleteSnapshot(s *session.Session, snap string) error {

	svc := ec2.New(s)

	params := &ec2.DeleteSnapshotInput{
		DryRun: aws.Bool(false),
		SnapshotId: aws.String(snap),
	}

	_, err := svc.DeleteSnapshot(params)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}
