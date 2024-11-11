package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter email addresses to check DNS records:")
	fmt.Println("domain, hasMax, hasSPF, sprRecord, hasSMARC, dmarcRecord etc")
	for scanner.Scan() {
		email := scanner.Text()
		// Extract domain from email
		domain := extractDomainFromEmail(email)
		if domain != "" {
			checkDomain(domain)
		} else {
			fmt.Println("Invalid email address:", email)
		}
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal("Error: Could not read from input:", err)
	}
}

func extractDomainFromEmail(email string) string {
	// Split the email by '@' and return the domain part
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1] // Return the domain part
	}
	return "" // Return empty if it's not a valid email
}

func checkDomain(domain string) {
	var hasMax, hasSMAR, hasSPF, hasDMARC bool
	var sprRecord, dmarcRecord string

	// Lookup MX records for the domain
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error looking up MX records for domain %s: %v\n", domain, err)
	} else if len(mxRecords) > 0 {
		hasMax = true
	}

	// Lookup TXT records to find SPF records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error looking up TXT records for domain %s: %v\n", domain, err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			sprRecord = record
			break
		}
	}

	// Lookup TXT records for DMARC
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records for domain %s: %v\n", domain, err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	// Print the results
	fmt.Printf("domain: %s, hasMax: %v, hasSMAR: %v, hasSPF: %v, hasDMARC: %v, sprRecord: %s, dmarcRecord: %s\n",
		domain, hasMax, hasSMAR, hasSPF, hasDMARC, sprRecord, dmarcRecord)
}
