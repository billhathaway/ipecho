ipecho is a service to return a caller's IP and hostname.

This is useful in cases where your host is behind a NAT or in an environment like OpenStack where it may not be able to determine locally what the external IP address is.

Endpoint
/ - returns the IP address on a single line

By default the process listens on 0.0.0.0:8080

To change the listen address or port, use
-host=<IP or name>
-port=< port>

Example:
./ipecho -host 127.0.0.1 -port 8888


Simple logging is now enabled by default, to disable it use
-quiet
