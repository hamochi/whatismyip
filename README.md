# whatismyip
Get your public Ip by looking it up on public  ip services.

[![GitHub License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hamochi/whatismyip/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/hamochi/whatismyip?status.svg)](https://godoc.org/github.com/hamochi/whatismyip)
[![Build Status](https://travis-ci.com/hamochi/whatismyip.svg?branch=master)](https://travis-ci.com/hamochi/whatismyip)

## Introduction
This package calls several public ip lookup services and when it matches ip from at least two sources it returns it, otherwise it returns an error.

## Install
```console
$ go get -u github.com/hamochi/whatismyip"
```

## Usage
```go
package main

import (
	"fmt"
	"log"
	"github.com/hamochi/whatismyip"
)

func main() {
	ip, err := whatismyip.Get()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ip.String())
}
```

The default list of IP lookup services are:
* https://checkip.amazonaws.com
* http://whatismyip.akamai.com
* https://api.ipify.org
* http://ifconfig.me/ip
* http://myexternalip.com/raw
* http://ipinfo.io/ip
* http://ipecho.net/plain
* http://icanhazip.com
* http://ifconfig.me/ip
* http://ident.me
* http://bot.whatismyipaddress.com
* http://wgetip.com
* http://ip.tyk.nu

If you want to pass your own list your can use whatismyip.GetWithCustomServices()

```go
package main

import (
	"fmt"
	"github.com/hamochi/whatismyip"
	"log"
)

func main() {
	ip, err := whatismyip.GetWithCustomServices([]string{
		"http://myexternalip.com/raw",
		"http://ipinfo.io/ip",
		"http://ipecho.net/plain",
		"http://icanhazip.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ip.String())
}
```

Please note that you need to pass at least two services into the list.

## Errors
You can cast the returned error to ApiErrors and inspect the returned error for each failed service call. Please take a look at the test file for more details.

## Credits
Based on https://github.com/chyeh/pubip
