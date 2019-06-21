# Messaging Realtime Chat
Messaging service is used by created simple message and retrieve message implementation.

## Cases
1. User Post a massage use API
2. Message will be send to websocket and store in database
3. Message `status` should be `delivered` if message success sent to websocket.

## Quick Start
### Technology Use
* Database Mysql
* Websocket (for data streaming)

### Installation Required
#### Install Golang
 I use go version for service compatible `go1.11.1 darwin/amd64`
#### Install Mysql
I common mysql

### Deploy and run
The first you must extract this file into your workspace project.

How To run :

    Run Application Server
    $ cd app
    $ go run main.go

    Run Websocket Server
    $ cd websocket
    $ go run main.go

## Documentation
This service has documentation following swagger for your try. So please se file in path ./doc/api_spech.yaml.

