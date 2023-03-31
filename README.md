
# Usage: tools to cleanup the blob in the docker registry

## Prerequest

1. Please mount the volumn of base directory to the container, and run the command in the container. 
2. If the database is in the k8s cluster, forward the db container's 5432 to localhost

```
docker build -t firstfloor/cleanup_blob:1.0 .
docker run -v <path of sha256>:/basedir firstfloor/cleanup_blob:1.0 cleanup_blob --help
Usage of cleanup_blob:
  -base_dir string
        Base directory to scan, for example: /var/lib/registry/docker/registry/v2/blobs/sha256
  -db_host string
        Postgres database host (default "localhost")
  -db_name string
        Postgres database name (default "registry")
  -db_pass string
        Postgres database password (default "root123")
  -db_port int
        Postgres database port (default 5432)
  -db_user string
        Postgres database user (default "postgres")
  -dry_run
        Whether to skip deleting files

```

