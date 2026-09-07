# Dynamic DNS for Cloudflare Services
## How to Run - Docker Compose

Create a compose file in a new directory and copy in the following.
``` yaml
services:
  dynamicDNS:
    image: ghcr.io/justinbiviano/dynamicdns
    container_name: dynamicDNS
    restart: unless-stopped
    env_file:
      - .env
```

Then create an .env file and use the following template.
``` yaml
CLOUDFLARE_API_TOKEN="abc123"
DNS_RECORD_ID="READ BELOW"
ZONE_ID="FIND IN DOMAIN OVERVIEW NEAR THE BOTTOM"
DOMAIN_NAME="example.com"
CHECK_INTERVAL=5m
```

Cloudflare DNS record id is tricky to find but once you have the zone ID and api token you can run. Then look for the id before the domain name you want to find.
``` bash
curl -X GET "https://api.cloudflare.com/client/v4/zones/[ZONE_ID]/dns_records" \
  -H "Authorization: Bearer [YOUR_API_TOKEN]" \
  -H "Content-Type: application/json" | jq
```
