package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	guuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type photo struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
}

// Handler
func postImageHandler(c echo.Context) error {
	//fmt.Println(c.FormFile("file"))

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	buf := make([]byte, file.Size)
	io.ReadFull(src, buf)

	b64 := base64.StdEncoding.EncodeToString(buf)

	// Debug Base64 Output
	//fmt.Println(b64)
	//ioutil.WriteFile("b64.txt", []byte(b64), 0777)

	guid := guuid.New().String()

	ret := photo{
		Src: file.Filename,
		ID:  guid,
		URL: fmt.Sprintf("%s/%s", url, guid),
	}

	i := image{
		ID:  ret.ID,
		Src: ret.Src,
		URL: ret.URL,
		B64: b64,
	}

	insertImage(i)

	return c.JSON(http.StatusOK, ret)
}

func getImageHandler(c echo.Context) error {
	id := c.QueryParam("ID")

	ret, err := getImage(id)
	if err != nil {
		c.Logger().Error(err)
	}
	fmt.Println(ret)

	var img image
	json.Unmarshal(ret, &img)

	return c.String(http.StatusOK, img.Src)
}
