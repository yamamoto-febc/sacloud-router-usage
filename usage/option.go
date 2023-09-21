// Copyright 2022-2023 The sacloud/sacloud-router-usage Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package usage

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

type Option struct {
	Time          uint     `long:"time" description:"Get average traffic for a specified amount of time" default:"3"`
	Prefix        []string `long:"prefix" description:"prefix for router names. prefix accepts more than one." required:"true"`
	Zones         []string `long:"zone" description:"zone name" required:"true"`
	PercentileSet string   `long:"percentile-set" default:"99,95,90,75" description:"percentiles to dispaly"`
	Version       bool     `short:"v" long:"version" description:"Show version"`
	Query         string   `long:"query" description:"jq style query to result and display"`
	EnvFrom       string   `long:"env-from" description:"load environment values from this file"`
	percentiles   []percentile
}

type optionProvider interface {
	option() *Option
}

func (o *Option) option() *Option {
	return o
}

type percentile struct {
	str   string
	float float64
}

func ParseOption(o interface{}) error {
	psr := flags.NewParser(o, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()

	v, ok := o.(optionProvider)
	if !ok {
		return nil
	}
	opts := v.option()

	if opts.Version {
		return nil
	}

	if err != nil {
		return err
	}

	if opts.Time < 1 {
		opts.Time = 1
	}

	if opts.EnvFrom != "" {
		if err := godotenv.Load(opts.EnvFrom); err != nil {
			return err
		}
	}

	m := make(map[string]struct{})
	for _, z := range opts.Zones {
		if _, ok := m[z]; ok {
			return fmt.Errorf("zone %q is duplicated", z)
		}
		m[z] = struct{}{}
	}

	var percentiles []percentile
	percentileStrings := strings.Split(opts.PercentileSet, ",")
	for _, s := range percentileStrings {
		if s == "" {
			continue
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("could not parse --percentile-set: %v", err)
		}
		f /= 100
		percentiles = append(percentiles, percentile{s, f})
	}
	opts.percentiles = percentiles

	return nil
}
