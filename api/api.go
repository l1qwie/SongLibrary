package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/SongLibrary/app"
	"github.com/l1qwie/SongLibrary/app/logs"
	"github.com/l1qwie/SongLibrary/app/types"
)

type server struct {
	router *gin.Engine
}

func configuration() *server {
	s := new(server)
	s.router = gin.Default()
	return s
}

func majorLogic(ctx *gin.Context, firstLook func(*types.Song, *gin.Context) error, entryF func(*types.Song) ([]byte, error)) {
	var err error
	var response []byte
	query := new(types.Song)
	if err = firstLook(query, ctx); err == nil {
		if err = isDefault(query); err == nil {
			response, err = entryF(query)
		}
	}
	sendResponse(ctx, response, err)
}

func sendResponse(ctx *gin.Context, response []byte, err error) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	} else {
		ctx.JSON(http.StatusOK, response)
	}
}

func getSongLogic(ctx *gin.Context) {
	majorLogic(ctx, firstStringLook, func(s *types.Song) ([]byte, error) {
		logs.InputDataIsOK()
		return app.GetSong(s)
	})
}

func getCoupletLogic(ctx *gin.Context) {
	majorLogic(ctx, firstStringLook, func(s *types.Song) ([]byte, error) {
		var err error
		var response []byte
		if s.Text == "" {
			logs.FieldRequired("text")
			err = code500()
		} else {
			logs.InputDataIsOK()
			response, err = app.GetCouple(s)
		}
		return response, err
	})
}

func deleteSongLogic(ctx *gin.Context) {
	majorLogic(ctx, firstStringLook, func(s *types.Song) ([]byte, error) {
		var err error
		var response []byte
		if s.ID == 0 {
			logs.FieldRequired("id")
			err = code500()
		} else {
			logs.InputDataIsOK()
			response, err = app.DeleteSong(s)
		}
		return response, err
	})
}

func changeSongLogic(ctx *gin.Context) {
	majorLogic(ctx, firstJsonLook, func(s *types.Song) ([]byte, error) {
		var err error
		var response []byte
		if s.ID == 0 || nothingElse(s) {
			logs.FieldsRequired("id", "something else")
			err = code500()
		} else {
			logs.InputDataIsOK()
			response, err = app.ChangeSong(s)
		}
		return response, err
	})
}

func newSongLogic(ctx *gin.Context) {
	majorLogic(ctx, firstJsonLook, func(s *types.Song) ([]byte, error) {
		var err error
		var response []byte
		if (s.Name == "") || (s.GroupName == "") {
			logs.FieldsRequired("name", "group")
			err = code500()
		} else {
			logs.InputDataIsOK()
			response, err = app.NewSong(s)
		}
		return response, err
	})
}

func (s *server) getSong() {
	path := "/song"
	s.router.GET(path, getSongLogic)
	logs.StartPoint(path, "GET")
}

func (s *server) getCouplet() {
	path := "/couplet"
	s.router.GET(path, getCoupletLogic)
	logs.StartPoint(path, "GET")
}

func (s *server) deleteSong() {
	path := "/song"
	s.router.DELETE(path, deleteSongLogic)
	logs.StartPoint(path, "DELETE")
}

func (s *server) changeSong() {
	path := "/song"
	s.router.PUT(path, changeSongLogic)
	logs.StartPoint(path, "PUT")
}

func (s *server) newSong() {
	path := "/song"
	s.router.POST(path, newSongLogic)
	logs.StartPoint(path, "POST")
}

func StartAPI() {
	s := configuration()
	s.getSong()
	s.getCouplet()
	s.deleteSong()
	s.changeSong()
	s.newSong()

	s.router.Run(":8080")
}
