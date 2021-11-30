package utils

import (
	"encoding/json"
	"testing"

	"github.com/CollActionteam/collaction_backend/models"
	"github.com/stretchr/testify/assert"
)

var (
	jsonCommitmentOptions = []byte(`[
        {
			"id": "5-7-days-a-week",
            "label": "5/7 days a week",
			"description": "Lorem ipsum foo bar"
        },
        {
			"id": "vegan",
            "label": "vegan",
			"description": "Lorem ipsum foo bar",
            "requires": [
                {
					"id": "vegetarian",
                    "label": "vegetarian",
					"description": "Lorem ipsum foo bar",
                    "requires": [
                        {
							"id": "pescatarian",
                            "label": "pescatarian",
							"description": "Lorem ipsum foo bar",
							"requires": [
								{
									"id": "no-beef",
									"label": "no beef",
									"description": "Lorem ipsum foo bar"
								}
							]
                        }
                    ]
                },
                {
					"id": "no-dairy",
                    "label": "no dairy",
					"description": "Lorem ipsum foo bar",
                    "requires": [
                        {
							"id": "no-cheese",
                            "label": "no cheese",
							"description": "Lorem ipsum foo bar"
                        }
                    ]
                }
            ]
        }
    ]`)
	validCommitments          = []string{"vegetarian", "pescatarian", "no-beef"}
	invalidCommitmentsMissing = []string{"vegan", "vegetarian", "pescatarian", "no-beef"}
	invalidCommitmentsUnknown = []string{"only-caviar"}
)

func TestCommitments(t *testing.T) {
	t.Run("Deserialize commitment options from JSON", func(t *testing.T) {
		var options []models.CommitmentOption
		err := json.Unmarshal(jsonCommitmentOptions, &options)
		assert.Nil(t, err)
	})

	t.Run("Valid commitment", func(t *testing.T) {
		var options []models.CommitmentOption
		json.Unmarshal(jsonCommitmentOptions, &options)
		err := ValidateCommitments(validCommitments, options)
		assert.Nil(t, err)
	})

	t.Run("Required commitments missing", func(t *testing.T) {
		var options []models.CommitmentOption
		json.Unmarshal(jsonCommitmentOptions, &options)
		err := ValidateCommitments(invalidCommitmentsMissing, options)
		assert.NotNil(t, err)
	})

	t.Run("Invalid commitment", func(t *testing.T) {
		var options []models.CommitmentOption
		json.Unmarshal(jsonCommitmentOptions, &options)
		err := ValidateCommitments(invalidCommitmentsUnknown, options)
		assert.NotNil(t, err)
	})
}
