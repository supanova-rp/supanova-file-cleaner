# Supanova File Cleaner

Go service that runs periodically to remove unused files from Supanova learning platform S3 bucket

Setup:
```
make dep
```

Run:
```
make run
```

Generate db queries:
```
make sqlc
```

Run with docker:
```
make docker/local-build
make docker/local-run
```

