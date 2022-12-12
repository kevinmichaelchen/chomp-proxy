package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	chompv1beta1 "go.buf.build/bufbuild/connect-go/kevinmichaelchen/chompapis/chomp/v1beta1"
	"io"
	"net/http"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetFood(
	ctx context.Context,
	req *connect.Request[chompv1beta1.GetFoodRequest],
) (*connect.Response[chompv1beta1.GetFoodResponse], error) {
	// Get API key
	apiKey, err := getAPIKey(req.Header())
	if err != nil {
		return nil, err
	}

	logrus.WithField("barcode", req.Msg.GetCode()).Info("Retrieving food...")

	barcode := req.Msg.GetCode()
	url := fmt.Sprintf("https://chompthis.com/api/v2/food/branded/barcode.php?api_key=%s&code=%s", apiKey, barcode)

	// Hit Chomp API
	apiRes, err := hitAPI(url)
	if err != nil {
		return nil, err
	}

	// Check for Not Found
	if len(apiRes.Items) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no foods found"))
	}

	res := &chompv1beta1.GetFoodResponse{
		Food: convert(apiRes.Items[0]),
	}

	out := connect.NewResponse(res)
	out.Header().Set("API-Version", "v1beta1")
	return out, nil
}

func (s *Service) ListFoods(
	ctx context.Context,
	req *connect.Request[chompv1beta1.ListFoodsRequest],
) (*connect.Response[chompv1beta1.ListFoodsResponse], error) {
	// Get API key
	apiKey, err := getAPIKey(req.Header())
	if err != nil {
		return nil, err
	}

	name := req.Msg.GetName()
	url := fmt.Sprintf("https://chompthis.com/api/v2/food/branded/name.php?api_key=%s&name=%s", apiKey, name)

	// Hit Chomp API
	apiRes, err := hitAPI(url)
	if err != nil {
		return nil, err
	}

	// Check for Not Found
	if len(apiRes.Items) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no foods found"))
	}

	var items []*chompv1beta1.Food
	for _, item := range apiRes.Items {
		items = append(items, convert(item))
	}
	res := &chompv1beta1.ListFoodsResponse{
		Items: items,
	}

	logrus.WithField("query", req.Msg.GetName()).Info("Retrieving food...")
	out := connect.NewResponse(res)
	out.Header().Set("API-Version", "v1beta1")
	return out, nil
}

func getAPIKey(headers http.Header) (string, error) {
	h, ok := headers["api_key"]
	if !ok || len(h) == 0 {
		return "", connect.NewError(connect.CodePermissionDenied, errors.New("missing api_key header"))
	}
	return h[0], nil
}

func hitAPI(url string) (*ChompResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request against Chomp API: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response bytes from Chomp API: %w", err)
	}

	var res ChompResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload from Chomp API: %w", err)
	}

	return &res, nil
}

func convert(in ChompFoodItem) *chompv1beta1.Food {
	var nutrients []*chompv1beta1.Nutrient
	for _, n := range in.Nutrients {
		nutrients = append(nutrients, &chompv1beta1.Nutrient{
			Name:            n.Name,
			Per_100G:        int32(n.Per100G),
			MeasurementUnit: n.MeasurementUnit,
			Rank:            int32(n.Rank),
			DataPoints:      int32(n.DataPoints),
			Description:     n.Description,
		})
	}
	var dietFlags []*chompv1beta1.DietFlag
	for _, d := range in.DietFlags {
		dietFlags = append(dietFlags, &chompv1beta1.DietFlag{
			Ingredient:               d.Ingredient,
			IngredientDescription:    d.IngredientDescription,
			DietLabel:                d.DietLabel,
			IsCompatible:             d.IsCompatible,
			CompatibilityLevel:       int32(d.CompatibilityLevel),
			CompatibilityDescription: d.CompatibilityDescription,
			IsAllergen:               d.IsAllergen,
		})
	}
	return &chompv1beta1.Food{
		Barcode:     in.Barcode,
		Name:        in.Name,
		Brand:       in.Brand,
		Ingredients: in.Ingredients,
		Package: &chompv1beta1.Package{
			Quantity: int32(in.Package.Quantity),
			Size:     in.Package.Size,
		},
		Serving: &chompv1beta1.Serving{
			Size:            in.Serving.Size,
			MeasurementUnit: in.Serving.MeasurementUnit,
			SizeFulltext:    in.Serving.SizeFulltext,
		},
		Categories: in.Categories,
		Nutrients:  nutrients,
		DietLabels: &chompv1beta1.DietLabels{
			Vegan: &chompv1beta1.DietLabel{
				Name:                  in.DietLabels.Vegan.Name,
				IsCompatible:          in.DietLabels.Vegan.IsCompatible,
				CompatibilityLevel:    int32(in.DietLabels.Vegan.CompatibilityLevel),
				Confidence:            int32(in.DietLabels.Vegan.Confidence),
				ConfidenceDescription: in.DietLabels.Vegan.ConfidenceDescription,
			},
			Vegetarian: &chompv1beta1.DietLabel{
				Name:                  in.DietLabels.Vegetarian.Name,
				IsCompatible:          in.DietLabels.Vegetarian.IsCompatible,
				CompatibilityLevel:    int32(in.DietLabels.Vegetarian.CompatibilityLevel),
				Confidence:            int32(in.DietLabels.Vegetarian.Confidence),
				ConfidenceDescription: in.DietLabels.Vegetarian.ConfidenceDescription,
			},
			GlutenFree: &chompv1beta1.DietLabel{
				Name:                  in.DietLabels.GlutenFree.Name,
				IsCompatible:          in.DietLabels.GlutenFree.IsCompatible,
				CompatibilityLevel:    int32(in.DietLabels.GlutenFree.CompatibilityLevel),
				Confidence:            int32(in.DietLabels.GlutenFree.Confidence),
				ConfidenceDescription: in.DietLabels.GlutenFree.ConfidenceDescription,
			},
		},
		DietFlags: dietFlags,
		PackagingPhotos: &chompv1beta1.PackagingPhotos{
			Front: &chompv1beta1.Photo{
				Small:   in.PackagingPhotos.Front.Small,
				Thumb:   in.PackagingPhotos.Front.Thumb,
				Display: in.PackagingPhotos.Front.Display,
			},
			Nutrition: &chompv1beta1.Photo{
				Small:   in.PackagingPhotos.Nutrition.Small,
				Thumb:   in.PackagingPhotos.Nutrition.Thumb,
				Display: in.PackagingPhotos.Nutrition.Display,
			},
			Ingredients: &chompv1beta1.Photo{
				Small:   in.PackagingPhotos.Ingredients.Small,
				Thumb:   in.PackagingPhotos.Ingredients.Thumb,
				Display: in.PackagingPhotos.Ingredients.Display,
			},
		},
		Allergens: in.Allergens,
		BrandList: in.BrandList,
		Countries: in.Countries,
		CountryDetails: &chompv1beta1.CountryDetails{
			EnglishSpeaking:    int32(in.CountryDetails.EnglishSpeaking),
			NonEnglishSpeaking: int32(in.CountryDetails.NonEnglishSpeaking),
		},
		PalmOilIngredients:    in.PalmOilIngredients,
		IngredientList:        in.IngredientList,
		HasEnglishIngredients: in.HasEnglishIngredients,
		Minerals:              in.Minerals,
		Traces:                in.Traces,
		Vitamins:              in.Vitamins,
		Description:           in.Description,
		Keywords:              in.Keywords,
	}
}

