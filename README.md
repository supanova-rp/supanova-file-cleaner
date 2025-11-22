# Supanova Maintenance

Go service that runs periodic maintenance tasks for Supanova learning platform (e.g. remove unused files from S3, backup DB to S3)

#### Setup:
```
make dep
```

#### Run:
```
make run
```

#### Run with docker:
```
make docker/run
```

#### Generate db queries:
```
make sqlc
```
