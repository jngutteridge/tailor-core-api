package sqlx

import (
	"github.com/jmoiron/sqlx"
	. "jng.dev/tailor/core/core"
)

var _ AccessService = &AccessSQLXService{}

type AccessSQLXService struct {
	sqlx.DB
}

func NewAccessService() *AccessService {
	s := Database()
	return AccessSQLXService{
		DB: s
	}
}

func (s *AccessSQLXService) Access(u *User) (bool, Access) {

}