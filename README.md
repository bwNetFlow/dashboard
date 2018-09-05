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
