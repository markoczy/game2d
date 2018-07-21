package world

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func unmarshall(file string) (map[string]interface{}, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(dat, &data)
	return data, err
}

func getIntValue(data map[string]interface{}, selector string) (int, error) {
	ret := data[selector]
	if ret == nil {
		return 0, fmt.Errorf("Error in file: %s not found", selector)
	}
	k, ok := data[selector].(float64)
	if !ok {
		return 0, fmt.Errorf("Error in file: Bad format for %s ", selector)
	}
	return int(k), nil
}

func getStringValue(data map[string]interface{}, selector string) (string, error) {
	ret := data[selector]
	if ret == nil {
		return "", fmt.Errorf("Error in file: %s not found", selector)
	}
	k, ok := data[selector].(string)
	if !ok {
		return "", fmt.Errorf("Error in file: Bad format for %s ", selector)
	}
	return k, nil
}

func getBoolValue(data map[string]interface{}, selector string) (bool, error) {
	node := data[selector]
	if node == nil {
		return false, fmt.Errorf("Error in file: %s not found", selector)
	}
	ret, ok := node.(bool)
	if !ok {
		return false, fmt.Errorf("Error in file: Bad format for %s ", selector)
	}
	return ret, nil
}

func getChildNode(data map[string]interface{}, selector string) (map[string]interface{}, error) {
	node := data[selector]
	if node == nil {
		return nil, fmt.Errorf("Error in file: %s not found", selector)
	}
	ret, ok := data[selector].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error in file: %s is not a struct", selector)
	}
	return ret, nil
}

func getChildrenArray(data map[string]interface{}, selector string) ([]map[string]interface{}, error) {
	node := data[selector]
	if node == nil {
		return nil, fmt.Errorf("Error in file: %s not found", selector)
	}
	arr, ok := data[selector].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Error in file: %s is not an array", selector)
	}
	var ret []map[string]interface{}
	for i, v := range arr {
		cur, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Error in file: %s[%d] is not a struct", selector, i)
		}
		ret = append(ret, cur)
	}

	return ret, nil
}

// func getItemChildrenArray(data interface{}) ([]map[string]interface{}, error) {
// 	arr, ok := data.([]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("Error in file: %s is not an array", selector)
// 	}
// 	var ret []map[string]interface{}
// 	for i, v := range arr {
// 		cur, ok := v.(map[string]interface{})
// 		if !ok {
// 			return nil, fmt.Errorf("Error in file: %s[%d] is not a struct", selector, i)
// 		}
// 		ret = append(ret, cur)
// 	}

// 	return ret, nil
// }
