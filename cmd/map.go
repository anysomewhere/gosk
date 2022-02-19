/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/munnik/gosk/config"
	"github.com/munnik/gosk/logger"
	"github.com/munnik/gosk/mapper"
	"github.com/munnik/gosk/nanomsg"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	mapCmd = &cobra.Command{
		Use:   "map",
		Short: "Map raw data to meaningfull data",
		Long:  `Map raw data to meaningfull data based on the SignalK specification`,
		Run:   doMap,
	}
	mapSubscribeURI string
	mapPublishURI   string
)

func init() {
	rootCmd.AddCommand(mapCmd)
	mapCmd.Flags().StringVarP(&mapPublishURI, "publishURI", "u", "", "Nanomsg URI, the URI is used to publish the collected data on. It listens for connections.")
	mapCmd.MarkFlagRequired("publishURI")
	mapCmd.Flags().StringVarP(&mapSubscribeURI, "subscribeURI", "s", "", "Nanomsg URI, the URI is used to listen for subscribed data.")
	mapCmd.MarkFlagRequired("subscribeURI")
}

func doMap(cmd *cobra.Command, args []string) {
	protocol, err := getProtocol(cfgFile)
	if err != nil {
		logger.GetLogger().Fatal(
			"Could not determine the protocol, make sure the protocol is in the path of the config file",
			zap.String("Config file", cfgFile),
			zap.String("Error", err.Error()),
		)
	}

	subscriber, err := nanomsg.NewSub(mapSubscribeURI, []byte{})
	if err != nil {
		logger.GetLogger().Fatal(
			"Could not subscribe",
			zap.String("URI", mapSubscribeURI),
			zap.String("Error", err.Error()),
		)
	}
	publisher := nanomsg.NewPub(mapPublishURI)

	c := config.NewMapperConfig(cfgFile)
	var m mapper.Mapper
	switch protocol {
	case config.CsvType:
		cmc := config.NewCsvMappingConfig(cfgFile)
		m, err = mapper.NewCsvMapper(c, cmc)
	case config.ModbusType:
		rmc := config.NewRegisterMappingsConfig(cfgFile)
		m, err = mapper.NewModbusMapper(c, rmc)
	case config.NMEA0183Type:
		m, err = mapper.NewNmea0183Mapper(c)
	default:
		logger.GetLogger().Fatal(
			"Not a supported protocol",
			zap.String("Protocol", protocol),
		)
		return
	}
	if err != nil {
		logger.GetLogger().Fatal(
			"Error while creating the collector",
			zap.String("Config file", cfgFile),
			zap.String("Error", err.Error()),
		)
	}
	m.Map(subscriber, publisher)
}
