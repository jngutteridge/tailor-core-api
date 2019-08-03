package http

import (
	. "jng.dev/tailor/core/core"
	"jng.dev/tailor/core/sqlx"
	"net/http"
)

func TextHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlx.Database()
	defer sqlx.Close(db)
	var t []Text
	err := db.Select(&t, "SELECT slug, text FROM core.text")
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, ErrorMessage{"Error retrieving text"})
		return
	}
	WriteResponse(w, http.StatusOK, TextResponse{t})
}
