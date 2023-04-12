
# Usage: tools to cleanup the blob in the docker registry

## Prerequest

1. Please mount the volumn of base directory to the container, and run the command in the container. 
2. If the database is in the k8s cluster, forward the db container's 5432 to localhost
3. Change the harbor into read only mode, backup the PV and database.

4. Run the goloang program to cleanup the blob in the registry
```

go run ./cmd/cleanup/main.go -base_dir=/var/lib/registry/docker/registry/v2/blobs/sha256 -dry_run=true

```

5. cleanup redis cache in the registry

```
redis-cli>
# make sure most of the keys are start with blobs::sha25
>keys *
>flushdb
>dbsize
```
5. Try to push and pull images in the registry
