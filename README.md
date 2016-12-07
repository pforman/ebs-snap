# ebs-snap

Takes EBS snapshots of volumes.  Few frills.

---

`ebs-snap` is a single binary to take snapshots of a volume attached to an instance in AWS.

It specifies an `Expires:` tag for eventual cleanup.  The default expiration time is *7* days.

The `ebs-snap` binary is intended to be deployed on a server to create snapshots of local filesystems.  For remote operation, see `remote-snap` below.  `cron` is a good method to schedule these snapshots.  If necessary, pre-snapshot and post-snapshot scripts can be specified (for example to lock database tables).


### Usage

```
ebs-snap /
```

The mount point (in this case `/`) is matched to a device using /proc/mounts, then the corresponding volume-id is retrieved from the API.  A snapshot is created of that volume-id, then all tags from the source volume are copied, plus an `Expires` tag, formatted as an ISO8601 date string.

### Options

```
  -device string
    	device to snapshot (for unmounted volumes only)
  -expires int
    	sets the expiration time in days (default 7)
  -post string
    	command to run after snapshot
  -pre string
    	command to run before snapshot
  -v	verbose mode, provides more info
  -version
    	print version string, then exit
```

### Credentials

`ebs-snap` uses the Go AWS SDK, so it understands credentials in environment variables (including STS tokens).

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



# remote-snap

If desired, the `remote-snap` binary can be used to initiate a snapshot of any instance, however the device and mountpoint must be provided.  The mountpoint in this case is used only for tagging, but is still a required parameter.

### Usage
```
remote-snap 
```

### Options

```
  -device device
    	device to snapshot
  -expires int
    	sets the expiration time in days (default 7)
  -instance instance-id
    	instance-id of the instance to snapshot
  -region region
    	region of instance
  -v	verbose mode, provides more info
  -version
    	print version string, then exit
```

# clean-snaps

The `clean-snaps` command can be used to remove expired snapshots (those with an Expired tag in the past).  It can be run from anywhere with permissions, although outside EC2 the region must be provided as a parameter or by using the AWS_REGION environment variable.

The `-clean` flag serves as a safety, and must be provided to enable snapshot deletion.

### Usage
```
clean-snaps -clean
```

### Options
```
  -clean
    	clean mode, will remove expired snapshots
  -dryrun
    	dry run mode, shows what would be cleaned
  -region region
    	region of instance
  -v	verbose mode, provides more info
  -version
    	print version string, then exit
```