package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
)

func main() {
	client := cloudflare.NewClient(
		option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")),
	)
	recordResponse, err := client.DNS.Records.Update(
		context.TODO(),
		os.Getenv("DNS_RECORD_ID"),
		dns.RecordUpdateParams{
			ZoneID: cloudflare.F(os.Getenv("ZONE_ID")),
			Body: dns.ARecordParam{
				Name:    cloudflare.F(os.Getenv("DOMAIN_NAME")),
				TTL:     cloudflare.F(dns.TTL1),
				Type:    cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F("1.42.114.177"),
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", recordResponse)
}
