package main

import (
	"advent-of-code-2015/common"
	"encoding/json"
	"fmt"
)

func processValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case []interface{}:
		var newSlice []interface{}
		for _, elem := range val {
			if newElem := processValue(elem); newElem != nil {
				newSlice = append(newSlice, newElem)
			}
		}
		if len(newSlice) > 0 {
			return newSlice
		} else {
			return nil
		}

	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range val {
			if value == "red" {
				return nil
			}
			if newValue := processValue(value); newValue != nil {
				newMap[key] = newValue
			}
		}
		if len(newMap) > 0 {
			return newMap
		} else {
			return nil
		}
	}

	return v
}

func removeRed(s string) string {
	var data []interface{}
	json.Unmarshal([]byte(s), &data)
	//before, _ := json.MarshalIndent(data, "", "   ")
	//fmt.Println(string(before))

	newdata := processValue(data)

	str, _ := json.Marshal(newdata)
	//str, _ := json.MarshalIndent(newdata, "", "   ")
	//fmt.Println(string(str))
	return string(str)
}

func main() {
	lines := common.ReadInputLines()

	str := removeRed(lines[0])

	result := countNumbers(str)

	fmt.Println("Result:", result)
}
