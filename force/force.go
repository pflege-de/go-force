// Package force A Go package that provides bindings to the force.com REST API
//
// See http://www.salesforce.com/us/developer/docs/api_rest/
package force

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
)

const (
	DefaultAPIVersion = "v53.0"
)

type APIConfig func(*ForceApi)

func WithClient(c *http.Client) APIConfig {
	return func(f *ForceApi) {
		if c != nil {
			f.httpClient = c
		}
	}
}

var versionCheck = regexp.MustCompile(`v\d+\.\d+`)

func WithApiVersion(v string) APIConfig {
	return func(f *ForceApi) {
		f.apiVersion = v
	}
}

func WithInstance(instance string) APIConfig {
	return func(f *ForceApi) {
		f.instance = instance
	}
}

func WithAccessToken(clientId, accessToken, instanceUrl string) APIConfig {
	return func(f *ForceApi) {
		// NOTE: IMPORTANT: we're keeping support for (at least for the time being)
		// this to make the transition a bit easier since we have many dependents.
		f.instance = instanceUrl
		f.accessTokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	}
}

func WithTokenSource(source oauth2.TokenSource) APIConfig {
	return func(f *ForceApi) {
		f.accessTokenSource = source
	}
}

func NewClient(cfg ...APIConfig) (ForceApiInterface, error) {
	f := &ForceApi{
		apiResources:           make(map[string]string),
		apiSObjects:            make(map[string]*SObjectMetaData),
		apiSObjectDescriptions: make(map[string]*SObjectDescription),
		apiVersion:             "v53.0",
		httpClient:             http.DefaultClient,
	}

	for _, c := range cfg {
		c(f)
	}

	if !versionCheck.MatchString(f.apiVersion) {
		return nil, fmt.Errorf("invalid API version '%s' specified", f.apiVersion)
	}

	if f.accessTokenSource == nil {
		return nil, fmt.Errorf("missing access token source")
	}

	if _, err := f.accessTokenSource.Token(); err != nil {
		// NOTE: we attempt to immediately acquire a valid access token as a
		// sanity check (and so that we can fail early and loudly on error).
		// ... this will also explicitly "prime" the token source "cache"...
		return nil, fmt.Errorf("failed to acquire access token: %w", err)
	}

	// Init ForceApi Resources
	err := f.getApiResources()
	if err != nil {
		return nil, err
	}
	err = f.getApiSObjects()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func CreateWithAccessToken(version, clientId, accessToken, instanceUrl string, httpClient *http.Client) (ForceApiInterface, error) {
	return NewClient(
		WithAccessToken(clientId, accessToken, instanceUrl),
		WithClient(httpClient),
	)
}

func CreateWithTokenSource(source oauth2.TokenSource, instance string, client *http.Client) (ForceApiInterface, error) {
	return NewClient(WithTokenSource(source), WithInstance(instance), WithClient(client))
}

// Used when running tests.
func createTest() ForceApiInterface {
	const (
		testVersion     = "v53.0"
		testEnvironment = "production"
		testLoginUri    = "https://login.salesforce.com/services/oauth2/token"

		testClientId      = "3MVG9A2kN3Bn17hs8MIaQx1voVGy662rXlC37svtmLmt6wO_iik8Hnk3DlcYjKRvzVNGWLFlGRH1ryHwS217h"
		testClientSecret  = "4165772184959202901"
		testUserName      = "go-force@jalali.net"
		testPassword      = "golangrocks3"
		testSecurityToken = "kAlicVmti9nWRKRiWG3Zvqtte" //nolint:gosec Just for testing purpose
	)

	token, err := (&oauth2.Config{
		ClientID:     testClientId,
		ClientSecret: testClientSecret,

		Endpoint: oauth2.Endpoint{
			TokenURL: testLoginUri,
		},
	}).PasswordCredentialsToken(context.Background(), testUserName, strings.Join([]string{testPassword, testSecurityToken}, ""))

	if err != nil {
		fmt.Printf("Unable to create ForceApi for test: %v", err)
		os.Exit(1)
	}

	instance, _ := token.Extra("instance_url").(string)

	forceApi, err := CreateWithAccessToken(testVersion, testClientId, token.AccessToken, instance, http.DefaultClient)
	if err != nil {
		fmt.Printf("Unable to create ForceApi for test: %v", err)
		os.Exit(1)
	}

	return forceApi
}

type ApiLogger interface {
	Printf(format string, v ...interface{})
}

// TraceOn turns on logging for this ForceApi. After this is called, all
// requests, responses, and raw response bodies will be sent to the logger.
// If prefix is a non-empty string, it will be written to the front of all
// logged strings, which can aid in filtering log lines.
//
// Use TraceOn if you want to spy on the ForceApi requests and responses.
//
// Note that the base log.Logger type satisfies ApiLogger, but adapters
// can easily be written for other logging packages (e.g., the
// golang-sanctioned glog framework).
func (forceApi *ForceApi) TraceOn(prefix string, logger ApiLogger) {
	forceApi.logger = logger
	if prefix == "" {
		forceApi.logPrefix = prefix
	} else {
		forceApi.logPrefix = fmt.Sprintf("%s ", prefix)
	}
}

// TraceOff turns off tracing. It is idempotent.
func (forceApi *ForceApi) TraceOff() {
	forceApi.logger = nil
	forceApi.logPrefix = ""
}

func (forceApi *ForceApi) trace(name string, value interface{}, format string) {
	if forceApi.logger != nil {
		logMsg := "%s%s " + format + "\n"
		forceApi.logger.Printf(logMsg, forceApi.logPrefix, name, value)
	}
}
