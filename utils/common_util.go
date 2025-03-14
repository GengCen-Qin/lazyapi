package utils

import (
	"encoding/json"
	"log"
)

func FormatJSON(jsonString string) string {
    var jsonObj map[string]interface{}

    // Parse the JSON into a map
    err := json.Unmarshal([]byte(jsonString), &jsonObj)
    if err != nil {
        log.Fatalf("Error occured during unmarshalling. %s", err)
    }

    // Format the json
    formattedJSON, err := json.MarshalIndent(jsonObj, "", "    ")
    if err != nil {
        log.Fatalf("Error occured during marshalling. %s", err)
    }

    return string(formattedJSON)
}
