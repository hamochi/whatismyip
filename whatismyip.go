package whatismyip

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// Timeout is the default time before each request is cancelled
var Timeout = time.Second * 2

var defaultIpServices = []string{
	"https://checkip.amazonaws.com",
	"http://whatismyip.akamai.com",
	"https://api.ipify.org",
	"http://ifconfig.me/ip",
	"http://myexternalip.com/raw",
	"http://ipinfo.io/ip",
	"http://ipecho.net/plain",
	"http://icanhazip.com",
	"http://ifconfig.me/ip",
	"http://ident.me",
	"http://bot.whatismyipaddress.com",
	"http://wgetip.com",
	"http://ip.tyk.nu",
}

type apiResult struct {
	endPoint string
	ip       net.IP
	err      error
}

type ApiErrors []apiResult

func (a ApiErrors) Is(target error) bool {
	_, ok := target.(ApiErrors)
	if !ok {
		return false
	}
	return true
}

func (a ApiErrors) Error() string {
	s := strings.Builder{}
	for _, result := range a {
		if result.err != nil {
			s.WriteString(fmt.Sprintf("\nEndpoint: %s, Returned IP: nil, error: %s", result.endPoint, result.err.Error()))
		} else {
			s.WriteString(fmt.Sprintf("\nEndpoint: %s, Returned IP: %s, error: nil", result.endPoint, result.ip.String()))
		}
	}
	return s.String()
}

// Get returns the Ip if it get 2 matched IPs from the default ip Lookup services
//
// Usage:
//
//package main
//
//import (
//	"fmt"
//	"log"
//	"github.com/hamochi/whatismyip"
//)
//
//func main() {
//	ip, err := whatismyip.Get()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(ip.String())
//}
func Get() (net.IP, error) {
	return GetWithCustomServices(defaultIpServices)
}

// GetWithCustomServices returns the Ip if it get 2 matched IPs from the provided ip Lookup services
// Usage:
//
//package main
//
//import (
//	"fmt"
//	"github.com/hamochi/whatismyip"
//	"log"
//)
//
//func main() {
//	ip, err := whatismyip.GetWithCustomServices([]string{
//		"http://myexternalip.com/raw",
//		"http://ipinfo.io/ip",
//		"http://ipecho.net/plain",
//		"http://icanhazip.com",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(ip.String())
//}
func GetWithCustomServices(services []string) (net.IP, error) {
	resultCh := make(chan apiResult)
	ctx, cancel := context.WithCancel(context.Background())
	for _, endpoint := range services {
		go func(e string) {
			var result apiResult
			result.endPoint = endpoint

			ip, err := call(ctx, e)
			if err != nil {
				result.err = err
			}
			result.ip = ip
			resultCh <- result
		}(endpoint)
	}

	count := make(map[string]bool)
	allResults := ApiErrors{}
	for i := 0; i < len(services); i++ {
		select {
		case result := <-resultCh:
			allResults = append(allResults, result)
			if result.ip != nil {
				_, ok := count[result.ip.String()]
				if ok {
					cancel()
					return result.ip, nil
				}
				count[result.ip.String()] = true
			}
		}
	}

	return nil, fmt.Errorf("could not get two matched IPs: returned %w", allResults)
}

func call(ctx context.Context, endpoint string) (net.IP, error) {
	client := http.Client{
		Timeout: Timeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ip, err := parseIp(body)
	if err != nil {
		return nil, err
	}

	return ip, err
}

func parseIp(body []byte) (net.IP, error) {
	ip := net.ParseIP(strings.TrimSpace(string(body)))
	if ip == nil {
		return nil, errors.New("could not parse ip")
	}
	return ip, nil
}
