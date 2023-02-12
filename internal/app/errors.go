package app

type Error struct {
	Message string `json:"message"`
}

var (
	invalidPostId = Error{Message: "invalid post ID."}
)
