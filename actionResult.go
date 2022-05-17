package main

import "encoding/json"

type ActionResult struct {
	Succeeded    bool
	ErrorMessage string
}

func GetSucceeded() ActionResult {
	return ActionResult{
		Succeeded: true,
	}
}
func GetFailed(errorMessage string) ActionResult {
	return ActionResult{
		Succeeded:    false,
		ErrorMessage: errorMessage,
	}
}
func GetSucceededJson() string {
	result, _ := json.Marshal(GetSucceeded())
	return string(result)
}
func GetFailedJson(errorMessage string) string {
	result, _ := json.Marshal(GetFailed(errorMessage))
	return string(result)
}
