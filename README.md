# TinyMacDns

DNS & rDNS server utilizing an arp table stored in redis.  
All settings are configured through environment variables.


## Settings

### Required

* TMD_REDIS_ADDRESS =  redis server address

### Overview

| Name                | Type    | Default | Description |
|---------------------|----------|-------|--------------|
| TMD_VERBOSE         | bool     | False | Enables extensive logging |
| TMD_REDIS_ADDRESS   | string   |       | Redis connection ip:port |
| TMD_REDIS_USERNAME  | string   |       | Redis username |
| TMD_REDIS_PASSWORD  | string   |       | Redis password |
| TMD_REDIS_DATABASE  | int      | 0     | Redis database |
| TMD_REDIS_ATTEMPTS  | int      | 3     | Redis connection attempts |
| TMD_REDIS_COOLDOWN  | duration | 1s    | Coooldown between redis connection attempts |
| TMD_HOSTS_hostname_MAC_n | string   |       | Hostname/MAC mapping (n = number of MAC mapping order) |


## Docker stack example

```yaml
version: "3.8"

services:
  redis:
    image: redis:alpine
    networks:
      - int
    ports:
      - 127.0.0.1:6379:6379
  collector:
    image: ghcr.io/kwitsch/arprediscollector
    environment:
      - ARC_REDIS_ADDRESS=127.0.0.1:6379
      - ARC_ARP_SUBNET_1=192.168.0.0/24
    cap_add:
      - CAP_NET_ADMIN
    networks:
      - host
  server:
    image: ghcr.io/kwitsch/tinymacdns
    environment:
      - TMD_VERBOSE=True
      - TMD_REDIS_ADDRESS=redis:6379
      - TMD_HOSTS_ExamplePC1_MAC_1=1a:2b:3c:4d:2a:4b
      - TMD_HOSTS_ExampleLaptop_MAC_1=3c:4d:2a:4b:1a:2b
    networks:
      - int
    ports:
      - 53:53/tcp
      - 53:53/udp

networks:
  int:
    internal: true
  host:
    name: host
    external: true
```