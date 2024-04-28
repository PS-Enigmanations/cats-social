package helper

import (
	"enigmanations/cats-social/model/domain"
	"enigmanations/cats-social/model/web"
)

func ToCatResponse(cat domain.Cat) web.CatResponse {
	return web.CatResponse{
		Id:   cat.Id,
		Name: cat.Name,
	}
}
