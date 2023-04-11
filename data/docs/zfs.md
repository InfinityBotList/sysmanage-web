# ZFS setup

The optimal ZFS structure for MongoDB in IBL has been found to currently be:

- ZFS pool named silverpelt created on a 100GB DO volume
- On silverpelt, a XFS filesystem attached to a zvol named nightheart (silverpelt/nightheart) of 90-95GB
- Ubuntu 22.04 LTS with zfsutils-linux used (not FUSE)
- Sunbeam as ZPOOL. ARC of 5368709120 following https://bun.uptrace.dev/postgres/tuning-zfs-aws-ebs.html#zfs-config. Set data_directory in postgres to /silverpelt/sunbeam. Use ZFS native encryption here if possible
