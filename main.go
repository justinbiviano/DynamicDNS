package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

	currentIP, publicIP, err := GetIPs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current IP", currentIP)
	fmt.Println("Public IP", publicIP)
	if currentIP != publicIP {
		fmt.Print("FALSE")
	}

	// recordResponse, err := UpdateDNSRecord(os.Getenv("IP_ADDRESS"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// neatJSON, err := json.MarshalIndent(recordResponse, "", "	")
	// if err != nil {
	// 	log.Printf("Failed to neaten JSON: %w", err)
	// }

	// fmt.Print(string(neatJSON))
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

func GetIPs() (string, string, error) {
	ipResponse, err := http.Get("https://api.ipify.org")
	if err != nil {
		defer ipResponse.Body.Close()
		return "", "", fmt.Errorf("Error fetching IP: %v\n", err)
	}
	defer ipResponse.Body.Close()

	body, err := io.ReadAll(ipResponse.Body)
	if err != nil {
		return "", "", fmt.Errorf("Error reading response: %w\n", err)
	}

	publicIP := string(body)

	resolveIP, err := net.LookupIP(os.Getenv("DOMAIN_NAME"))
	currentIP := resolveIP[0].String()
	return currentIP, publicIP, nil
}
