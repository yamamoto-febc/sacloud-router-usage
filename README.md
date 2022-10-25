# sacloud/sacloud-router-usage

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/sacloud-router-usage.svg)](https://pkg.go.dev/github.com/sacloud/sacloud-router-usage)
[![Tests](https://github.com/sacloud/sacloud-router-usage/workflows/Tests/badge.svg)](https://github.com/sacloud/sacloud-router-usage/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/sacloud-router-usage)](https://goreportcard.com/report/github.com/sacloud/sacloud-router-usage)

[sacloud/sacloud-cpu-usage](https://github.com/sacloud/sacloud-cpu-usage)のルータ+スイッチ版

## 概要

スイッチ+ルータの使用状況をさくらのクラウドAPI経由で取得しMax/Min/Averageを算出します。

## Usage

```bash
Usage:
  sacloud-router-usage [OPTIONS]

Application Options:
      --time=           Get average usage for a specified amount of time (default: 3)
      --prefix=         Prefix for router names. prefix accepts more than one.
      --zone=           Zone name
      --percentile-set= Percentiles to dispaly (default: 99,95,90,75)
  -v, --version         Show version
      --query=          jq style query to result and display
      --env-from=       Load envrionment values from this file

Help Options:
  -h, --help            Show this help message
```

## Examples

```bash
% ./sacloud-router-usage --prefix example- --zone tk1b --zone is1b 
2022/10/25 10:57:25 example-r1 zone:tk1b traffic:0.000000 time:2022-10-25 10:45:00 +0900 JST
2022/10/25 10:57:25 example-r1 zone:tk1b traffic:0.000000 time:2022-10-25 10:50:00 +0900 JST
2022/10/25 10:57:25 example-r1 zone:tk1b traffic:0.000000 time:2022-10-25 10:55:00 +0900 JST
2022/10/25 10:57:25 example-r1 average_traffic:0.000000
2022/10/25 10:57:25 example-r2 zone:is1b traffic:0.000000 time:2022-10-25 10:40:00 +0900 JST
2022/10/25 10:57:25 example-r2 zone:is1b traffic:0.000000 time:2022-10-25 10:45:00 +0900 JST
2022/10/25 10:57:25 example-r2 zone:is1b traffic:0.000000 time:2022-10-25 10:50:00 +0900 JST
2022/10/25 10:57:25 example-r2 average_traffic:0.000000
{"75pt":0,"90pt":0,"95pt":0,"99pt":0,"avg":0,"max":0,"min":0,"routers":[{"avg":0,"monitors":[{"time":"2022-10-25 10:45:00 +0900 JST","traffic":0},{"time":"2022-10-25 10:50:00 +0900 JST","traffic":0},{"time":"2022-10-25 10:55:00 +0900 JST","traffic":0}],"name":"example-r1","zone":"tk1b"},{"avg":0,"monitors":[{"time":"2022-10-25 10:40:00 +0900 JST","traffic":0},{"time":"2022-10-25 10:45:00 +0900 JST","traffic":0},{"time":"2022-10-25 10:50:00 +0900 JST","traffic":0}],"name":"example-r2","zone":"is1b"}]}
```

## License

`sacloud-router-usage` Copyright (C) 2022 The sacloud/sacloud-router-usage authors.

This project is published under [MIT](LICENSE).