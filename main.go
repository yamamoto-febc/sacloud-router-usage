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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/sacloud-router-usage/usage"
	"github.com/sacloud/sacloud-router-usage/version"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := &commandOpts{
		Option: &usage.Option{},
	}
	if err := usage.ParseOption(opts); err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}
	if opts.Version {
		usage.PrintVersion()
		return usage.ExitOk
	}

	client, err := routerClient()
	if err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}

	resources, err := fetchResources(client, opts)
	if err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}

	if err := usage.OutputMetrics(os.Stdout, resources.Metrics(), opts.Query); err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}
	return usage.ExitOk
}

type commandOpts struct {
	*usage.Option
	Item string `long:"item" description:"Item name" required:"true" choice:"in" choice:"out" default:"in"`
}

func routerClient() (iaasRouterAPI, error) {
	options := api.OptionsFromEnv()
	if options.AccessToken == "" {
		return nil, fmt.Errorf("environment variable %q is required", "SAKURACLOUD_ACCESS_TOKEN")
	}
	if options.AccessTokenSecret == "" {
		return nil, fmt.Errorf("environment variable %q is required", "SAKURACLOUD_ACCESS_TOKEN_SECRET")
	}

	if options.UserAgent == "" {
		options.UserAgent = fmt.Sprintf(
			"sacloud/sacloud-router-uage/v%s (%s/%s; +https://github.com/sacloud/sacloud-router-uage) %s",
			version.Version,
			runtime.GOOS,
			runtime.GOARCH,
			iaas.DefaultUserAgent,
		)
	}

	caller := api.NewCallerWithOptions(options)
	return iaas.NewInternetOp(caller), nil
}

type iaasRouterAPI interface {
	Find(ctx context.Context, zone string, conditions *iaas.FindCondition) (*iaas.InternetFindResult, error)
	MonitorRouter(ctx context.Context, zone string, id types.ID, condition *iaas.MonitorCondition) (*iaas.RouterActivity, error)
}

func fetchResources(client iaasRouterAPI, opts *commandOpts) (*usage.Resources, error) {
	rs := &usage.Resources{Label: "routers", Option: opts.Option}
	for _, prefix := range opts.Prefix {
		for _, zone := range opts.Zones {
			condition := &iaas.FindCondition{
				Filter: map[search.FilterKey]interface{}{},
			}
			condition.Filter[search.Key("Name")] = search.PartialMatch(prefix)
			result, err := client.Find(
				context.Background(),
				zone,
				condition,
			)
			if err != nil {
				return nil, err
			}
			for _, r := range result.Internet {
				if !strings.HasPrefix(r.Name, prefix) {
					continue
				}
				monitors, err := fetchRouterActivities(client, zone, r.ID, opts)
				if err != nil {
					return nil, err
				}
				rs.Resources = append(rs.Resources, &usage.Resource{
					ID:             r.ID,
					Name:           r.Name,
					Zone:           zone,
					Monitors:       monitors,
					Label:          "traffic",
					AdditionalInfo: nil,
				})
			}
		}
	}
	return rs, nil
}

func fetchRouterActivities(client iaasRouterAPI, zone string, id types.ID, opts *commandOpts) ([]usage.MonitorValue, error) {
	b, _ := time.ParseDuration(fmt.Sprintf("-%dm", (opts.Time+3)*5))
	condition := &iaas.MonitorCondition{
		Start: time.Now().Add(b),
		End:   time.Now(),
	}
	activity, err := client.MonitorRouter(context.Background(), zone, id, condition)
	if err != nil {
		return nil, err
	}
	usages := activity.GetValues()
	if len(usages) > int(opts.Time) {
		usages = usages[len(usages)-int(opts.Time):]
	}

	var results []usage.MonitorValue
	for _, u := range usages {
		v := u.Out
		if opts.Item == "in" {
			v = u.In
		}
		if v > 0 {
			v = v / 1000 / 1000 // 単位変換: bps->Mbps
		}
		results = append(results, usage.MonitorValue{Time: u.Time, Value: v})
	}
	return results, nil
}
