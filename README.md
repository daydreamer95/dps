# dps
## Overview: DPS stands for Distributed Priority Queue
```
Inspired: https://engineering.fb.com/2021/02/22/production-engineering/foqs-scaling-a-distributed-priority-queue/
```
Thanks for sharing with community. Using golang to implement this. Same dependency MySql using
in this applications
This is just a small application without database sharding and monolith service contains Dequeue worker, Prefetch buffer, ack and nack worker all in one.

## How to start

```
# For quick start please run this command
docker build --tag dps .
docker-compose up -d

#For local development
cd src && go mod tidy
go run src/.


Enjoy and have fun.
Make file update soon
```

## Authorize
```
Me, my self and I. On learning purpose lead me here
```

## Contributing
```
Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.
```
## License


