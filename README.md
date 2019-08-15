[![Build Status](https://travis-ci.org/bwNetFlow/dashboard.svg?branch=master)](https://travis-ci.org/bwNetFlow/dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/bwNetFlow/dashboard)](https://goreportcard.com/report/github.com/bwNetFlow/dashboard)
[![GoDoc](https://godoc.org/github.com/bwNetFlow/dashboard?status.svg)](https://godoc.org/github.com/bwNetFlow/dashboard)
[![Docker Pulls](https://img.shields.io/docker/pulls/bwnetflow/dashboard.svg)](https://hub.docker.com/r/bwnetflow/dashboard/)

# Dashboard Consumer

This Kafka consumer reads enriched bwNetFlow flows and exports relevant
information to a Prometheus database. Relevant information means only a subset
of information from the flows, e.g. aggregated counter values etc.

## Requested Dashboard Metrics

For one customer dashboard (flows primarily are filtered by cid).

* Bytes / Packets / Flows Total
* Bytes By well known Ports (http, e-mail, etc)
* Bytes By Protocol
* Bytes By IP Version
* Bytes By Peer

Rank based:

* Top N IP addresses by Bytes
* Top N IP addresses by Connections
