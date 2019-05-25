# Docker Swarm REST API

[![CircleCI](https://circleci.com/gh/syntaqx/swarm-api.svg?style=svg)](https://circleci.com/gh/syntaqx/swarm-api)
![GitHub release](https://img.shields.io/github/release/syntaqx/swarm-api.svg)

A dead simple REST API built to solve provisoner configuration issues when
creating a Docker Swarm with [Terraform][].

> See https://github.com/hashicorp/terraform/issues/19509

## Installation

### Download the binary

Quickly download install the latest release:

```sh
curl -sfL https://install.goreleaser.com/github.com/syntaqx/swarm-api.sh | sh
```

Or manually download the [latest release binary][releases] for your
system/architecture and install it into your `$PATH`

### Homebrew on macOS

If you are using [Homebrew][] on macOS, you can install the `swarm-api` with the
following command:

```sh
brew install syntaqx/tap/swarm-api
```

## Usage

By default the `swarm-api` server will only respond to requests on `localhost`.
It is designed to be used on private/protected network interfaces so that newly
created docker hosts can easily connect to a swarm without the need to persist
join tokens while provisioning hosts.

In practice, a DigitalOcean Droplet might start the `swarm-api` server like:

```sh
nohup swarm-api serve --host $(curl http://169.254.169.254/metadata/v1/interfaces/private/0/ipv4/address) &
```

This leverages DigitalOcean's [Droplet Metadata API][metadata-api] and should
work on a Droplet that supports it.

> __⚠ Important__: You should only change the listening host to a
> private/protected network interface. This API is not protected for use on a
> public network and doing so would be a massive security hole in your swarm.
> Please be sure to configure the correct host.

The `swarm-api` is intended to be a running alongside a docker swarm leader or
connection configuration, configured through the environment:

```sh
docker swarm init --advertise-addr $(curl http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address)
```

Once both the swarm is initialzied and the `swarm-api` server started, you're
then able to join new nodes by simply leveraging `curl` or equivelant tool,
specifying `worker` or `manager` as the last parameter in the example below:

```sh
docker swarm join --token $(curl http://$SWARM_LEADER_PRIVATE_IP:8080/swarm/token/worker) $SWARM_LEADER_PUBLIC_IP
```

You are responsible for knowing the values for `$SWARM_LEADER_PRIVATE_IP` and
`$SWARM_LEADER_PUBLIC_IP`, which are generally available during provisioning.

## Security vulnerabilities

If you discover a security vulnerability within the project, please send an
email to Chase Pierce at syntaqx@gmail.com. All security vulnerabilities will be
promptly addressed.

## License

`swarm-api` is open source software released under the [MIT license][MIT].

[MIT]: https://opensource.org/licenses/MIT
[terraform]: https://www.terraform.io/
[homebrew]: https://brew.sh/
[releases]: https://github.com/syntaqx/swarm-api/releases
[metadata-api]: https://developers.digitalocean.com/documentation/metadata/
