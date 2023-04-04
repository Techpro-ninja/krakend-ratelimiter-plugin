// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// pluginName is the plugin name
var pluginName = "krakend-ratelimiter-plugin"

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer(pluginName)

type registerer string
var BucketCapacity int = 10
var BucketStock int = 10
var TokenRate int = 2


func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {

	// The config variable contains all the keys you have defined in the configuration
	// if the key doesn't exists or is not a map the plugin returns an error and the default handler
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	// The plugin will look for these configurations:
	trackerPath, _ := config["trackerpath"].(string)


	// Initiate the plugin background jobs
	go initPlugin()

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// If the requested path is not what we defined, continue.
		if req.URL.Path == ""{
			http.NotFound(w, req)
		} 
		
		if req.URL.Path != trackerPath{

			status, eMsg := Consume()

			if !status {
				http.Error(w, eMsg, http.StatusForbidden)
				return 
			}
			h.ServeHTTP(w, req)
			return
		}  

		if trackerPath != ""{
			fmt.Fprintf(w, Tracker())
		}
	}), nil
}

func main() {}



// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}

// Initializes the custom plugin's background jobs
func initPlugin(){
	time.Now()
	for true {
		time.Sleep(time.Minute)
		go MinuteUpdates()
	}
}

// Function for minute updates
func MinuteUpdates(){
	if BucketStock < BucketCapacity{
		if BucketStock + TokenRate >= BucketCapacity{
			BucketStock = BucketCapacity
		} else {
			BucketStock += TokenRate
		}
	}
}

func Consume()(bool, string) {
	if BucketStock > 0 {
		BucketStock--
		return true, ""
	}
	return false, "Token limit reached. Please wait for a minute"
}

func Tracker() string {
	return fmt.Sprintf("Total available tokens: %d/%d", BucketStock, BucketCapacity)
}
