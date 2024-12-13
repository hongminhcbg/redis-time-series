# service-gin-template

## Run demo

```sh
docker-compose up # init redis time seris
go run cmd/* server # init app
```

## Request to init data time series

```sh
curl --location --request POST 'http://localhost:8080/api/v1/velocity-rule-in' \
--header 'User-Agent: Apidog/1.0.0 (https://apidog.com)' \
--header 'Content-Type: application/json' \
--data-raw '{
    "trans_type": "A",
    "data": {
        "a": {
            "client_id": "xxx",
            "casa_id": "xxx",
            "amount": 1025,
            "timestamp": 1000
        },
        "b": {
            "client_id": "xxx",
            "casa_id": "xxx",
            "amount": 1009,
            "timestamp": 1000
        },
        "c": {
            "client_id": "xxx",
            "casa_id": "xxx",
            "amount": 10005,
            "timestamp": 1000
        }
    }
}'

--- resp ---
{
    "output": [
        {
            "UserId": "1",
            "DataKey": "a",
            "Ts": "1000",
            "Data": 1025
        },
        {
            "UserId": "1",
            "DataKey": "a2",
            "Ts": "1000",
            "Data": 1025
        }
    ],
    "status": "ok"
}
```

## Exec into redis and check redis time series data

```sh
ts.range 1a2 1734081093978 1734081130788 AGGREGATION sum 1000
# resp 
# 127.0.0.1:6379> ts.range 1a2 1734081093978 1734081130788 AGGREGATION sum 1000
#1) 1) (integer) 1734081093000
#   2) 1003
#2) 1) (integer) 1734081106000
#   2) 1005
#3) 1) (integer) 1734081113000
#   2) 1007
#4) 1) (integer) 1734081120000
#   2) 1009
#5) 1) (integer) 1734081130000
#   2) 102

```
