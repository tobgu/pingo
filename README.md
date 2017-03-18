Pingo
=====
Simple application level network probing for TCP and UDP.

It consists of two parts:
* An echo server that returns whatever data is sent to it over TCP and UDP .
* A probe that sends traffic to selected echo servers, measures the round trip time and publishes the results over an HTTP API.

Run Server
----------
`pingo server --config=server-config.yaml`

Run Probe
---------
`pingo probe --config=probe-config.yaml`

Configuration
--------------
Commented example configuration files for server and probe are available under `examples/`

Licence
=======
MIT