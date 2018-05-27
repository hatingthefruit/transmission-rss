/*
 * Copyright (C) 2018 Aurélien Chabot <aurelien@chabot.fr>
 *
 * SPDX-License-Identifier: MIT
 */

package main

import (
	"fmt"
	"os"
)

import "github.com/jessevdk/go-flags"

type Options struct {
	Config string `short:"c" long:"conf" description:"Config file" default:"/etc/transmission-rss.conf"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	config := NewConfig(options.Config)

	client := NewTransmission(fmt.Sprintf("%s:%s",
		config.Server.Host, config.Server.Port))

	cache := NewCache()

	for _, feed := range config.Feeds {
		aggregator := NewAggregator(feed, cache)

		urls := aggregator.GetNewTorrentURL()
		for _, url := range urls {
			client.Add(url)
		}
	}
}
