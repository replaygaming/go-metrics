# Game metrics service
[![Build Status](https://semaphoreci.com/api/v1/projects/dba15a7d-a543-4860-b8c0-a6b64d15b840/563329/shields_badge.svg)](https://semaphoreci.com/paulgould/go-metrics) [![Coverage Status](https://coveralls.io/repos/replaygaming/go-metrics/badge.svg?branch=master&service=github)](https://coveralls.io/github/replaygaming/go-metrics?branch=master)

Translates Replay Poker events and forward them to 3rd-party APIs

## Integrations supported

  - [Amplitude](http://www.amplitude.com)

## Usage

### Pre-built Binary
Get the latest binary for your [distribution](https://github.com/replaygaming/go-metrics/releases)

### Building from source

####  Get project dependencies

```shell
export GOPATH=~/go
go get github.com/replaygaming/go-metrics
cd ~/go/src/github.com/replaygaming/go-metrics
go get .
```

#### Compile

```
go build
```

#### Run

```shell
./go-metrics
```

#### Configuration

```shell
# AMQP URL (for example "amqp://guest:guest@127.0.0.1:5672/metrics")
AMQP_URL

# AMQP Queue name (for example "metrics")
AMQP_QUEUE

# Amplitude API key
AMPLITUDE_API_KEY
```

## Configure RabbitMQ

### Install `rabbitmq` and `rabbitmqadmin`

Download and installation guide from [RabbitMQ site](https://www.rabbitmq.com/download.html).
rabbitmqadmin is binary, found as part of [rabbitmq-management](https://github.com/rabbitmq/rabbitmq-management) project.

### Enable the management plugin:
Enable rabbitmq_management plugin.

```
    [sudo] rabbitmq-plugins enable rabbitmq_management
```

Then (re)start the rabbitmq daemon.

```
    [sudo] sudo rabbitmqctl stop
    [sudo] rabbitmq-server -detached
```

Declare the host and exchange for the metrics

```shell
rabbitmqadmin declare vhost name=metrics
rabbitmqadmin declare permission vhost=metrics user=guest configure=".*" write=".*" read=".*"
rabbitmqadmin -V metrics declare exchange name=metrics_ex type=fanout durable=true
```

## Development Resources

### Go

[Go Installation](https://golang.org/doc/install)
[Go Code Documentation](https://golang.org/doc/code.html)
[Go + Docker](https://blog.golang.org/docker)

### Docker

[Docker Installation](https://docs.docker.com/engine/installation/)
[Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)

### Codeship

[Codeship Steps Configuration](https://codeship.com/documentation/docker/steps/)
[Codeship Services Configuration](https://codeship.com/documentation/docker/services/)

## Contributing

We would love to see contributions from the community. Please feel free to raise an issue or send your PR to this Github project.
