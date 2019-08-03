package core

import "jng.dev/tailor/core/core"

type PropertyID string
type ModuleID string

type Property struct {
	ID       PropertyID `db:"id" json:"id"`
	Slug     string     `db:"slug" json:"slug"`
	Name     string     `db:"name" json:"name"`
	ImageURL string     `db:"image_url" json:"imageUrl"`
	LiveURL  string     `db:"live_url" json:"liveUrl"`
	Modules  []Module   `json:"modules"`
}

type Module struct {
	Slug string `db:"slug" json:"slug"`
	Name string `db:"name" json:"name"`
}

type Access struct {
	Properties []Property `json:"properties"`
}

type AccessService interface {
	Access(user *core.User) (bool, Access)
}
