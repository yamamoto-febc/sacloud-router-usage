// Copyright 2022 The sacloud/sacloud-router-usage Authors
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
	"bytes"
	"context"
	"io"
	"log"
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func Test_fetchMetrics(t *testing.T) {
	log.SetOutput(io.Discard)

	type args struct {
		client iaasRouterAPI
		opts   *commandOpts
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				client: &stubIaasRouterAPI{},
				opts: &commandOpts{
					Prefix: []string{"test"},
					Zones:  []string{"is1a"},
					Item:   "out",
					Time:   0,
				},
			},
			want: map[string]interface{}{
				"avg":     0.,
				"max":     0.,
				"min":     0.,
				"routers": []interface{}{},
			},
			wantErr: false,
		},
		{
			name: "single resource - single value",
			args: args{
				client: &stubIaasRouterAPI{
					zone: "is1a",
					routers: map[string][]*iaas.Internet{
						"is1a": {
							{
								ID:   types.ID(1),
								Name: "test1",
							},
						},
					},
					activity: map[string]map[string]*iaas.RouterActivity{
						"is1a": {
							"1": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(1, 0),
										In:   0,
										Out:  3 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Prefix: []string{"test"},
					Zones:  []string{"is1a"},
					Item:   "out",
					Time:   1,
				},
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
			wantErr: false,
		},
		{
			name: "with percentiles",
			args: args{
				client: &stubIaasRouterAPI{
					zone: "is1a",
					routers: map[string][]*iaas.Internet{
						"is1a": {
							{
								ID:   types.ID(1),
								Name: "test1",
							},
						},
					},
					activity: map[string]map[string]*iaas.RouterActivity{
						"is1a": {
							"1": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(1, 0),
										In:   0,
										Out:  3 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Prefix: []string{"test"},
					Zones:  []string{"is1a"},
					Item:   "out",
					Time:   1,
					percentiles: []percentile{
						{
							str:   "90",
							float: 0.9,
						},
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
			wantErr: false,
		},
		{
			name: "single resource - multi values",
			args: args{
				client: &stubIaasRouterAPI{
					zone: "is1a",
					routers: map[string][]*iaas.Internet{
						"is1a": {
							{
								ID:   types.ID(1),
								Name: "test1",
							},
						},
					},
					activity: map[string]map[string]*iaas.RouterActivity{
						"is1a": {
							"1": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(1, 0),
										Out:  1 * 1000 * 1000,
									},
									{
										Time: time.Unix(2, 0),
										Out:  2 * 1000 * 1000,
									},
									{
										Time: time.Unix(3, 0),
										Out:  3 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Prefix: []string{"test"},
					Zones:  []string{"is1a"},
					Item:   "out",
					Time:   3,
				},
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
			wantErr: false,
		},
		{
			name: "multi resources - single value",
			args: args{
				client: &stubIaasRouterAPI{
					zone: "is1a",
					routers: map[string][]*iaas.Internet{
						"is1a": {
							{
								ID:   types.ID(1),
								Name: "test1",
							},
						},
						"is1b": {
							{
								ID:   types.ID(2),
								Name: "test2",
							},
						},
					},
					activity: map[string]map[string]*iaas.RouterActivity{
						"is1a": {
							"1": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(3, 0),
										Out:  2 * 1000 * 1000,
									},
								},
							},
						},
						"is1b": {
							"2": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(3, 0),
										Out:  4 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Prefix: []string{"test"},
					Zones:  []string{"is1a", "is1b"},
					Item:   "out",
					Time:   1,
				},
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
			wantErr: false,
		},
		{
			name: "multi resources - multi values",
			args: args{
				client: &stubIaasRouterAPI{
					zone: "is1a",
					routers: map[string][]*iaas.Internet{
						"is1a": {
							{
								ID:   types.ID(1),
								Name: "test1",
							},
						},
						"is1b": {
							{
								ID:   types.ID(2),
								Name: "test2",
							},
						},
					},
					activity: map[string]map[string]*iaas.RouterActivity{
						"is1a": {
							"1": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(1, 0),
										Out:  1 * 1000 * 1000,
									},
									{
										Time: time.Unix(2, 0),
										Out:  2 * 1000 * 1000,
									},
									{
										Time: time.Unix(3, 0),
										Out:  3 * 1000 * 1000,
									},
								},
							},
						},
						"is1b": {
							"2": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(4, 0),
										Out:  4 * 1000 * 1000,
									},
									{
										Time: time.Unix(5, 0),
										Out:  5 * 1000 * 1000,
									},
									{
										Time: time.Unix(6, 0),
										Out:  6 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Prefix: []string{"test"},
					Zones:  []string{"is1a", "is1b"},
					Item:   "out",
					Time:   3,
					percentiles: []percentile{
						{
							str:   "90",
							float: 0.9,
						},
					},
				},
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchMetrics(tt.args.client, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.want, got)
		})
	}
}

type stubIaasRouterAPI struct {
	zone     string
	routers  map[string][]*iaas.Internet
	activity map[string]map[string]*iaas.RouterActivity // map[zone][resource-id]
}

var _ iaasRouterAPI = (*stubIaasRouterAPI)(nil)

func (s *stubIaasRouterAPI) Find(_ context.Context, zone string, _ *iaas.FindCondition) (*iaas.InternetFindResult, error) {
	if _, ok := s.routers[zone]; !ok {
		return &iaas.InternetFindResult{}, nil
	}
	return &iaas.InternetFindResult{
		Total:    len(s.routers[zone]),
		From:     0,
		Count:    0,
		Internet: s.routers[zone],
	}, nil
}

func (s *stubIaasRouterAPI) MonitorRouter(_ context.Context, zone string, id types.ID, _ *iaas.MonitorCondition) (*iaas.RouterActivity, error) {
	if _, ok := s.activity[zone]; !ok {
		return &iaas.RouterActivity{}, nil
	}
	if _, ok := s.activity[zone][id.String()]; !ok {
		return &iaas.RouterActivity{}, nil
	}

	return s.activity[zone][id.String()], nil
}

func Test_outputMetrics(t *testing.T) {
	log.SetOutput(io.Discard)

	type args struct {
		metrics map[string]interface{}
		query   string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "without query",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: "",
			},
			wantW:   `{"90pt":1,"avg":2,"max":3,"min":4,"routers":[]}`,
			wantErr: false,
		},
		{
			name: "with query",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: ".avg",
			},
			wantW:   `2`,
			wantErr: false,
		},
		{
			name: "invalid query",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: "invalid-query",
			},
			wantW:   ``,
			wantErr: true,
		},
		{
			name: "query returns no value",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: ".not_exists",
			},
			wantW:   ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := outputMetrics(w, tt.args.metrics, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				require.Equal(t, tt.wantW+"\n", w.String())
			}
		})
	}
}
