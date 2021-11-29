package model

import (
	"fmt"
	"strings"
)

type CommitmentOption struct {
	Label    string             `json:"label"`
	Requires []CommitmentOption `json:"requires,omitempty"`
}

func ValidateCommitments(commitments []string, rootOptions []CommitmentOption) error {
	err := validateCommitments(&commitments, rootOptions, false)
	if err == nil && len(commitments) > 0 {
		err = fmt.Errorf("commitments \"%s\" not in options", strings.Join(commitments, ", "))
	}
	return err
}

// TODO move to shared utils module
func indexOf(s []string, v string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// TODO move to shared utils module
func remove(s *[]string, i int) {
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
}

func validateCommitments(commitments *[]string, rootOptions []CommitmentOption, requireAll bool) (err error) {
	if len(rootOptions) == 1 {
		// Is root (of subtree)
		option := rootOptions[0]
		if i := indexOf(*commitments, option.Label); i != -1 {
			remove(commitments, i)
			// Validate branches
			err = validateCommitments(commitments, option.Requires, true)
		} else if requireAll {
			err = fmt.Errorf("required commitment \"%s\"", option.Label)
		} else {
			// Validate branches
			err = validateCommitments(commitments, option.Requires, false)
		}
	} else {
		// Are branches
		for _, option := range rootOptions {
			err = validateCommitments(commitments, []CommitmentOption{option}, requireAll)
			if err != nil {
				break
			}
		}
	}
	return
}
