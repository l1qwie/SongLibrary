package integrated

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/l1qwie/SongLibrary/api"
	"github.com/l1qwie/SongLibrary/app/logs"
	"github.com/l1qwie/SongLibrary/app/types"
)

type response struct {
	status       int
	havingAnswer bool
	msg          interface{}
	structer     interface{}
}

type httpTestReq struct {
	method            string
	havingContentType bool
	path              []string
	body              []io.Reader
	response          []response
}

func checkResponse(respSer *http.Response, respExpected response, t *testing.T) {
	defer respSer.Body.Close()
	if respSer.StatusCode != respExpected.status {
		t.Fatal("unexpected status of the request")
	}
	respbody, err := io.ReadAll(respSer.Body)
	if err != nil {
		t.Fatal(err)
	}
	if respExpected.havingAnswer {
		err = json.Unmarshal([]byte(respbody), respExpected.structer)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(respExpected.structer, respExpected.msg) {
			t.Fatal("reponse isn't the same as expected")
		}
	}
}

func makeReq(obj httpTestReq, i int, t *testing.T) (*http.Response, error) {
	req, err := http.NewRequest(obj.method, obj.path[i], obj.body[i])
	if obj.havingContentType {
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	return client.Do(req)
}

func doer(container httpTestReq, t *testing.T) {
	for i := range container.path {
		if resp, err := makeReq(container, i, t); err != nil {
			t.Fatal(err)
		} else {
			checkResponse(resp, container.response[i], t)
		}
	}
}

func (htr *httpTestReq) statusGetSong() {
	htr.method = "GET"
	htr.body = []io.Reader{nil, nil, nil}
	htr.path = []string{"http://localhost:8080/song?id=1&name=Wouldn't%20It%20Be%20Nice", "http://localhost:8080/song?",
		"http://localhost:8080/song?id=1&chtoto=lichnee"}
	htr.response = []response{{200, false, nil, nil}, {400, false, nil, nil}, {200, false, nil, nil}}
}

func (htr *httpTestReq) statusGetCouplet() {
	htr.method = "GET"
	htr.body = []io.Reader{nil, nil, nil, nil}
	htr.path = []string{"http://localhost:8080/couplet?id=1&text=Wouldn't%20It%20Be%20Nice", "http://localhost:8080/couplet?",
		"http://localhost:8080/couplet?id=1&chtoto=lichnee", "http://localhost:8080/couplet?text=text"}
	htr.response = []response{{200, false, nil, nil}, {400, false, nil, nil}, {400, false, nil, nil}, {200, false, nil, nil}}
}

func (htr *httpTestReq) statusDeleteSong() {
	htr.method = "DELETE"
	htr.body = []io.Reader{nil, nil, nil, nil, nil}
	htr.path = []string{"http://localhost:8080/song?id=1&text=Wouldn't%20It%20Be%20Nice", "http://localhost:8080/song?",
		"http://localhost:8080/song?group=1&chtoto=lichnee", "http://localhost:8080/song?text=text", "http://localhost:8080/song?id=222"}
	htr.response = []response{{200, false, nil, nil}, {400, false, nil, nil}, {400, false, nil, nil}, {400, false, nil, nil}, {200, false, nil, nil}}
}

func getRequestData(songcontainer []types.Song, t *testing.T) []io.Reader {
	res := make([]io.Reader, len(songcontainer))
	for i, song := range songcontainer {
		body, err := json.Marshal(song)
		if err != nil {
			t.Fatal(err)
		}
		res[i] = bytes.NewBuffer(body)
	}
	return res
}

func (htr *httpTestReq) statusChangeSong(t *testing.T) {
	htr.method = "PUT"
	htr.havingContentType = true
	htr.body = getRequestData([]types.Song{{ID: 11, Name: "Holy Molly!"}, {},
		{Name: "Wouldn't It Be Nice"}, {ID: 512, GroupName: "The Beach Boys", Name: "Wouldn't It Be Nice"}}, t)
	htr.path = []string{"http://localhost:8080/song", "http://localhost:8080/song", "http://localhost:8080/song", "http://localhost:8080/song"}
	htr.response = []response{{200, false, nil, nil}, {400, false, nil, nil}, {400, false, nil, nil}, {200, false, nil, nil}}
}

func (htr *httpTestReq) statusNewSong(t *testing.T) {
	htr.method = "POST"
	htr.havingContentType = true
	htr.body = getRequestData([]types.Song{{ID: 11}, {},
		{Name: "Wouldn't It Be Nice"}, {ID: 512, GroupName: "The Beach Boys", Name: "Wouldn't It Be Nice"}}, t)
	htr.path = []string{"http://localhost:8080/song", "http://localhost:8080/song", "http://localhost:8080/song", "http://localhost:8080/song"}
	htr.response = []response{{400, false, nil, nil}, {400, false, nil, nil}, {400, false, nil, nil}, {200, false, nil, nil}}
}

func TestStatusGetSong(t *testing.T) {
	htr := new(httpTestReq)
	htr.statusGetSong()
	doer(*htr, t)
}

func TestStatusGetCouplet(t *testing.T) {
	htr := new(httpTestReq)
	htr.statusGetCouplet()
	doer(*htr, t)
}

func TestStatusDeleteSong(t *testing.T) {
	htr := new(httpTestReq)
	htr.statusDeleteSong()
	doer(*htr, t)
}

func TestStatusChangeSong(t *testing.T) {
	htr := new(httpTestReq)
	htr.statusChangeSong(t)
	doer(*htr, t)
}

func TestStatusNewSong(t *testing.T) {
	htr := new(httpTestReq)
	htr.statusNewSong(t)
	doer(*htr, t)
}

func init() {
	logs.SetDebug()
	go api.StartAPI()
	time.Sleep(100 * time.Millisecond)
}
