# calculator-client

 A simple calculator program that calculate simple math tasks by getting input via program arguments and send it to the [Server](https://github.com/JosephmBassey/calculator-service). 
## Clone  the project

```
$ git clone https://github.com/JosephmBassey/calculator-client
$ cd calculator-client
```
## Set the calculator server environment variables.
 Assuming the calculator server is running on `0.0.0.0:9081`
```
$ export CALCULATOR_GRPC_SERVER=0.0.0.0:9081
```

## build and run the project.
- make sure the calculator  [Server](https://github.com/JosephmBassey/calculator-service).  is running.
- still in the calculator-client dir.
```
$ make build
$ ./client -method add -a 2 -b 5 
```