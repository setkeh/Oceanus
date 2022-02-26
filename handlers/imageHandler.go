package handlers

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"os"

	guuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/setkeh/Oceanus/bucket"
	"github.com/setkeh/Oceanus/db"
	"github.com/setkeh/Oceanus/models"
)

type ImageStore interface {
	InsertImage(models.Image)
	Image(id string) ([]byte, error)
	ImageList() ([]models.Photo, error)
}

var (
	DB   ImageStore
	obst = os.Getenv("OBJECT_STORAGE_URL")
)

func PostImageHandler(c echo.Context) error {
	//fmt.Println(c.FormFile("file"))

	d := new(db.Mysql)
	d.Init()

	b := new(bucket.Bucket)
	b.Init()

	file, err := c.FormFile("file")
	if err != nil {
		c.Logger().Errorf("c.Formfile Error: %s\n", err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("file.Open error: %s\n", err)
		return err
	}

	defer src.Close()
	im, err := png.Decode(src)
	if err != nil {
		c.Logger().Errorf("PNG Decode Error: %s\n", err)
		return err
	}

	//buf := make([]byte, file.Size)
	//io.ReadFull(src, buf)

	//b64 := base64.StdEncoding.EncodeToString(buf)

	// Debug Base64 Output
	//fmt.Println(b64)
	//ioutil.WriteFile("b64.txt", []byte(b64), 0777)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, im)
	if err != nil {
		return err
	}
	imgr := bytes.NewReader(buf.Bytes())

	guid := guuid.New().String()

	ret := models.Photo{
		Src: file.Filename,
		ID:  guid,
		URL: fmt.Sprintf("https://%s/%s/%s", b.Endpoint, b.BucketName, file.Filename),
	}

	i := models.Image{
		ID:    ret.ID,
		Src:   ret.Src,
		URL:   ret.URL,
		Image: im,
	}

	e := b.CheckBucketExists()
	if e != true {
		b.CreateBucket()
	}

	rb, err := b.UploadImage(i.Src, imgr, imgr.Size())
	if err != nil {
		c.Logger().Infof("Bytes Uploaded: %d\n", rb)
		c.Logger().Errorf("UploadImage Error: %s\n", err)
		return err
	}
	c.Logger().Infof("Bytes Uploaded: %d\n", rb)

	r, err := d.Insert(i.ID, i.URL, i.Src)
	if err != nil {
		c.Logger().Infof("Bytes Uploaded: %d\n", r)
		c.Logger().Errorf("DB Insert Error: %s\n", err)
		return err
	}
	c.Logger().Infof("Bytes Uploaded: %d\n", r)

	return c.JSON(http.StatusOK, ret)
}

/*
func GetImageHandler(c echo.Context) error {

	id := c.QueryParam("ID")

	ret, err := DB.Image(id)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	//fmt.Println(ret)

	var img models.Image
	jsonerr := json.Unmarshal(ret, &img)
	if err != nil {
		return c.String(http.StatusInternalServerError, jsonerr.Error())
	}

	i := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img.B64))

	return c.Stream(http.StatusOK, "image/png", i)
}

func GetImageListHandler(c echo.Context) error {
	ret, err := DB.ImageList()

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONPretty(http.StatusOK, ret, "\n")
}
*/
