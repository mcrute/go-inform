# Ubiquiti Inform Protocol in Golang
This repo contains a Golang implementation of the Ubiquiti Networks
Inform protocol used by the Unifi access points and the mFi machine
networking components. The primary purpose of this repository is to
implement the core inform protocol for interoperability with Ubiquiti
products. It is not recommended to use this protocol as the base for
any new product development, it's rather Ubiquiti specific and not very
secure in the bootstrap phase (uses well-known keys over plaintext
channels).

This work is largely based on 
[ubntmfi](https://github.com/mcrute/ubntmfi/blob/master/inform_protocol.md)

## Status
This repo is a work in stasis and is semi-maintained, it is not yet
considered stable. Though the repository is archived, if there is ever
interest in this project I will un-archive it and find time to more
actively maintain it. If you find it useful patches are welcome, just
open a pull request.

This repository should work well with older mFi components and modern
Unifi components as well. The author has tested with devices as new as
the U6-Mesh devices.

## Missing Features
If you need these then feel free to implement them. Pull requests are
accepted.

- Support for writing AES GCM packets
- Support for writing zlib compressed packets
- Support for writing snappy compressed packets

## Protocol
The inform protocol works over HTTP or HTTPS in the case of modern
equipment, older equipment did not support HTTPS. The `mcad` daemon
on the devices is responsible for speaking the inform protocol to the
controller.

A device will `POST` an HTTP request to the controller with a
content-type of `application/x-binary` which contains a payload encoded
in inform format. The server will respond with a payload that is inform
encoded.

## Packet Format
Packets are binary and transmitted in big-endian format.

| Size     | Type     | Purpose                  | Notes                                |
| -------- | -------- | ------------------------ | ------------------------------------ |
| 4 bytes  | `int32`  | Protocol Magic Number    | Must always be `1414414933` (`UBNT`) |
| 4 bytes  | `int32`  | Packet Version           | Currently this is `0`                |
| 6 bytes  | `[]byte` | Device MAC Address       | Used for crypto-key lookup           |
| 2 bytes  | `int16`  | Flags                    | See below                            |
| 16 bytes | `[]byte` | Encryption IV            |                                      |
| 4 bytes  | `int32`  | Data Version             | Currently this is `1`                |
| 4 bytes  | `int32`  | Encrypted Payload Length |                                      |
| n bytes  | `[]byte` | Encrypted Payload        | See below                            |

### Flags
| Flag | Name              | Purpose                                         |
| ---- | ----------------- | ----------------------------------------------- |
| `1`  | Encrypted         | Indicates that the payload is encrypted         |
| `2`  | Zlib Compressed   | Indicates that payload is zlib compressed       |
| `4`  | Snappy Compressed | Indicates that payload is snappy compressed     |
| `8`  | GCM Encrypted     | Indicates that packet is encrypted with AES GCM |

### Encryption
There are two encryption modes AES 128 CBC and AES 128 GCM. GCM is used
in newer devices and CBC in older devices. The same key is used for
either mode. It's stored as `x_authkey` in the Unifi database and is
encoded as hex. The key must be decoded before decryption.

AES GCM requires authentication data to decrypt the packet. The
authentication data is the following fields encoded in big-endian binary
format.

- Protocol Magic Number
- Packet Version
- Device MAC Address
- Flags
- Encryption IV
- Data Version
- Data Length

In CBC mode the data is padded to fit a full AES block. In GCM mode the
final block appears to be a padding/garbage block.

### Payload
The payload is a JSON string that is device and application specific.
The payload may also (but is not required to be) be compressed using
either [snappy](https://github.com/golang/snappy) or zlib compression as
indicated by the flags.
