package tests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
	"update-service/internal/model"
	"update-service/internal/utils"
)

type UpdateServerClientTest struct{}

func NewUpdateServerClientTest() *UpdateServerClientTest {
	return &UpdateServerClientTest{}
}

func (inst *UpdateServerClientTest) Login() error {
	return nil
}

func (inst *UpdateServerClientTest) UpdateList(pkgType utils.PackageType) ([]model.RrUpdates, error) {
	switch pkgType {
	case utils.Rules:
		return []model.RrUpdates{
			model.RrUpdates{
				Name:   "Old-Rules.3.9.gz",
				Hw:     []string{SuccessAllVersion},
				Sw:     []string{SuccessAllVersion},
				Latest: false,
				Link:   "123",
				Size:   1,
				Date:   time.Now().Add(6 * time.Hour),
			},
			model.RrUpdates{
				Name:   "Latest-Rules.3.9.gz",
				Hw:     []string{SuccessAllVersion},
				Sw:     []string{SuccessAllVersion},
				Latest: true,
				Link:   "123",
				Size:   1,
				Date:   time.Now(),
			},
			model.RrUpdates{
				Name:   "NotNew-Rules.3.10.gz",
				Hw:     []string{"3.8", "3.9", "3.10"},
				Sw:     []string{"3.8", "3.9", "3.10"},
				Latest: false,
				Link:   "123",
				Size:   1,
				Date:   time.Now().Add(6 * time.Hour),
			},
			model.RrUpdates{
				Name:   "New-Rules.3.10.gz",
				Hw:     []string{"3.8", "3.9", "3.10"},
				Sw:     []string{"3.8", "3.9", "3.10"},
				Latest: true,
				Link:   "123",
				Size:   1,
				Date:   time.Now(),
			},
			model.RrUpdates{
				Name:   "NotNew-Rules.3.11.gz",
				Hw:     []string{"3.9", "3.10", "3.11"},
				Sw:     []string{"3.9", "3.10", "3.11"},
				Latest: false,
				Link:   "123",
				Size:   1,
				Date:   time.Now().Add(6 * time.Hour),
			},
			model.RrUpdates{
				Name:   "New-Rules.3.11.gz",
				Hw:     []string{"3.9", "3.10", "3.11"},
				Sw:     []string{"3.9", "3.10", "3.11"},
				Latest: true,
				Link:   "123",
				Size:   1,
				Date:   time.Now(),
			},
			model.RrUpdates{
				Name:   "NotNew-Rules.3.12.gz",
				Hw:     []string{"3.10", "3.11", "3.12"},
				Sw:     []string{"3.10", "3.11", "3.12"},
				Latest: false,
				Link:   "123",
				Size:   1,
				Date:   time.Now().Add(6 * time.Hour),
			},
			model.RrUpdates{
				Name:   "New-Rules.3.12.gz",
				Hw:     []string{"3.10", "3.11", "3.12"},
				Sw:     []string{"3.10", "3.11", "3.12"},
				Latest: true,
				Link:   "123",
				Size:   1,
				Date:   time.Now(),
			},
		}, nil
	}

	return []model.RrUpdates{
		model.RrUpdates{
			Name:   "Malware.gz",
			Hw:     []string{"3.10", "3.11", "3.12", "3.9"},
			Sw:     []string{"3.10", "3.11", "3.12", "3.9"},
			Latest: false,
			Link:   "123",
			Size:   1,
		},
		model.RrUpdates{
			Name:   "Malware.gz",
			Hw:     []string{"3.10", "3.11", "3.12", "3.9"},
			Sw:     []string{"3.10", "3.11", "3.12", "3.9"},
			Latest: true,
			Link:   "123",
			Size:   1,
		},
	}, nil
}

func (inst *UpdateServerClientTest) Download(pkgType utils.PackageType, pkgInfo *model.RrUpdates, dir4Save string) (string, error) {
	cache, err := os.Create(fmt.Sprintf("%v/%v", dir4Save, pkgInfo.Name))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(cache, bytes.NewBuffer([]byte(fmt.Sprintf("new file: %s", pkgInfo.Name))))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v/%v", dir4Save, pkgInfo.Name), nil
}
