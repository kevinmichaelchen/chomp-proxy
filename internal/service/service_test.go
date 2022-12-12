package service

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	chompv1beta1 "go.buf.build/bufbuild/connect-go/kevinmichaelchen/chompapis/chomp/v1beta1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	h := http.Header{
		"Api_key": []string{"foobar"},
	}
	key, err := getAPIKey(h)
	require.NoError(t, err)
	require.Equal(t, "foobar", key)
}

func TestUnmarshal(t *testing.T) {
	s := `
{
  "barcode": "string",
  "name": "string",
  "brand": "string",
  "ingredients": "string",
  "package": {
    "quantity": 0,
    "size": "string"
  },
  "serving": {
    "size": "string",
    "measurement_unit": "string",
    "size_fulltext": "string"
  },
  "categories": [
    "string"
  ],
  "nutrients": [
    {
      "name": "string",
      "per_100g": 0,
      "measurement_unit": "string",
      "rank": 0,
      "data_points": 0,
      "description": "string"
    }
  ],
  "diet_labels": {
    "vegan": {
      "name": "string",
      "is_compatible": true,
      "compatibility_level": 0,
      "confidence": 0,
      "confidence_description": "string"
    },
    "vegetarian": {
      "name": "string",
      "is_compatible": true,
      "compatibility_level": 0,
      "confidence": 0,
      "confidence_description": "string"
    },
    "gluten_free": {
      "name": "string",
      "is_compatible": true,
      "compatibility_level": 0,
      "confidence": 0,
      "confidence_description": "string"
    }
  },
  "diet_flags": [
    {
      "ingredient": "string",
      "ingredient_description": "string",
      "diet_label": "string",
      "is_compatible": "string",
      "compatibility_level": 0,
      "compatibility_description": "string",
      "is_allergen": true
    }
  ],
  "packaging_photos": {
    "front": {
      "small": "string",
      "thumb": "string",
      "display": "string"
    },
    "nutrition": {
      "small": "string",
      "thumb": "string",
      "display": "string"
    },
    "ingredients": {
      "small": "string",
      "thumb": "string",
      "display": "string"
    }
  },
  "allergens": [
    "string"
  ],
  "brand_list": [
    "string"
  ],
  "countries": [
    "string"
  ],
  "country_details": {
    "english_speaking": 0,
    "non_english_speaking": 0
  },
  "palm_oil_ingredients": [
    "string"
  ],
  "ingredient_list": [
    "string"
  ],
  "has_english_ingredients": true,
  "minerals": [
    "string"
  ],
  "traces": [
    "string"
  ],
  "vitamins": [
    "string"
  ],
  "description": "string",
  "keywords": [
    "string"
  ]
}
`

	var actual chompv1beta1.Food
	err := protojson.Unmarshal([]byte(s), &actual)
	require.NoError(t, err)

	expected := &chompv1beta1.Food{
		Barcode:     "string",
		Name:        "string",
		Brand:       "string",
		Ingredients: "string",
		Package: &chompv1beta1.Package{
			Quantity: 0,
			Size:     "string",
		},
		Serving: &chompv1beta1.Serving{
			Size:            "string",
			MeasurementUnit: "string",
			SizeFulltext:    "string",
		},
		Categories: []string{"string"},
		Nutrients: []*chompv1beta1.Nutrient{
			{
				Name:            "string",
				Per_100G:        0,
				MeasurementUnit: "string",
				Rank:            0,
				DataPoints:      0,
				Description:     "string",
			},
		},
		DietLabels: &chompv1beta1.DietLabels{
			Vegan: &chompv1beta1.DietLabel{
				Name:                  "string",
				IsCompatible:          true,
				CompatibilityLevel:    0,
				Confidence:            0,
				ConfidenceDescription: "string",
			},
			Vegetarian: &chompv1beta1.DietLabel{
				Name:                  "string",
				IsCompatible:          true,
				CompatibilityLevel:    0,
				Confidence:            0,
				ConfidenceDescription: "string",
			},
			GlutenFree: &chompv1beta1.DietLabel{
				Name:                  "string",
				IsCompatible:          true,
				CompatibilityLevel:    0,
				Confidence:            0,
				ConfidenceDescription: "string",
			},
		},
		DietFlags: []*chompv1beta1.DietFlag{
			{
				Ingredient:               "string",
				IngredientDescription:    "string",
				DietLabel:                "string",
				IsCompatible:             "string",
				CompatibilityLevel:       0,
				CompatibilityDescription: "string",
				IsAllergen:               true,
			},
		},
		PackagingPhotos: &chompv1beta1.PackagingPhotos{
			Front: &chompv1beta1.Photo{
				Small:   "string",
				Thumb:   "string",
				Display: "string",
			},
			Nutrition: &chompv1beta1.Photo{
				Small:   "string",
				Thumb:   "string",
				Display: "string",
			},
			Ingredients: &chompv1beta1.Photo{
				Small:   "string",
				Thumb:   "string",
				Display: "string",
			},
		},
		Allergens: []string{"string"},
		BrandList: []string{"string"},
		Countries: []string{"string"},
		CountryDetails: &chompv1beta1.CountryDetails{
			EnglishSpeaking:    0,
			NonEnglishSpeaking: 0,
		},
		PalmOilIngredients:    []string{"string"},
		IngredientList:        []string{"string"},
		HasEnglishIngredients: true,
		Minerals:              []string{"string"},
		Traces:                []string{"string"},
		Vitamins:              []string{"string"},
		Description:           "string",
		Keywords:              []string{"string"},
	}
	diff := cmp.Diff(expected, &actual, protocmp.Transform())
	require.Empty(t, diff)
}
