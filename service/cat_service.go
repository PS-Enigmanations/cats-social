package service

import (
	"context"
	"enigmanations/cats-social/model/web"
)

type CatService interface {
	Create(ctx context.Context, request web.CatCreateRequest) web.CatResponse
	Update()
	Delete()
	FindById()
	Get()
}
