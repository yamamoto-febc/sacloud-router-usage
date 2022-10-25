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
# TODO
```

## License

`sacloud-router-usage` Copyright (C) 2022 The sacloud/sacloud-router-usage authors.

This project is published under [MIT](LICENSE).