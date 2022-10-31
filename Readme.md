# Simple HTTP Mock

Simple http mock sample, take it as it is, no real functionality.


# Mock the API

## Build and start

```shell
cd src
go build main.go
./main -p 9999
```

## Usage and examples

```shell
Usage of ./main:
  -c, --howmany int32           HowMany users mock will return. Default:1000 (default 1000)
  -l, --myserverstatus string   Status of the API. Default:up (default "up")
  -p, --port int32              TCP port for the HTTP listener to bind to. Default: 8080 (default 8080)
  -u, --usercheck string        User to check if exists. Default:gigel (default "gigel")
  -r, --userresponse string     Response value checking user. Default:true (default "true")
```


```shell
./main -p 8080 -c 100 -l up -u gigel -r true
```

This will return:

```shell
curl -X 'GET' http://localhost:8080/v1/usercount -H 'accept:application/json'
{"code":200,"id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","message":"Success","data":100}

curl -X 'GET' http://localhost:8080/v1/status -H 'accept:application/json'
{"code":200,"id":"f21609a2-643a-4dc4-9c30-7e63c08d8283","message":"Success","data":{"ServerStatus":"up","ProcessId":1234}}

 curl -X 'GET' http://localhost:8080/v1/usercheck/gigel -H 'accept:application/json'
{"code":200,"id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","message":"Success","data":"true"}
      
```


## Curl

Query user

```shell
curl -X 'GET' http://localhost:8082/api/v1/usercheck/test -H 'accept:application/json'
{"Code":"200","Id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","Data":"false","Message":"Success"}
```

Count users

```shell
curl -X 'GET' http://localhost:8082/api/v1/usercount -H 'accept:application/json'
{"Code":"200","Id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","Data":999999,"Message":"Success"}%
```

Check server status

```shell
curl -X 'GET' http://localhost:8082/api/v1/status -H 'accept:application/json'
{"Code":"200","Id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","Data":{"ServerStatus":"up","ProcessId":1234},"Message":"Success"}
```

