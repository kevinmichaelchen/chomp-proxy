package service

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/sirupsen/logrus"
	chompv1beta1 "go.buf.build/bufbuild/connect-go/kevinmichaelchen/chompapis/chomp/v1beta1"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetFood(
	ctx context.Context,
	req *connect.Request[chompv1beta1.GetFoodRequest],
) (*connect.Response[chompv1beta1.GetFoodResponse], error) {
	res := &chompv1beta1.GetFoodResponse{
		Food: &chompv1beta1.Food{
			Name: "Hard-coded Test Food",
		},
	}
	logrus.WithField("barcode", req.Msg.GetCode()).Info("Retrieving food...")
	out := connect.NewResponse(res)
	out.Header().Set("API-Version", "v1beta1")
	return out, nil
}

func (s *Service) ListFoods(
	ctx context.Context,
	req *connect.Request[chompv1beta1.ListFoodsRequest],
) (*connect.Response[chompv1beta1.ListFoodsResponse], error) {
	res := &chompv1beta1.ListFoodsResponse{
		Items: []*chompv1beta1.Food{
			{
				Name: "Hard-coded Test Food",
			},
		},
	}
	logrus.WithField("query", req.Msg.GetName()).Info("Retrieving food...")
	out := connect.NewResponse(res)
	out.Header().Set("API-Version", "v1beta1")
	return out, nil
}
