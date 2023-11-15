# Go RabbitMQ Exercise
- Publish/Subscribe
- QoS
- Topic
- Durable messages
- Publish confirm
- Consume confirm

## Run
```bash 
cd go-rabbitmq-exercise

# Initialize apps dependencies.
go mod download

# Run RabbitMQ instance.
docker compose up -d

# Run consumer.
run go consumer.go

# Run producer.
run go producer.go
```