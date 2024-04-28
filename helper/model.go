package helper

import (
	"PS-Enigmanations/cats-social/model/domain"
	"PS-Enigmanations/cats-social/model/web"
)

func ToCatResponse(cat domain.Cat) web.CatResponse {
	return web.CatResponse{
		Id:   cat.Id,
		Name: cat.Name,
	}
}
