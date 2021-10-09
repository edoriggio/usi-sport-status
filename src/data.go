// Copyright 2021 Edoardo Riggio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"time"
	"encoding/json"
)

type Result struct {
	ValidTo   string `json:"valid-to"`
	Available bool   `json:"available"`
}

func checkIfExpired() bool {
	var result Result

	jsonFile, _ := os.ReadFile("../data/result.json")
	json.Unmarshal(jsonFile, &result)

	now := time.Now().Local()
	parsed, _ := time.Parse("01-02-2006 15:04:05 MST", result.ValidTo)

	if !now.Before(parsed) {
		var newResult Result

		expTime := time.Now().Local().Add(time.Hour * time.Duration(2)).Format("01-02-2006 15:04:05 MST")
	  available := checkWebsite()

		newResult.ValidTo = expTime
		newResult.Available = available

		file, _ := json.MarshalIndent(newResult, "", " ")
		os.WriteFile("../data/result.json", file, 0644)

		return available
	}

	return result.Available
}
