package handlers

import (
	"bytes"
	"fmt"
	"image"
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
	DB     ImageStore
	obst   = os.Getenv("OBJECT_STORAGE_URL")
	betype = os.Getenv("BACKEND_TYPE")
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

	var url string
	if betype == "LOCAL" {
		url = fmt.Sprintf("/img/%s", file.Filename)
	}

	if betype == "BUCKET" {
		url = fmt.Sprintf("https://%s/%s/%s", b.Endpoint, b.BucketName, file.Filename)
	}

	ret := models.Photo{
		Src: file.Filename,
		ID:  guid,
		URL: url, //fmt.Sprintf("https://%s/%s/%s", b.Endpoint, b.BucketName, file.Filename),
	}

	i := models.Image{
		ID:    ret.ID,
		Src:   ret.Src,
		URL:   ret.URL,
		Image: im,
	}

	switch betype {
	case "LOCAL":
		//buf := make([]byte, file.Size)
		//io.ReadFull(src, buf)

		//b64 := base64.StdEncoding.EncodeToString(buf)

		// Debug Base64 Output
		//fmt.Println(b64)
		//ioutil.WriteFile("b64.txt", []byte(b64), 0777)
		save, err := os.Create(ret.URL)
		if err != nil {
			c.Logger().Errorf("Error Creating Image on Disk: %s", err)
		}
		defer save.Close()

		/*_, serr := io.Copy(save, src)
		if serr != nil {
			c.Logger().Errorf("Error Saving Image on Disk: %s", serr)
		}*/

		png.Encode(save, im)
	case "BUCKET":
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
	}

	r, err := d.Insert(i.ID, i.URL, i.Src)
	if err != nil {
		c.Logger().Infof("Bytes Uploaded: %d\n", r)
		c.Logger().Errorf("DB Insert Error: %s\n", err)
		return err
	}
	c.Logger().Infof("Bytes Uploaded: %d\n", r)

	return c.JSON(http.StatusOK, ret)
}

func GetImageHandler(c echo.Context) error {

	id := c.QueryParam("ID")

	d := new(db.Mysql)
	d.Init()

	path, err := d.GetByID(id)
	if err != nil {
		c.Logger().Errorf("Error fetching path from DB: %s", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	c.Logger().Debug(path)

	if betype == "LOCAL" {
		c.Logger().Infof("Image Path: %s", path)
		file, ferr := os.Open(path)
		if ferr != nil {
			c.Logger().Errorf("Error Opening Image Path From Disk: %s", ferr)
			return c.JSON(http.StatusInternalServerError, ferr)
		}

		defer file.Close()

		/*fStat, _ := file.Stat()
		buf := make([]byte, fStat.Size())*/
		buf1 := new(bytes.Buffer)
		/*io.ReadFull(file, buf)

		b64 := base64.StdEncoding.EncodeToString(buf)

		i := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64))*/
		/*im, err := png.Decode(file)
		if err != nil {
			c.Logger().Errorf("PNG Decode Error: %s\n", err)
			return err
		}*/

		file.Seek(0, 0)
		src, _, ierr := image.Decode(file)
		if ierr != nil {
			c.Logger().Errorf("Image Decode Failed: %s", ierr)
		}

		png.Encode(buf1, src)

		//return c.Stream(http.StatusOK, "image/png", i)
		return c.Blob(http.StatusOK, "image/png", buf1.Bytes())
	}

	if betype == "BUCKET" {
		return c.JSON(http.StatusOK, path)
	}

	/*
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
	*/

	return c.JSON(http.StatusInternalServerError, "Somthing has Gone wrong the code should never make it here This Likely Means you forgot to set the BACKEND_TYPE environment variable.")
}

/*
func GetImageListHandler(c echo.Context) error {
	ret, err := DB.ImageList()

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSONPretty(http.StatusOK, ret, "\n")
}
*/
