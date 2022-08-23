# Golang Training for OWT

Steps to test:
1. Run docker containers:
```sh
$ docker-compose up -d
```
2. Build the application:
```sh
$ make build
```
3. Define cadence domain (wait until docker containers are fully deployed):
```sh
$ make register 
```
4. Run the application:
```sh
$ make run
```

When the application is started, it can be tested using the curl.sh file provided in the repository:
```sh
$ ./curl.sh
```
Use keys 1-6 which map to the following curl calls:
1. call the start Fibonacci endpoint with input 100
2. call the polling endpoint with the last returned id from 1. or 4.
3. call the get result endpoint with the last returned id from 1. or 4.
4. call the start Fibonacci endpoint and provide your own input
5. call the polling endpoint and provide the id yourself
6. call the get result endpoint and provide the id yourself
