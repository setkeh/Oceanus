package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	guuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

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

	ret := models.photo{
		Src: file.Filename,
		ID:  guid,
		URL: fmt.Sprintf("%s/%s", url, guid),
	}

	i := models.image{
		ID:  ret.ID,
		Src: ret.Src,
		URL: ret.URL,
		B64: b64,
	}

	db.insertImage(i)

	return c.JSON(http.StatusOK, ret)
}

func getImageHandler(c echo.Context) error {
	id := c.QueryParam("ID")

	ret, err := getImage(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	//fmt.Println(ret)

	var img image
	jsonerr := json.Unmarshal(ret, &img)
	if err != nil {
		return c.String(http.StatusInternalServerError, jsonerr.Error())
	}

	i := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img.B64))

	return c.Stream(http.StatusOK, "image/png", i)
}

func getImageListHandler(c echo.Context) error {
	ret, err := getImageList()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONPretty(http.StatusOK, ret, "\n")
}
