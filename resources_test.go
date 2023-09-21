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
	"io"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_iaasResources_toMetrics(t *testing.T) {
	log.SetOutput(io.Discard)

	type args struct {
		resources   *iaasResources
		percentiles []percentile
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "empty",
			args: args{
				resources:   &iaasResources{Label: "routers"},
				percentiles: []percentile{},
			},
			want: map[string]interface{}{
				"avg":     0.,
				"max":     0.,
				"min":     0.,
				"routers": []interface{}{},
			},
		},
		{
			name: "single resource - single value",
			args: args{
				resources: &iaasResources{
					Label: "routers",
					Resources: []*iaasResource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []monitorValue{
								{Time: time.Unix(1, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
				percentiles: []percentile{},
			},
			want: map[string]interface{}{
				"avg": 3.,
				"max": 3.,
				"min": 3.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  3.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(1, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "with percentiles",
			args: args{
				resources: &iaasResources{
					Label: "routers",
					Resources: []*iaasResource{
						{

							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []monitorValue{
								{Time: time.Unix(1, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
				percentiles: []percentile{
					{
						str:   "90",
						float: 0.9,
					},
				},
			},
			want: map[string]interface{}{
				"avg":  3.,
				"max":  3.,
				"min":  3.,
				"90pt": 3.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  3.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(1, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "single resource - multi values",
			args: args{
				resources: &iaasResources{
					Label: "routers",
					Resources: []*iaasResource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []monitorValue{
								{Time: time.Unix(1, 0), Value: 1},
								{Time: time.Unix(2, 0), Value: 2},
								{Time: time.Unix(3, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
				percentiles: []percentile{},
			},
			want: map[string]interface{}{
				"avg": 2.,
				"max": 2.,
				"min": 2.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  2.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 1.,
								"time":    time.Unix(1, 0).String(),
							},
							map[string]interface{}{
								"traffic": 2.,
								"time":    time.Unix(2, 0).String(),
							},
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "multi resources - single value",
			args: args{
				resources: &iaasResources{
					Label: "routers",
					Resources: []*iaasResource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []monitorValue{
								{Time: time.Unix(3, 0), Value: 2},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
						{
							ID:   2,
							Name: "test2",
							Zone: "is1b",
							Monitors: []monitorValue{
								{Time: time.Unix(3, 0), Value: 4},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
				percentiles: []percentile{},
			},
			want: map[string]interface{}{
				"avg": 3.,
				"max": 4.,
				"min": 2.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  2.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 2.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
					map[string]interface{}{
						"name": "test2",
						"zone": "is1b",
						"avg":  4.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 4.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "multi resources - multi values",
			args: args{
				resources: &iaasResources{
					Label: "routers",
					Resources: []*iaasResource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []monitorValue{
								{Time: time.Unix(1, 0), Value: 1},
								{Time: time.Unix(2, 0), Value: 2},
								{Time: time.Unix(3, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
						{
							ID:   2,
							Name: "test2",
							Zone: "is1b",
							Monitors: []monitorValue{
								{Time: time.Unix(4, 0), Value: 4},
								{Time: time.Unix(5, 0), Value: 5},
								{Time: time.Unix(6, 0), Value: 6},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
				percentiles: []percentile{{str: "90", float: 0.9}},
			},
			want: map[string]interface{}{
				"avg":  3.5,
				"max":  5.,
				"min":  2.,
				"90pt": 5.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  2.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 1.,
								"time":    time.Unix(1, 0).String(),
							},
							map[string]interface{}{
								"traffic": 2.,
								"time":    time.Unix(2, 0).String(),
							},
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
					map[string]interface{}{
						"name": "test2",
						"zone": "is1b",
						"avg":  5.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 4.,
								"time":    time.Unix(4, 0).String(),
							},
							map[string]interface{}{
								"traffic": 5.,
								"time":    time.Unix(5, 0).String(),
							},
							map[string]interface{}{
								"traffic": 6.,
								"time":    time.Unix(6, 0).String(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.resources.toMetrics(tt.args.percentiles)
			require.EqualValues(t, tt.want, got)
		})
	}
}
