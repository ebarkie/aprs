# APRS

Go package for working with APRS string and byte packets.  It can send those
packets via APRS-IS or transmit them via TNC KISS.

## Installation

```
$ go get github.com/ebarkie/aprs
```

## Usage

See [USAGE](USAGE.md).

## Examples
  - Position report
    - [Sending](examples/position_send)
  - Weather observation
    - [Sending](examples/wx_send)
    - [Receiving](examples/wx_recv)

## License

Copyright (c) 2016-2020 Eric Barkie.  All rights reserved.  
Use of this source code is governed by the MIT license
that can be found in the [LICENSE](LICENSE) file.
