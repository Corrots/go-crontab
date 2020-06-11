package model

import "fmt"

type Job struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	Expression string `json:"expression"`
}

func (j *Job) Validation() error {
	if j.Name == "" || j.Command == "" || j.Expression == "" {
		return fmt.Errorf("invalid job field")
	}
	return nil
}
