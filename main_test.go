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
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	usage "github.com/sacloud/sacloud-usage-lib"
	"github.com/stretchr/testify/require"
)

type stubIaasRouterAPI struct {
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

func Test_fetchResources(t *testing.T) {
	type args struct {
		client iaasRouterAPI
		opts   *commandOpts
	}
	tests := []struct {
		name    string
		args    args
		want    *usage.Resources
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				client: &stubIaasRouterAPI{
					routers:  map[string][]*iaas.Internet{},
					activity: nil,
				},
				opts: &commandOpts{Option: &usage.Option{}},
			},
			want: &usage.Resources{
				Resources: nil,
				Label:     "routers",
			},
			wantErr: false,
		},
		{
			name: "single resource - single value",
			args: args{
				client: &stubIaasRouterAPI{
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
										Out:  1 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Option: &usage.Option{
						Prefix: []string{"test"},
						Zones:  []string{"is1a"},
						Time:   1,
					},
					Item: "out",
				},
			},
			want: &usage.Resources{
				Resources: []*usage.Resource{
					{
						ID:   1,
						Name: "test1",
						Zone: "is1a",
						Monitors: []usage.MonitorValue{
							{
								Time:  time.Unix(1, 0),
								Value: 1,
							},
						},
						Label:          "traffic",
						AdditionalInfo: nil,
					},
				},
				Label: "routers",
			},
		},
		{
			name: "single resource - multi values",
			args: args{
				client: &stubIaasRouterAPI{
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
					Option: &usage.Option{
						Prefix: []string{"test"},
						Zones:  []string{"is1a"},
						Time:   3,
					},
					Item: "out",
				},
			},
			want: &usage.Resources{
				Resources: []*usage.Resource{
					{
						ID:   1,
						Name: "test1",
						Zone: "is1a",
						Monitors: []usage.MonitorValue{
							{
								Time:  time.Unix(1, 0),
								Value: 1,
							},
							{
								Time:  time.Unix(2, 0),
								Value: 2,
							},
							{
								Time:  time.Unix(3, 0),
								Value: 3,
							},
						},
						Label:          "traffic",
						AdditionalInfo: nil,
					},
				},
				Label: "routers",
			},
		},
		{
			name: "multi resources - single value",
			args: args{
				client: &stubIaasRouterAPI{
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
								},
							},
						},
						"is1b": {
							"2": {
								Values: []*iaas.MonitorRouterValue{
									{
										Time: time.Unix(2, 0),
										Out:  2 * 1000 * 1000,
									},
								},
							},
						},
					},
				},
				opts: &commandOpts{
					Option: &usage.Option{
						Prefix: []string{"test"},
						Zones:  []string{"is1a", "is1b"},
						Time:   1,
					},
					Item: "out",
				},
			},
			want: &usage.Resources{
				Resources: []*usage.Resource{
					{
						ID:   1,
						Name: "test1",
						Zone: "is1a",
						Monitors: []usage.MonitorValue{
							{
								Time:  time.Unix(1, 0),
								Value: 1,
							},
						},
						Label:          "traffic",
						AdditionalInfo: nil,
					},
					{
						ID:   2,
						Name: "test2",
						Zone: "is1b",
						Monitors: []usage.MonitorValue{
							{
								Time:  time.Unix(2, 0),
								Value: 2,
							},
						},
						Label:          "traffic",
						AdditionalInfo: nil,
					},
				},
				Label: "routers",
			},
		},
		{
			name: "multi resources - multi values",
			args: args{
				client: &stubIaasRouterAPI{
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
					Option: &usage.Option{
						Prefix: []string{"test"},
						Zones:  []string{"is1a", "is1b"},
						Time:   3,
					},
					Item: "out",
				},
			},
			want: &usage.Resources{
				Resources: []*usage.Resource{
					{
						ID:   1,
						Name: "test1",
						Zone: "is1a",
						Monitors: []usage.MonitorValue{
							{
								Time:  time.Unix(1, 0),
								Value: 1,
							},
							{
								Time:  time.Unix(2, 0),
								Value: 2,
							},
							{
								Time:  time.Unix(3, 0),
								Value: 3,
							},
						},
						Label:          "traffic",
						AdditionalInfo: nil,
					},
					{
						ID:   2,
						Name: "test2",
						Zone: "is1b",
						Monitors: []usage.MonitorValue{
							{
								Time:  time.Unix(4, 0),
								Value: 4,
							},
							{
								Time:  time.Unix(5, 0),
								Value: 5,
							},
							{
								Time:  time.Unix(6, 0),
								Value: 6,
							},
						},
						Label:          "traffic",
						AdditionalInfo: nil,
					},
				},
				Label: "routers",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchResources(context.Background(), tt.args.client, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Option = nil
			tt.want.Option = nil
			require.EqualValues(t, tt.want, got)
		})
	}
}
