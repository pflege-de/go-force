# go-force


[![Go Report Card](https://goreportcard.com/badge/github.com/pflege-de/go-force)](https://goreportcard.com/report/github.com/pflege-de/go-force)
[![build status](https://github.com/pflege-de/go-force/workflows/build/badge.svg)](https://github.com/pflege-de/go-force/actions?query=workflow%3Abuild)
[![Go version](https://img.shields.io/github/go-mod/go-version/pflege-de/go-force)](https://github.com/pflege-de/go-force/blob/master/go.mod)
[![Current Release](https://img.shields.io/github/v/release/pflege-de/go-force.svg)](https://github.com/pflege-de/go-force/releases)
[![godoc](https://godoc.org/github.com/pflege-de/go-force?status.svg)](https://godoc.org/github.com/pflege-de/go-force)
<!-- [![Go Coverage](https://github.com/pflege-de/go-force/wiki/coverage.svg)](https://raw.githack.com/wiki/pflege-de/go-force/coverage.html) -->
[![License](https://img.shields.io/github/license/pflege-de/go-force)](https://github.com/pflege-de/go-force/blob/master/LICENSE)

[Golang](http://golang.org/) API wrapper for [Force.com](http://www.force.com/), [Salesforce.com](http://www.salesforce.com/)


This repo is based on https://github.com/nimajalali/go-force which seems to be dormant with the last commit 4 years ago at this time (01/2024).

## Installation

	go get github.com/pflege-de/go-force/force

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/pflege-de/go-force/force"
	"github.com/pflege-de/go-force/sobjects"
)

type SomeCustomSObject struct {
	sobjects.BaseSObject
	
	Active    bool   `force:"Active__c"`
	AccountId string `force:"Account__c"`
}

func (t *SomeCustomSObject) ApiName() string {
	return "SomeCustomObject__c"
}

type SomeCustomSObjectQueryResponse struct {
	sobjects.BaseQuery

	Records []*SomeCustomSObject `force:"records"`
}

func main() {
	// Init the force
	forceApi, err := force.Create(
		"YOUR-API-VERSION",
		"YOUR-CLIENT-ID",
		"YOUR-CLIENT-SECRET",
		"YOUR-USERNAME",
		"YOUR-PASSWORD",
		"YOUR-SECURITY-TOKEN",
		"YOUR-ENVIRONMENT",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get somCustomSObject by ID
	someCustomSObject := &SomeCustomSObject{}
	err = forceApi.GetSObject("Your-Object-ID", nil, someCustomSObject)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", someCustomSObject)

	// Query
	someCustomSObjects := &SomeCustomSObjectQueryResponse{}
	err = forceApi.Query("SELECT Id FROM SomeCustomSObject__c LIMIT 10", someCustomSObjects)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v", someCustomSObjects)
}
```

## Documentation 

* [Package Reference](http://godoc.org/github.com/pflege-de/go-force/force)
* [Force.com API Reference](http://www.salesforce.com/us/developer/docs/api_rest/)
