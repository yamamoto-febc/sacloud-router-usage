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
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/sacloud/iaas-api-go/types"
)

type monitorValue struct {
	Time  time.Time
	Value float64
}

type iaasResource struct {
	ID   types.ID
	Name string
	Zone string

	Monitors []monitorValue
	Label    string

	AdditionalInfo map[string]interface{}
}

func (r *iaasResource) toMetrics() map[string]interface{} {
	sum := float64(0)
	monitors := make([]interface{}, 0)
	for _, p := range r.Monitors {
		m := map[string]interface{}{
			r.Label: p.Value,
			"time":  p.Time.String(),
		}
		monitors = append(monitors, m)
		sum += p.Value
		log.Printf("%s zone:%s %s:%f time:%s", r.Name, r.Zone, r.Label, p.Value, p.Time.String())
	}

	avg := sum / float64(len(r.Monitors))
	log.Printf("%s average_%s:%f", r.Name, r.Label, avg)

	metrics := map[string]interface{}{
		"name":     r.Name,
		"zone":     r.Zone,
		"avg":      avg,
		"monitors": monitors,
	}

	for k, v := range r.AdditionalInfo {
		metrics[k] = v
	}

	return metrics
}

type iaasResources struct {
	Resources []*iaasResource
	Label     string
}

func (rs *iaasResources) toMetrics(percentiles []percentile) map[string]interface{} {
	var fs sort.Float64Slice
	routers := make([]interface{}, 0)
	total := float64(0)
	for _, t := range rs.Resources {
		metrics := t.toMetrics()
		avg := metrics["avg"].(float64)

		fs = append(fs, avg)
		total += avg

		routers = append(routers, metrics)
	}

	if len(fs) == 0 {
		result := map[string]interface{}{}
		result["max"] = float64(0)
		result["avg"] = float64(0)
		result["min"] = float64(0)
		for _, p := range percentiles {
			result[fmt.Sprintf("%spt", p.str)] = float64(0)
		}
		result["routers"] = routers
		return result
	}

	sort.Sort(fs)
	fl := float64(len(fs))
	result := map[string]interface{}{}
	result["max"] = fs[len(fs)-1]
	result["avg"] = total / fl
	result["min"] = fs[0]
	for _, p := range percentiles {
		result[fmt.Sprintf("%spt", p.str)] = fs[round(fl*(p.float))]
	}
	result["routers"] = routers
	return result
}

func round(f float64) int64 {
	return int64(math.Round(f)) - 1
}
