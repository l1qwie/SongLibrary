package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/SongLibrary/app/logs"
	"github.com/l1qwie/SongLibrary/app/types"
)

func code500() error {
	err := new(types.Err)
	err.Code = http.StatusBadRequest
	err.Msg = "invalid query parameters"
	return err
}

func firstStringLook(song *types.Song, ctx *gin.Context) error {
	var err error
	if err = ctx.ShouldBindQuery(song); err != nil {
		err = code500()
	}
	return err
}

func firstJsonLook(song *types.Song, ctx *gin.Context) error {
	var err error
	var body []byte
	if body, err = io.ReadAll(ctx.Request.Body); err == nil {
		err = json.Unmarshal(body, song)
	}
	return err
}

func isDefault(song *types.Song) error {
	var err error
	if (song.ID == 0) && (song.Name == "") &&
		(song.GroupName == "") && (song.ReleaseDate == "") &&
		(song.Text == "") && (song.Link == "") && (song.Page == 0) {
		logs.Nothing()
		err = code500()
	}
	return err
}

func nothingElse(song *types.Song) bool {
	var res bool
	if (song.Name == "") &&
		(song.GroupName == "") && (song.ReleaseDate == "") &&
		(song.Text == "") && (song.Link == "") && (song.Page == 0) {
		logs.Nothing()
		res = true
	}
	return res
}
