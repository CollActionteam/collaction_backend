package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	jsonCommitmentOptions = []byte(`[
        {
            "label": "5/7 days a week"
        },
        {
            "label": "vegan",
            "requires": [
                {
                    "label": "vegetarian",
                    "requires": [
                        {
                            "label": "pescatarian"
                        },
                        {
                            "label": "no beef"
                        }
                    ]
                },
                {
                    "label": "no dairy",
                    "requires": [
                        {
                            "label": "no cheese"
                        }
                    ]
                }
            ]
        }
    ]`)
	validCommitments          = []string{"vegetarian", "pescatarian", "no beef"}
	invalidCommitmentsMissing = []string{"vegan", "vegetarian", "pescatarian", "no beef"}
	invalidCommitmentsUnknown = []string{"only caviar"}
)

func TestCommitments(t *testing.T) {
	t.Run("Deserialize commitment options from JSON", func(t *testing.T) {
		var options []CommitmentOption
		err := json.Unmarshal(jsonCommitmentOptions, &options)
		assert.Nil(t, err)
	})

	t.Run("Valid commitment", func(t *testing.T) {
		var options []CommitmentOption
		json.Unmarshal(jsonCommitmentOptions, &options)
		err := ValidateCommitments(validCommitments, options)
		assert.Nil(t, err)
	})

	t.Run("Required commitments missing", func(t *testing.T) {
		var options []CommitmentOption
		json.Unmarshal(jsonCommitmentOptions, &options)
		err := ValidateCommitments(invalidCommitmentsMissing, options)
		assert.NotNil(t, err)
	})

	t.Run("Invalid commitment", func(t *testing.T) {
		var options []CommitmentOption
		json.Unmarshal(jsonCommitmentOptions, &options)
		err := ValidateCommitments(invalidCommitmentsUnknown, options)
		assert.NotNil(t, err)
	})
}
