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
      --time=           Get average traffic for a specified amount of time (default: 3)
      --item=[in|out]   Item name (default: in)
      --prefix=         prefix for router names. prefix accepts more than one.
      --zone=           zone name
      --percentile-set= percentiles to dispaly (default: 99,95,90,75)
  -v, --version         Show version
      --query=          jq style query to result and display
      --env-from=       load environment values from this file

Help Options:
  -h, --help            Show this help message
```

## Examples

```bash
% ./sacloud-router-usage --prefix example- --item out  --zone tk1b 
2022/10/25 12:02:06 example-r1 zone:tk1b traffic:5305794.000000 time:2022-10-25 11:50:00 +0900 JST
2022/10/25 12:02:06 example-r1 zone:tk1b traffic:5305794.000000 time:2022-10-25 11:55:00 +0900 JST
2022/10/25 12:02:06 example-r1 zone:tk1b traffic:0.000000 time:2022-10-25 12:00:00 +0900 JST
2022/10/25 12:02:06 example-r1 average_traffic:3537196.000000
{"75pt":3537196,"90pt":3537196,"95pt":3537196,"99pt":3537196,"avg":3537196,"max":3537196,"min":3537196,"routers":[{"avg":3537196,"monitors":[{"time":"2022-10-25 11:50:00 +0900 JST","traffic":5305794},{"time":"2022-10-25 11:55:00 +0900 JST","traffic":5305794},{"time":"2022-10-25 12:00:00 +0900 JST","traffic":0}],"name":"example-r1","zone":"tk1b"}]}
```

## License

`sacloud-router-usage` Copyright (C) 2022 The sacloud/sacloud-router-usage authors.

This project is published under [MIT](LICENSE).