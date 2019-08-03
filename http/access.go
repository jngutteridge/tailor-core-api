package http

import (
	. "jng.dev/tailor/core/core"
)

var _ AccessService = &AccessHTTPService{}

type AccessHTTPService struct{}

func (s *AccessHTTPService) Access(u *User) (bool, Access) {

}
