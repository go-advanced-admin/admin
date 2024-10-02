package form

import "encoding/json"

type HTMLType string

func ConvertFormDataToHTMLTypeMap(data map[string][]string) (map[string]HTMLType, error) {
	result := make(map[string]HTMLType)

	for key, values := range data {
		if len(values) > 1 {
			jsonEncoded, err := json.Marshal(values)
			if err != nil {
				return result, err
			}
			result[key] = HTMLType(jsonEncoded)
		} else if len(values) == 1 {
			result[key] = HTMLType(values[0])
		}
	}

	return result, nil
}
