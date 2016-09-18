Ubiquiti Inform Protocol in Golang
==================================
This repo contains a Golang implemntation of the Ubiquiti Networks Inform
protocol used by the Unifi access points and the mFi machine networking
components. The primary focus is an implemenation compatible with mFi
components but the library should also be compatible with Unifi devices.

This repo is a work in progress and is not yet considered stable. If you find
it useful patches are welcome, just open a pull request.

There is a feature-complete Python version of this API, a set of tools useful
for reverse engineering the protocol and an in-progress protocol spec at
[ubntmfi](https://github.com/mcrute/ubntmfi/blob/master/inform_protocol.md)
