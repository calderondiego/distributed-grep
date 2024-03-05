# CS425 MP1: Distributed Log Querier

### Building
`make build`

### Running
To run the server: `./mp1 server -port 4250`
To run the client: `./mp1 client -servers servers.txt -command "grep -c 191.251.168.5 *.log"`

### Manual
```
$ ./mp1 server -help
Usage of server:
  -port int
    	port to listen on (default 4250)

$ ./mp1 client -help
Usage of client:
  -command string
    	grep command to run on server
  -servers string
    	file containing list of servers (default "servers.txt")
```
