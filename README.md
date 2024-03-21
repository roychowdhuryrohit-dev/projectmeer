# Meer

 An distributed peer-to-peer collaborative text editor built as a part of final project for Distributed Systems (CSEN 371) at Santa Clara University.

 ## Features

  - Uses FugueMax CRDT to minimize the issue of interleaving.
  - Causal Broadcast algorithm is used to ensure eventual consistency.
  - Frontend is built with Lexical Text Editor framework to handle editor state changes smoothly.
  - Gracefully shuts down the server when OS signals like _SIGINT_ & _SIGTERM_ are caught. The maximum amount of time taken by the server to wait for pending requests is configurable.

## Usage

Make sure `go` is installed and set in the path.
```
$ brew install go
```
To build the project locally, 
```
$ make build_local
```
An executable will be generated in `bin/` directory. To run it,
```
# run slug server and visit localhost:8080 on browser
$ ./bin/slug -document_root=/Users/www/scu.edu -port=8080 -timeout=5
```
To build docker image.
```
$ make build_docker
```
