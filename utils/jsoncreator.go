package utils

import (
	"github.com/brianvoe/gofakeit/v6"
	"l0/internal/datastorage/structure"
	"log"
)

func FakeModel() (error, structure.Model) {
	model := structure.Model{}
	err := gofakeit.Struct(&model)
	if err != nil {
		log.Fatalf("GOFAKEIT ERROR: %s", err)
	}
	return err, model
}
