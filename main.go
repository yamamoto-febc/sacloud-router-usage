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
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sacloud/go-otelsetup"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/sacloud-router-usage/version"
	usage "github.com/sacloud/sacloud-usage-lib"
	"go.opentelemetry.io/otel"
)

func main() {
	// test
	os.Exit(_main())
}

const appName = "github.com/sacloud/sacloud-router-usage"

func _main() int {
	// initialize OTel SDK
	otelShutdown, err := otelsetup.InitWithOptions(context.Background(), otelsetup.Options{
		ServiceName:      "sacloud-router-usage",
		ServiceVersion:   version.Version,
		ServiceNamespace: "sacloud",
	})
	if err != nil {
		log.Println("Error in initializing OTel SDK: " + err.Error())
		return usage.ExitUnknown
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
		if err != nil {
			log.Println("Error in initializing OTel SDK: " + err.Error())
		}
	}()

	// init root span
	ctx, span := otel.Tracer(appName).Start(otelsetup.ContextForTrace(context.Background()), "main")
	defer span.End()

	opts := &commandOpts{
		Option: &usage.Option{},
	}
	if err := usage.ParseOption(opts); err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}
	if opts.Option.Version {
		usage.PrintVersion(version.Version)
		return usage.ExitOk
	}

	caller, err := usage.SacloudAPICaller("sacloud-router-usage", version.Version)
	if err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}

	resources, err := fetchResources(ctx, iaas.NewInternetOp(caller), opts)
	if err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}

	if err := usage.OutputMetrics(os.Stdout, resources.Metrics(), opts.Option.Query); err != nil {
		log.Println(err)
		return usage.ExitUnknown
	}
	return usage.ExitOk
}

type commandOpts struct {
	*usage.Option
	Item string `long:"item" description:"Item name" required:"true" choice:"in" choice:"out" default:"in"`
}

type iaasRouterAPI interface {
	Find(ctx context.Context, zone string, conditions *iaas.FindCondition) (*iaas.InternetFindResult, error)
	MonitorRouter(ctx context.Context, zone string, id types.ID, condition *iaas.MonitorCondition) (*iaas.RouterActivity, error)
}

func fetchResources(ctx context.Context, client iaasRouterAPI, opts *commandOpts) (*usage.Resources, error) {
	ctx, span := otel.Tracer(appName).Start(ctx, "fetchResources")
	defer span.End()

	rs := &usage.Resources{Label: "routers", Option: opts.Option}
	for _, prefix := range opts.Option.Prefix {
		for _, zone := range opts.Option.Zones {
			condition := &iaas.FindCondition{
				Filter: map[search.FilterKey]interface{}{},
			}
			condition.Filter[search.Key("Name")] = search.PartialMatch(prefix)
			result, err := client.Find(
				ctx,
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
				monitors, err := fetchRouterActivities(ctx, client, zone, r.ID, opts)
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

func fetchRouterActivities(ctx context.Context, client iaasRouterAPI, zone string, id types.ID, opts *commandOpts) ([]usage.MonitorValue, error) {
	ctx, span := otel.Tracer(appName).Start(ctx, "fetchRouterActivities")
	defer span.End()

	b, _ := time.ParseDuration(fmt.Sprintf("-%dm", (opts.Option.Time+3)*5))
	condition := &iaas.MonitorCondition{
		Start: time.Now().Add(b),
		End:   time.Now(),
	}
	activity, err := client.MonitorRouter(ctx, zone, id, condition)
	if err != nil {
		return nil, err
	}
	usages := activity.GetValues()
	if len(usages) > int(opts.Option.Time) { //nolint:gosec
		usages = usages[len(usages)-int(opts.Option.Time):] //nolint:gosec
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
