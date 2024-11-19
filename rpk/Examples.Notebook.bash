# Examples 

## clsuter Health
```shell
rpk cluster health
rpk redparpk redpanda admin brokers list

```

### create a topic 
```shell
rpk topic create win-some -p 6 -r 1 
rpk topic create some-more -p 6 -r 1
```

### List partitions
```shell
rpk redpanda admin partitions list 2
```

# Run Single Connect Service 
```shell
clear ; rpk connect run generate.yaml
```

# Run Rpk in streams mode 
```shell
clear; rpk connect streams server.yaml streams/*.yaml 
```

# Json 
```shell
curl http://localhost:18082/topics 
```

# List Partitions
```shell
curl -s \
   -X 'GET' \
  'http://localhost:18082/offsets' \
  -H 'accept: application/vnd.kafka.v2+json' \
  -H 'Content-Type: application/vnd.kafka.v2+json' \
  -d '{
  "partitions": [
    { "topic": "win-some", "partition": 0 },
    { "topic": "win-some", "partition": 1 },
    { "topic": "win-some", "partition": 2 },
    { "topic": "win-some", "partition": 3 },
    { "topic": "win-some", "partition": 4 },
    { "topic": "win-some", "partition": 5 },
    { "topic": "some-more", "partition": 0 },
    { "topic": "some-more", "partition": 1 },
    { "topic": "some-more", "partition": 2 },
    { "topic": "some-more", "partition": 3 },
    { "topic": "some-more", "partition": 4 },
    { "topic": "some-more", "partition": 5 },
  ]
}'

```