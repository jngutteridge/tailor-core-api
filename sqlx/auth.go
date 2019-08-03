package sqlx

import (
	. "jng.dev/tailor/core/core"
)


func GetUser(u *User, email string) bool {
	db := Database()
	defer Close(db)
	err := db.Get(u, "SELECT id, name, email, image_url, verification_code FROM core.user WHERE email = $1", email)
	return err == nil && u.ID != ""
}