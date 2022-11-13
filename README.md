# uploadfile-test

## Setup
### Requirements
- Install [Docker Engine](https://docs.docker.com/engine/install/ubuntu/)
- Install [Docker compose](https://docs.docker.com/compose/install/other/)
- Install `make`: `sudo apt install build-essential`
- Install [go](https://golang.org/dl/), version at least `1.19.1`
- Install [vegeta](https://github.com/tsenart/vegeta#source)

## Explain my application:
### Concept
- Our application will serve a api, user can call this api to send some data, we will forward these data to gcp or s3.
- The payload size limit is 10KB
- API can still handle/recieve requests from clients without having to wait for the uploading process to complete.

=> For these reasons above, I decided use [Nats-Jetstream](https://docs.nats.io/nats-concepts/jetstream) pub/sub in this project. 
It is same like queue, it has publisher and subscriber. So when user make a request, the data will be published by publisher to the queue, 
subscriber listens the queue, then the data will be processed by subscriber and send it to gcp/s3. Btw, because I don't have gcp account,
so I used [mockery](https://github.com/vektra/mockery) to generate the mock when we upload data to gcp.

- Pros:
    - User don't need to wait the uploading process complete.
    - Easy to scale, we can easy to scale the subscriber more resource to process the logic faster.
    - If publisher or subscriber is crash/down, the data are still storing in nats, then we don't lost the data.
- Cons:
    - It makes our application quite complexity.

![follow](https://user-images.githubusercontent.com/36435846/201523208-9fe4c101-65c9-42c4-9d54-1bbe5251cd20.png)


### Skeleton project

```
├── configs // variable environment
├── docker-compose.yaml // run program with docker
├── files // payload to run api
├── internals
│   ├── logs // store logs into file for tracing.
│   ├── transport
│   │   ├── route.go // the api and middleware to serve
│   │   └── route_test.go // unit test of route
│   ├── usecase // all logics at here 
│   │   ├── entities
│   │   ├── mocks // we use mockery to generate these files
│   │   ├── publisher
|   │   │   ├── publisher_test.go // the unit test of publisher logic
|   │   │   └──publisher.go // publish message to nats
│   │   └── subscriber
|   │   │   ├── gcp.go // upload data to gcp
|   │   │   ├── s3.go // upload data to s3
|   │   │   ├── subscriber_test.go // the unit test of subscriber logic
|   │   │   ├── subscriber.go // receive message from nats then call gcp/s3 to upload data.
|   │   │   └── upload-file.go // interface to define upload method.
│   ├── util 
│   │   ├── config_test.go // the unit test of config logic
│   │   ├── config.go // define some config fields of service
│   │   ├── constants.go // almost elements of nats
│   │   └── nats.go // interface and nats implement.
└── main.go
```

## Run the application:
### In order to run the application:
1. Run command `make init` to start the nats and application. 
2. Run command `make call-api-with-normal-payload` to make a request to server with 8080 port. I can go to `logs.log` to see the process.
```
2022-11-13T10:18:10.603Z	INFO	publisher/publisher.go:32	UploadPublish: Publish msg success	{"msg_id": "c39f8d06-9e4f-4e40-afaa-635b76782840"}
2022-11-13T10:18:10.607Z	INFO	util/nats.go:114	QueueSubscribe: Ack msg success	{"sequence_id": 180223, "msg_id": "c39f8d06-9e4f-4e40-afaa-635b76782840"}
```
3. Run command `make call-api-with-heavy-payload` to make a request with heavy payload, then we got this log:
```
2022-11-13T10:18:06.464Z	ERROR	transport/route.go:36	Payload is over 10KB
```
4. Run command `make load-testing` to load testing our api:
```
vegeta attack -targets=./files/target.json -duration=120s -rate=0 -max-workers=3 | tee results.bin | vegeta report
Requests      [total, rate, throughput]         174220, 1451.80, 1451.79
Duration      [total, attack, wait]             2m0s, 2m0s, 1.181ms
Latencies     [min, mean, 50, 90, 95, 99, max]  230.92µs, 2.042ms, 1.348ms, 4.199ms, 5.968ms, 11.381ms, 68.868ms
Bytes In      [total, mean]                     5749260, 33.00
Bytes Out     [total, mean]                     595832400, 3420.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:174220 
```
- My system information:
```
PRETTY_NAME="Ubuntu 22.04.1 LTS"
NAME="Ubuntu"
VERSION_ID="22.04"
VERSION="22.04.1 LTS (Jammy Jellyfish)"
VERSION_CODENAME=jammy
ID=ubuntu
ID_LIKE=debian
RAM=32G
CPU=8core
```

5. Run command `unit-test` for unit test:
```
go test -count=3 ./internals/...
?       github-com/edarha/uploadfile-test/internals/logs        [no test files]
ok      github-com/edarha/uploadfile-test/internals/transport   0.018s
?       github-com/edarha/uploadfile-test/internals/usecases/entities   [no test files]
?       github-com/edarha/uploadfile-test/internals/usecases/mocks      [no test files]
ok      github-com/edarha/uploadfile-test/internals/usecases/publisher  0.017s
ok      github-com/edarha/uploadfile-test/internals/usecases/subscriber 0.016s
ok      github-com/edarha/uploadfile-test/internals/util        0.011s
```
6. In order to monitor your nats. You can use [nats-box](https://github.com/nats-io/nats-box)
Run it by docker `docker run --rm --network host -it natsio/nats-box:latest`

Then we will connect to nats-server by this command: `nats context save s1 --user=admin --password=admin --server=nats://127.0.0.1:4223 --select`

After connect success, we can run `nats --help` to use some commands.
```
LAP00335:~# nats str info upload
Information for Stream upload created 2022-11-13T09:55:18Z

Configuration:

             Subjects: upload.send
     Acknowledgements: true
            Retention: File - Interest
             Replicas: 1
       Discard Policy: Old
     Duplicate Window: 2m0s
    Allows Msg Delete: true
         Allows Purge: true
       Allows Rollups: false
     Maximum Messages: unlimited
        Maximum Bytes: unlimited
          Maximum Age: unlimited
 Maximum Message Size: unlimited
    Maximum Consumers: unlimited


State:

             Messages: 0
                Bytes: 0 B
             FirstSeq: 180,224
              LastSeq: 180,223 @ 2022-11-13T10:18:10 UTC
     Active Consumers: 1

LAP00335:~# nats con report upload
╭───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│                                                          Consumer report for upload with 1 consumers                                                          │
├──────────────┬──────┬────────────┬──────────┬─────────────┬─────────────┬─────────────┬───────────┬───────────────────────────────────────────────────────────┤
│ Consumer     │ Mode │ Ack Policy │ Ack Wait │ Ack Pending │ Redelivered │ Unprocessed │ Ack Floor │ Cluster                                                   │
├──────────────┼──────┼────────────┼──────────┼─────────────┼─────────────┼─────────────┼───────────┼───────────────────────────────────────────────────────────┤
│ upload-queue │ Push │ Explicit   │ 30.00s   │ 0           │ 0           │ 0           │ 180,223   │ NAQT6UPEWCFNJES4PVKEJYSKEVY7UTB35FOPAVUZ22JOPM6XZ6NGQBHB* │

```