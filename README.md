
# Usage: tools to cleanup the blob in the docker registry

## Prerequest

1. Mount the current nfs storage to the current VM
2. If the database is in the k8s cluster, forward the db container's 5432 to localhost
```
kubectl get pods -n tanzu-system-registry
# Get the database password
kubectl exec -it harbor-database-0 -n tanzu-system-registry -- env |grep PASS
kubectl port-forward harbor-database-0  5432:5432 -n tanzu-system-registry
```

3. Change the harbor into read only mode, backup the NFS folder and database.
```
create table artifact_blob_backup as select * from artifact_blob;
```
4. Run the goloang program to cleanup the blob in the registry
```
go run ./cmd/cleanup/main.go -base_dir=/var/lib/registry/docker/registry/v2 -dry_run=true
go run ./cmd/cleanup/main.go -base_dir=/var/lib/registry/docker/registry/v2 -dry_run=false

```

5. cleanup redis cache in the registry

```
kubectl exec -it harbor-redis-0 -n tanzu-system-registry -- redis-cli
redis-cli>
>select 2
# make sure most of the keys are start with blobs::sha25
>keys *
>flushdb
>dbsize
```
6. Try to push and pull images in the registry
