# Game metrics service

[ ![Codeship Status for replaygaming/go-metrics](https://codeship.com/projects/2d93ed00-10a6-0134-1a45-32602de4173e/status?branch=kubernetes/master)](https://codeship.com/projects/157065)

Translates Replay Poker events and forward them to 3rd-party APIs

## Integrations supported

- [Amplitude](http://www.amplitude.com)

## Usage

### Building from source

We use [`glide`](http://glide.sh/) package manager for go. [Install](https://glide.readthedocs.io/en/latest/) the latest version of `glide`.

####  Get project dependencies

    export GOPATH=~/go
    go get github.com/replaygaming/go-metrics
    cd ~/go/src/github.com/replaygaming/go-metrics

Initialize and install dependencies

    glide install		# generates glide.lock

Glide will install all your go project dependencies in `vendor/` directory.

When you update the `glide.yaml` to update the dependencies or add specific `glide` configuration manually, you should update the `glide.lock` file, by using following command.

    glide update        # updates the glide.lock


#### Compile

    go build -v

#### Configuration

Configuration is done using environment variables.

    # Topic (defaults to "metrics")
    export METRICS_TOPIC="my_topic"

    # Subscription (defaults to "metrics_workers")
    export METRICS_SUBSCRIPTION="my_subscription"

    # Amplitude API key. Required
    export AMPLITUDE_API_KEY="asdf1234"

#### Run


    ./go-metrics

## Development Resources

### Go

- [Go Installation](https://golang.org/doc/install)
- [Go Code Documentation](https://golang.org/doc/code.html)
- [Go + Docker](https://blog.golang.org/docker)
- [Glide](http://glide.sh)

### Docker

- [Docker Installation](https://docs.docker.com/engine/installation/)
- [Dockerfile Reference](https://docs.docker.com/engine/reference/builder/)

### Codeship

- [Codeship Steps Configuration](https://codeship.com/documentation/docker/steps/)
- [Codeship Services Configuration](https://codeship.com/documentation/docker/services/)

## Contributing

We would love to see contributions from the community. Please feel free to raise an issue or send your PR to this Github project.
