# Meer
<img align="right" width="100px" src="https://github.com/roychowdhuryrohit-dev/projectmeer/blob/main/assets/logo.png">

 An distributed truly peer-to-peer collaborative text editor based on Conflict-Free Replicated Data Type (CRDT) and built as a part of final project for Distributed Systems (CSEN 371) at Santa Clara University.

 ## Features

  - Uses [FugueMax CRDT](https://arxiv.org/abs/2305.00583) to minimize the issue of interleaving.
  - [Causal Broadcast](https://www.youtube.com/watch?v=A8oamrHf_cQ) algorithm is used to ensure eventual consistency.
  - React frontend is built with [Lexical](https://lexical.dev) framework to handle editor state changes seamlessly.
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


<img align="left" width="550px" src="https://github.com/roychowdhuryrohit-dev/projectmeer/assets/24897721/7d50f168-939e-44f6-b43d-1adc399c2e99">


