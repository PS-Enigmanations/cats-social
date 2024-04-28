package service

import (
	"PS-Enigmanations/cats-social/model/web"
	"context"
)

type CatService interface {
	Create(ctx context.Context, request web.CatCreateRequest) web.CatResponse
	Update()
	Delete()
	FindById()
	Get()
}
