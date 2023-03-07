# Whoami DNS

whoami-dns is a clever webserver that runs the backend for https://debug.charliejonas.co.uk/dns.html. It's designed to work in tandem with an authoritative DNS server using [dnstap](https://dnstap.info) and wildcard domains so that clients' recursive DNS servers can be identified by the source IP addresses of their queries.

## Principle

1. A zonefile contains a `* 1 A 192.0.2.1` wildcard record where 192.0.2.1 is the public IP address of the host.
2. The name server that is authoritative for that zone is configured to use dnstap via a UNIX socket.
3. whoami-dns listens on that socket and detects incoming DNS queries.
4. The source IP address and domain of the query are parsed and stored by whoami-dns.
5. whoami-dns listens for incoming HTTP requests and uses the `Host` header to determine the source IP address of the DNS query.
6. A plaintext HTTP response body is sent by whoami-dns containing the client's recursive resolver's IP address.

## Usage

```
Usage:
  whoami-dns [flags]

Flags:
  -b, --bind string   path to dnstap UNIX socket (default "/var/lib/knot/dnstap.sock")
  -h, --help          help for whoami-dns
  -p, --port string   port on which to listen for HTTP requests (default "6780")
```

## Installation

Pre-built binaries for a variety of operating systems and architectures are available to download from [GitHub Releases](https://github.com/CHTJonas/whoami-dns/releases). If you wish to compile from source then you will need a suitable [Go toolchain installed](https://golang.org/doc/install). After that just clone the project using Git and run Make! Cross-compilation is easy in Go so by default we build for all targets and place the resulting executables in `./bin`:

```bash
git clone https://github.com/CHTJonas/whoami-dns.git
cd whoami-dns
make clean && make all
```

## Copyright

whoami-dns is licensed under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).

Copyright (c) 2021 Charlie Jonas.
