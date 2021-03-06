# Keychain Storage for N26 API Client

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/n26keychain)](https://github.com/nhatthm/n26keychain/releases/latest)
[![Build Status](https://github.com/nhatthm/n26keychain/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/n26keychain/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/n26keychain/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/n26keychain)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/httpmock)](https://goreportcard.com/report/github.com/nhatthm/httpmock)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/n26keychain)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

`n26keychain` uses [system keyring](https://github.com/zalando/go-keyring#go-keyring-library) as a storage for persisting/getting credentials and token. It supports OS X, Linux 
(dbus) and Windows. 

## Prerequisites

- `Go >= 1.17`

## Install

```bash
go get github.com/nhatthm/n26keychain
```

## Usage

### `n26api.CredentialsProvider`

**Examples**

Build `n26api.Client`:

```go
package mypackage

import (
	"github.com/google/uuid"
	"github.com/nhatthm/n26api"
	"github.com/nhatthm/n26keychain/credentials"
)

func buildClient() (*n26api.Client, error) {
	deviceID := uuid.New()

	c := n26api.NewClient(
		n26api.WithDeviceID(deviceID),
		credentials.WithCredentialsProvider(),
	)

	return c, nil
}
```

Persist credentials in system keyring:

```go
package mypackage

import (
	"github.com/google/uuid"
	"github.com/nhatthm/n26keychain/credentials"
)

func persist(deviceID uuid.UUID, username, password string) error {
	c := credentials.New(deviceID)
	
	return c.Update(username, password)
}
```

### `auth.TokenStorage`

```go
package mypackage

import (
	"github.com/nhatthm/n26api"
	"github.com/nhatthm/n26keychain/token"
)

func buildClient() *n26api.Client {
	return n26api.NewClient(
		token.WithTokenStorage(),
	)
}
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
