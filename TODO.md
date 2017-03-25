A collection of smaller and larger tasks that would be nice to implement in this project.

* Persistent TCP connections
  - Use protocol header to signal how much data is being transmitted

* Asymmetrical payload size
  - Specify the payload size in the header
  - Specify the requested return size in header
  - Server should take config parameter stating the maximum allowed size, the actual size retured
    should be min of these.

* Different ping intervals for UDP and TCP
  - Two new config parameters

* Online configuration of probes through an API to be able to dynamically adjust the traffic pattern.
  New command "pingo ctl" to manage this. Read and write configuration.
  
* Metrics collector with plugins for different metric sinks (InfluxDB, Prometheus, Graphite, ...).
  New command "pingo collect" for this.
  
* TLS support with server and client certificates, both between server and probe and probe and API
  clients.
  
* White listing of IP addresses that may access servers and probes.

* Add "server less" pinging against arbitrary HTTP(s) or ICMP endpoints that do not have to run pingo servers.

Request header format v1
------------------------
A 20 byte header will always be present 

```
XXXXXXXX|YYYY|ZZZZ|BBBB
 |       |    |    |
 |       |    |    - Config bitmap, 4 bytes (descibed below)
 |       |    - Requested response content len, 4 bytes
 |       - Request content len, 4 bytes
 - Magic number, 8 bytes (also identifies header version)
 
 B0 - Should the connection remain open after response has been sent (eg. be persistent)? 
      1 - True, 0 - False
 ```
 
Response header format v1
-------------------------
A 28 byte header will always be present
```
XXXXXXXX|YYYY|ZZZZZZZZZZZZZZZZ
 |       |    |
 |       |    |
 |       |    - MD5 checksum of request payload
 |       - Response content len, 4 bytes
 - Magic number, 8 bytes (also identifies header version)
```