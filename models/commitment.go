package models

import (
	"fmt"
	"strings"

	"github.com/CollActionteam/collaction_backend/utils"
)

type CommitmentOption struct {
	Id          string             `json:"id"`
	Label       string             `json:"label"`
	Description string             `json:"description"`
	Requires    []CommitmentOption `json:"requires,omitempty"`
}

func ValidateCommitments(commitments []string, rootOptions []CommitmentOption) error {
	err := validateCommitments(&commitments, rootOptions, false)
	if err == nil && len(commitments) > 0 {
		err = fmt.Errorf("commitments \"%s\" not in options", strings.Join(commitments, ", "))
	}
	return err
}

func validateCommitments(commitments *[]string, rootOptions []CommitmentOption, requireAll bool) (err error) {
	if len(rootOptions) == 1 {
		// Is root (of subtree)
		option := rootOptions[0]
		if i := utils.IndexOf(*commitments, option.Id); i != -1 {
			utils.Remove(commitments, i)
			// Validate branches
			err = validateCommitments(commitments, option.Requires, true)
		} else if requireAll {
			err = fmt.Errorf("required commitment \"%s\"", option.Id)
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
