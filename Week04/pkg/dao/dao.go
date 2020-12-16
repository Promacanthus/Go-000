package dao

import (
	"database/sql"

	xerrors "github.com/pkg/errors"
)

type Repo interface {
	SaveStringLength(s string) error
	GetAmount(s string) (int, error)
}

type repo struct {
	db *sql.DB
}

var _ Repo = (*repo)(nil)

func NewRepo(db *sql.DB) Repo {
	return &repo{db: db}
}

func (r *repo) SaveStringLength(s string) error {
	_, err := r.db.Query("insert into Strings(name) values (?)", s)
	return xerrors.Wrapf(err, "failed querying database")
}

func (r *repo) GetAmount(s string) (amount int, err error) {
	rows, err := r.db.Query("select Length from Strings where name = ?", s)
	if err != nil {
		return 0, xerrors.Wrapf(err, "failed querying database")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&amount); err != nil {
			return 0, xerrors.Wrapf(err, "failed scanning rows")
		}
	}

	return amount, xerrors.Wrapf(rows.Err(), "rows error")
}
