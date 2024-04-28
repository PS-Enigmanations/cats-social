package repository

import (
	"context"
	"database/sql"
	"enigmanations/cats-social/model/domain"
)

type CatRepositoryImpl struct {
}

func NewCatRepository() CatRepository {
	return &CatRepositoryImpl{}
}

func (c *CatRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, cat domain.Cat) domain.Cat {
	SQL := "INSERT into cats (name, race, sex, age_in_month, description, image_urls) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, cat.Name, cat.Race, cat.AgeInMonth, cat.Description, cat.ImageUrls)
	if err != nil {
		panic(err)
	}
	catId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	cat.Id = int(catId)
	return cat
}

func (c *CatRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, cat domain.Cat) domain.Cat {
	return domain.Cat{}
}

func (c *CatRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, cat domain.Cat) {

}

func (c *CatRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, catId int) domain.Cat {
	return domain.Cat{}
}

func (c *CatRepositoryImpl) Get(ctx context.Context, tx *sql.Tx) []domain.Cat {
	return []domain.Cat{}
}
