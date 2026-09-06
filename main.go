package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No env file found.")
	}

	interval := 5 * time.Minute
	v := os.Getenv("CHECK_INTERVAL")
	if v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			interval = d
		}
	}

	for {
		runCheck()
		time.Sleep(interval)
	}
}

func runCheck() {
	currentIP, publicIP, err := GetIPs()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Current IP", currentIP)
	fmt.Println("Public IP", publicIP)

	if currentIP != publicIP {
		fmt.Printf("Record update required! Updating to: %s\n", publicIP)
		recordResponse, err := UpdateDNSRecord(publicIP)
		if err != nil {
			log.Printf("Error during record update: %v", err)
			return
		}

		neatJSON, err := json.MarshalIndent(recordResponse, "", "	")
		if err != nil {
			log.Printf("Failed to neaten JSON: %v", err)
			return
		}
		fmt.Print(string(neatJSON))
	} else {
		fmt.Println("IP's Match")
	}
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
		return nil, fmt.Errorf("Failed to update DNS record: %w\n", err)
	}

	return recordResponse, nil
}

func GetIPs() (string, string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	ipResponse, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", "", fmt.Errorf("Error fetching IP: %v\n", err)
	}
	defer ipResponse.Body.Close()
	body, err := io.ReadAll(ipResponse.Body)
	if err != nil {
		return "", "", fmt.Errorf("Error reading response: %w\n", err)
	}
	publicIP := strings.TrimSpace(string(body))

	resolveIP, err := net.LookupIP(os.Getenv("DOMAIN_NAME"))
	if err != nil {
		return "", "", fmt.Errorf("Failed resolving Domain: %w", err)
	}
	if len(resolveIP) == 0 {
		return "", "", fmt.Errorf("No records found for domain")
	}

	currentIP := resolveIP[0].String()
	return currentIP, publicIP, nil
}
