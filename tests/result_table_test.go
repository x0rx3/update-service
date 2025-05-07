package tests

import (
	"update-service/pkg/models"

	"github.com/google/uuid"
)

type ResultTableTest struct {
}

func NewResultTableTest() *ResultTableTest {
	return &ResultTableTest{}
}
func (inst *ResultTableTest) Insert(result *models.Result) (string, error) {
	return uuid.NewString(), nil
}
