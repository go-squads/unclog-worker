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
- PostgreSQL 10.4
### How to install Kafka
Download kafka from https://kafka.apache.org/downloads and unzip the folder

### How to install PostgreSQL for MacOS
- Note: 


Using Homebrew:
```
brew install postgresql
```

### Login PostgreSQL
```
sudo -u postgres psql
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
Run Unclog Migration (setting up database)
* Note: Create ```config.yaml``` according to your environment, follow ```config.yaml.example```
```
$GOPATH/bin/unclog-worker migrate
```
Run Unclog Worker Stream Processor
```
$GOPATH/bin/unclog-worker sp
```
Run Unclog Worker Log Count Stream Processor
```
$GOPATH/bin/unclog-worker splc
```
## Troubleshooting

### Killing a process in the background
```
lsof -i @localhost:[port]
```
```
kill -9 [pid]
```
### Clearing kafka logs
```
cd ~/../../tmp
```
```
rm -rf kafka-logs
```
```
rm -rf zookeeper/version-2
```
