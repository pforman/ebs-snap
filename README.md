# ebs-snap

Takes EBS snapshots of volumes.  Few frills.

---

`ebs-snap` is a single binary to take snapshots of a volume attached to an instance in AWS.

It specifies an `Expires:` tag for eventual cleanup.  The default expiration time is *7* days.

The options are set to prefer on-instance operation, however snapshots an be initated from anywhere if you provide the correct device and mount mapping for the instance.

## Usage

```
ebs-snap /
```

The mount point (in this case `/`) is matched to a device using /proc/mounts, then the corresponding volume-id is retrieved from the API.  RA snapshot is created of that volume-id, then all tags from the source volume are copied, plus an `Expires` tag, formatted as an ISO8601 date string.

## Options

```
Usage of ./ebs-snap:
  -device string
    	device to snapshot (for remote snaps only, be careful!)
  -expires int
    	sets the expiration time in days (default 7)
  -instance string
    	instance-id (for remote snaps only)
  -region string
    	region of instance (for remote snaps only)
  -v	verbose mode, provides more info
  -version
    	print version string, then exit
```

## Credentials

`ebs-snap` use the Go AWS SDK, so it understands credentials in environment variables (including STS tokens).

However, the preferred method is to make use of an instance profile with appropriate permissions to create snapshots.

The permissions required are (at minimum):
```
{
  "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "SnapshotActions",
        "Effect": "Allow",
        "Action": [
          "ec2:CreateSnapshot",
          "ec2:CreateTags",
          "ec2:DescribeInstances",
          "ec2:DescribeSnapshots",
          "ec2:DescribeTags",
          "ec2:DescribeVolumes"
        ],
        "Resource": [
          "*"
        ]
      }
    ]
}
```

`DeleteSnapshots` and `DeleteTags` are deliberately excluded, as adding these to a host profile can compromise the security of all backups.  Don't do that.


## Build

`go build`

Testing and a more complex build process are in the works.
