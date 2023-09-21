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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/itchyny/gojq"
	"github.com/sacloud/sacloud-router-usage/version"
)

func PrintVersion() {
	fmt.Printf(`%s %s
Compiler: %s %s
`,
		os.Args[0],
		version.Version,
		runtime.Compiler,
		runtime.Version())
}

func OutputMetrics(w io.Writer, metrics map[string]interface{}, query string) error {
	if query == "" {
		v, _ := json.Marshal(metrics)
		fmt.Fprintln(w, string(v))
		return nil
	}

	parsed, err := gojq.Parse(query)
	if err != nil {
		return err
	}
	iter := parsed.Run(metrics)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}
		if v == nil {
			return fmt.Errorf("%s not found in result", query)
		}
		j2, _ := json.Marshal(v)
		fmt.Fprintln(w, string(j2))
	}

	return nil
}
