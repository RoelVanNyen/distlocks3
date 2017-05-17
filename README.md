# Distribute locking using s3 backend

A golang implementation of a locking mechanism for amazon S3 that leverages the versioning feature to create
locks. A new version of a file is uploaded first, then checks if the oldest
non-deleted version is the same as the one uploaded. Once it 'acquires' the lock
this way, runs the critical section, then removes its version.
