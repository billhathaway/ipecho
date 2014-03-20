ipecho is a service to return a caller's IP and hostname.

This is useful in cases where your host is behind a NAT or in an environment like
OpenStack where it may not be able to determine locally what the external IP address is.

Endpoints:
/     - renders main page
/text - returns just IP address with no other formatting
/json - returns IP and hostname in JSON


