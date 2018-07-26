# Project Unclog
![alt](https://travis-ci.org/go-squads/unclog-worker.svg?branch=master)

## Description
This project is integrated with BaritoLog found on this link: https://github.com/BaritoLog/

This project will create a stream processor which will process logs according to the contract (TimberWolf) and will send the processed logs into a new topic that will later be stored in a database. 


## Environment
- Golang

### How to install Golang for MacOS
Using Homebrew:
```
brew install go
```
### GoLang Dependencies
```
go get
```

## Services used
- Kafka
- InfluxDB
### How to install Kafka
Download kafka from https://kafka.apache.org/downloads and unzip the folder

### How to install InfluxDB for MacOS
- Note: for other binaries, follow instructions on https://portal.influxdata.com/downloads


Using Homebrew:
```
brew install influxdb
```

## Installation and Setup of Project Unclog
```
cd $GOPATH/src
```
```
git clone git@github.com/go-squads/unclog-worker
```
```
cd unclog-worker
```
```
go install
```

## Running
- Note: To kill processes, use control + c key, make sure you don't just exit the terminal


Run zookeeper (from kafka directory)
```
./bin/zookeeper-server-start.sh config/zookeeper.properties
```
Run kafka-cluster (from kafka directory)
```
./bin/kafka-server-start.sh config/server.properties
```
Run Unclog worker
```
$GOPATH/bin/unclog-worker
```
## Troubleshooting

### Killing a process in the background
```
lsof -i @localhost:[port]
```
```
kill -9 [pid]
```
