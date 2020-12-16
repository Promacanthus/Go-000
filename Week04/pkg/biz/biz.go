package biz

import (
	"errors"

	"github.com/Promacanthus/Go-000/Week04/pkg/dao"
	xerrors "github.com/pkg/errors"
)

var (
	// ErrEmpty is returned when input string is empty
	ErrEmpty = errors.New("empty string")
)

type StringService interface {
	Save(string) (int, error)
	GetAmount(s string) (int, error)
}

type stringService struct {
	dao.Repo
}

var _ StringService = (*stringService)(nil)

func NewStringService(repo dao.Repo) StringService {
	return &stringService{repo}
}

func (ss *stringService) Save(s string) (int, error) {
	if s == "" {
		return 0, ErrEmpty
	}

	return len(s), xerrors.Wrapf(ss.Repo.SaveStringLength(s), "failed saving string")
}

func (ss *stringService) GetAmount(s string) (int, error) {
	return ss.Repo.GetAmount(s)
}
