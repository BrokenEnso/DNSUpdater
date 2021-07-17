package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

type Configuration struct {
	ApiToken string
	Domain   string
}

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "Specify the configuration file location. Default is config.json")
	flag.Parse()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file (%v) not found: ", configFile))
	}

	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	if err := decoder.Decode(&configuration); err != nil {
		exitError(err)
	}

	var apiToken = configuration.ApiToken
	var domain = configuration.Domain
	var domainParts = strings.Split(domain, ".")
	var zoneName = strings.Join(domainParts[len(domainParts)-2:], ".")

	/**/

	configuredIP, err := lookupDominDnsIP(domain)
	if err != nil {
		exitError(err)
	}

	currentIP, err := lookupIPUsingAkamai()
	if err != nil {
		exitError(err)
	}

	if currentIP != configuredIP {
		fmt.Println("IP has changed and needs an updated configuration.")
		updateCloudFlare(apiToken, domain, zoneName, currentIP)

	} else {
		fmt.Println("IP hasn't changed. Nothing to do.")
	}

}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func lookupDominDnsIP(domain string) (string, error) {
	addr, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	var configuredIP string
	if len(addr) == 1 {
		configuredIP = addr[0].String()
	}

	return configuredIP, nil
}

func lookupIPUsingAkamai() (string, error) {
	resp, err := http.Get("http://whatismyip.akamai.com")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var currentIP = string(body)

	return currentIP, nil
}

func updateCloudFlare(apiToken string, domain string, zone string, ipAddress string) {
	/**/
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		exitError(err)
	}

	zoneID, err := api.ZoneIDByName(zone)
	if err != nil {
		exitError(err)
	}

	// Fetch records of any type with name "foo.example.com"
	// The name must be fully-qualified
	foo := cloudflare.DNSRecord{Name: domain, Type: "A"}

	recs, err := api.DNSRecords(context.Background(), zoneID, foo)
	if err != nil {
		exitError(err)
	}

	var rec cloudflare.DNSRecord

	if len(recs) == 1 {
		rec = recs[0]
	} else {
		fmt.Println("Cloudflare: No DNS  record found")
	}

	fmt.Printf("OLD: %s: %s\n", rec.Name, rec.Content)

	rec.Content = ipAddress

	api.UpdateDNSRecord(context.Background(), zoneID, rec.ID, rec)

	fmt.Printf("NEW: %s: %s\n", rec.Name, rec.Content)

	/**/
}
