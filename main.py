import os
from cloudflare import Cloudflare
from dotenv import load_dotenv

load_dotenv()

client = Cloudflare(
    api_token=os.environ.get("CLOUDFLARE_API_TOKEN"),  # This is the default and can be omitted
)
record_response = client.dns.records.update(
    dns_record_id=os.environ.get("DNS_RECORD_ID"),
    zone_id=os.environ.get("ZONE_ID"),
    name="jbivi.me",
    ttl=3600,
    type="A",
    content="1.42.114.177",
    proxied=False,
)
print(record_response)