type ChompResponse struct {
	Items []ChompFoodItem `json:"items"`
}

type ChompFoodItem struct {
	Barcode     string `json:"barcode"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Ingredients string `json:"ingredients"`
	Package     struct {
		Quantity int    `json:"quantity"`
		Size     string `json:"size"`
	} `json:"package"`
	Serving struct {
		Size            string `json:"size"`
		MeasurementUnit string `json:"measurement_unit"`
		SizeFulltext    string `json:"size_fulltext"`
	} `json:"serving"`
	Categories []string `json:"categories"`
	Nutrients  []struct {
		Name            string `json:"name"`
		Per100G         int    `json:"per_100g"`
		MeasurementUnit string `json:"measurement_unit"`
		Rank            int    `json:"rank"`
		DataPoints      int    `json:"data_points"`
		Description     string `json:"description"`
	} `json:"nutrients"`
	DietLabels struct {
		Vegan struct {
			Name                  string `json:"name"`
			IsCompatible          bool   `json:"is_compatible"`
			CompatibilityLevel    int    `json:"compatibility_level"`
			Confidence            int    `json:"confidence"`
			ConfidenceDescription string `json:"confidence_description"`
		} `json:"vegan"`
		Vegetarian struct {
			Name                  string `json:"name"`
			IsCompatible          bool   `json:"is_compatible"`
			CompatibilityLevel    int    `json:"compatibility_level"`
			Confidence            int    `json:"confidence"`
			ConfidenceDescription string `json:"confidence_description"`
		} `json:"vegetarian"`
		GlutenFree struct {
			Name                  string `json:"name"`
			IsCompatible          bool   `json:"is_compatible"`
			CompatibilityLevel    int    `json:"compatibility_level"`
			Confidence            int    `json:"confidence"`
			ConfidenceDescription string `json:"confidence_description"`
		} `json:"gluten_free"`
	} `json:"diet_labels"`
	DietFlags []struct {
		Ingredient               string `json:"ingredient"`
		IngredientDescription    string `json:"ingredient_description"`
		DietLabel                string `json:"diet_label"`
		IsCompatible             string `json:"is_compatible"`
		CompatibilityLevel       int    `json:"compatibility_level"`
		CompatibilityDescription string `json:"compatibility_description"`
		IsAllergen               bool   `json:"is_allergen"`
	} `json:"diet_flags"`
	PackagingPhotos struct {
		Front struct {
			Small   string `json:"small"`
			Thumb   string `json:"thumb"`
			Display string `json:"display"`
		} `json:"front"`
		Nutrition struct {
			Small   string `json:"small"`
			Thumb   string `json:"thumb"`
			Display string `json:"display"`
		} `json:"nutrition"`
		Ingredients struct {
			Small   string `json:"small"`
			Thumb   string `json:"thumb"`
			Display string `json:"display"`
		} `json:"ingredients"`
	} `json:"packaging_photos"`
	Allergens      []string `json:"allergens"`
	BrandList      []string `json:"brand_list"`
	Countries      []string `json:"countries"`
	CountryDetails struct {
		EnglishSpeaking    int `json:"english_speaking"`
		NonEnglishSpeaking int `json:"non_english_speaking"`
	} `json:"country_details"`
	PalmOilIngredients    []string `json:"palm_oil_ingredients"`
	IngredientList        []string `json:"ingredient_list"`
	HasEnglishIngredients bool     `json:"has_english_ingredients"`
	Minerals              []string `json:"minerals"`
	Traces                []string `json:"traces"`
	Vitamins              []string `json:"vitamins"`
	Description           string   `json:"description"`
	Keywords              []string `json:"keywords"`
}
