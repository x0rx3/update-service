package tests

import (
	"update-service/internal/model"

	"github.com/google/uuid"
)

type ResultTableTest struct {
}

func NewResultTableTest() *ResultTableTest {
	return &ResultTableTest{}
}
func (inst *ResultTableTest) Insert(result *model.Result) (string, error) {
	return uuid.NewString(), nil
}
