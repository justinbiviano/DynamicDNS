package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env file relying on system env.")
	}

	recordResponse, err := UpdateDNSRecord(os.Getenv("IP_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}

	neatJSON, err := json.MarshalIndent(recordResponse, "", "	")
	if err != nil {
		log.Printf("Failed to neaten JSON: %w", err)
	}

	fmt.Print(string(neatJSON))
}

func UpdateDNSRecord(NewIPAddress string) (*dns.RecordResponse, error) {
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
				Content: cloudflare.F(NewIPAddress),
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to update DNS record: %w", err)
	}

	return recordResponse, nil
}
