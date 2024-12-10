package types

import "fmt"

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("[ERROR:%d] %s", e.Code, e.Msg)
}

type Song struct {
	ID          int    `form:"id"`
	Name        string `form:"song"`
	GroupName   string `form:"group"`
	ReleaseDate string `form:"releaseDate"`
	Text        string `form:"text"`
	Link        string `form:"link"`
	Page        int    `form:"page"`
}
