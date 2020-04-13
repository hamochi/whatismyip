package whatismyip

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetWithCustomServices(t *testing.T) {
	getHandler := func(resp string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, resp)
		}
	}

	var tests = []struct {
		name          string
		returnHandler []http.HandlerFunc
		expectFinalIp net.IP
		expectedError error
	}{
		{"allGood", []http.HandlerFunc{getHandler("213.76.193.55"), getHandler("213.76.193.55"), getHandler("213.76.193.55")}, net.ParseIP("213.76.193.55"), nil},
		{"notAllReturnedTheSame", []http.HandlerFunc{getHandler("10.233.193.34"), getHandler("213.233.193.34"), getHandler("112.76.193.55")}, nil, ApiErrors{}},
		{"badIP", []http.HandlerFunc{getHandler("blabla"), getHandler("929384"), getHandler("112.76.193.55")}, nil, ApiErrors{}},
		{"twoOfThreeGood", []http.HandlerFunc{getHandler("404"), getHandler("213.76.193.55"), getHandler("213.76.193.55")}, net.ParseIP("213.76.193.55"), nil},
		{"twoOfManyGood", []http.HandlerFunc{getHandler("404"), getHandler("404"), getHandler("404"), getHandler("404"), getHandler("213.76.193.55"), getHandler("213.76.193.55")}, net.ParseIP("213.76.193.55"), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var endpoints []string
			for _, h := range tt.returnHandler {
				ts := httptest.NewServer(h)
				defer ts.Close()
				endpoints = append(endpoints, ts.URL)
			}
			ip, err := GetWithCustomServices(endpoints)
			if err != nil {
				if tt.expectedError == nil {
					t.Fatalf("exptected no errors, got: %s", err.Error())
				}

				if errors.Is(err, tt.expectedError) == false {
					t.Fatalf("miss matched errors, got: %v, want %v", err, tt.expectedError)
				}
			}

			if err == nil && tt.expectedError != nil {
				t.Fatal("expected err, got none")
			}

			if ip != nil && tt.expectFinalIp != nil && ip.String() != tt.expectFinalIp.String() {
				t.Fatalf("unexpected returned ip, expected %s got %s", tt.expectFinalIp.String(), ip.String())
			}
		})
	}

}

func TestGet(t *testing.T) {
	_, err := Get()
	if err != nil {
		t.Fatal("unexpected err", err)
	}
}

func TestTimeOut(t *testing.T) {
	Timeout = time.Nanosecond

	_, err := Get()
	if err == nil {
		t.Fatal("expected error got non")
	}

	var apiResults ApiErrors
	ok := errors.As(err, &apiResults)
	if !ok {
		t.Fatal("unable to cast error to ApiErrors")
	}

	for _, res := range apiResults {
		var netErr net.Error
		if res.Err != nil {
			ok := errors.As(res.Err, &netErr)
			if !ok {
				t.Fatalf("unexpected error type, want net.Error, got something else %s", res.Err.Error())
			}
			if netErr.Timeout() == false {
				t.Fatalf("unexpected net.Error, was looking for Timeout() to be true, but it's not")
			}
		}

	}

}

func TestDefaultServices(t *testing.T) {
	Timeout = time.Second * 2

	for _, s := range defaultIpServices {
		t.Run(s, func(t *testing.T) {
			ctx, _ := context.WithCancel(context.Background())
			_, err := call(ctx, s)
			if err != nil {
				t.Errorf("found err in endpoint: %s ", err)
			}
		})
	}
}
