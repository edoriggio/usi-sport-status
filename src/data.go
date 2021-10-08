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
