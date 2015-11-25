# Game metrics service
[![Build Status](https://semaphoreci.com/api/v1/projects/dba15a7d-a543-4860-b8c0-a6b64d15b840/563329/shields_badge.svg)](https://semaphoreci.com/paulgould/go-metrics) [![Coverage Status](https://coveralls.io/repos/replaygaming/go-metrics/badge.svg?branch=master&service=github)](https://coveralls.io/github/replaygaming/go-metrics?branch=master)

Translates Replay Poker events and forward them to 3rd-party APIs

## Integrations supported

  - [x] [Amplitude](http://www.amplitude.com)

## Usage

```shell
./bin/metrics -h

  -amplitude-api-key string
        Amplitude API Key
  -amqp-queue string
        AMQP Queue name (default "metrics")
  -amqp-url string
        AMQP URL (default "amqp://guest:guest@localhost:5672/metrics")
  -newrelic-license string
        NewRelic License Key (default "")
  -newrelic-app string
        NewRelic App Name (default "")
```

## Configure RabbitMQ

### Install `rabbitmq` and `rabbitmqadmin`

Download and installation guide from [RabbitMQ site](https://www.rabbitmq.com/download.html).
rabbitmqadmin is binary, found as part of [rabbitmq-management](https://github.com/rabbitmq/rabbitmq-management) project.

### Enable the management plugin:

    [sudo] rabbitmq-plugins enable rabbitmq_management

Then (re)start the rabbitmq daemon.

    [sudo] sudo rabbitmqctl stop
    [sudo] rabbitmq-server -detached

Declare the host and exchange for the metrics:

    rabbitmqadmin declare vhost name=metrics
    rabbitmqadmin declare permission vhost=metrics user=guest configure=".*" write=".*" read=".*"
    rabbitmqadmin -V metrics declare exchange name=metrics_ex type=fanout durable=true

## Contribuing

### Install `go`

Follow the instructions at [Golang.org](https://golang.org). **DO NOT** install using your distro pkg manager.

### Get project dependencies

    go get .

### Running

    make
    cd bin
    LD_LIBRARY_PATH=lib ./metrics
