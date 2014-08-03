package migrate

import (
	"appengine"
	"appengine/blobstore"
	"appengine/delay"
	aeimage "appengine/image"
	"appengine/urlfetch"
	"bytes"
	"encoding/json"
	"github.com/danielphan/ae/logger"
	"github.com/danielphan/dandubois/painting"
	"github.com/danielphan/dandubois/painting/legacy"
	"github.com/danielphan/ae/object"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strconv"
)

const paintingsURL = baseURL + "paintings.json"

type paintings []*legacy.Painting

func fetchPaintings(c *http.Client) (paintings, error) {
	res, err := c.Get(paintingsURL)
	if err != nil {
		return nil, err
	}
	return newPaintings(res.Body)
}

func newPaintings(r io.Reader) (paintings, error) {
	var ps paintings
	err := json.NewDecoder(r).Decode(&ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (p paintings) convert(w http.ResponseWriter, c appengine.Context, cs categories, ms media) {
	for _, o := range p {
		old := o.Fields
		p := &painting.Painting{
			Object: object.New(
				painting.Kind,
				strconv.Itoa(old.Number),
			),
			Title: old.Title,
			Image: painting.Image{
				BlobKey: appengine.BlobKey(" "),
			},
			Description: old.Description,
			Year:        old.Year,
			Categories:  cs.convert(old.Categories),
			Media:       ms.convert(old.Medium),
			Price:       old.Price,
			Width:       old.Width,
			Height:      old.Height,
			ForSale:     old.Status == "F",
		}
		if old.Image != "" {
			p.Image.URL = baseURL + old.Image
		}
		saveLater.Call(c, p)
	}
}

var saveLater *delay.Function

func init() {
	saveLater = delay.Func("save", save)
}

func save(c appengine.Context, p *painting.Painting) error {
	if p.Image.URL != "" {
		// Fetch the image.
		res, err := urlfetch.Client(c).Get(p.Image.URL)
		if err != nil {
			logger.Error(c, err)
			return err
		}
		defer res.Body.Close()

		// Save what we read to decode the image config,
		// so we can save the whole image to the blobstore.
		buf := &bytes.Buffer{}
		t := io.TeeReader(res.Body, buf)

		// Decode the config to get the size and content type.
		conf, ext, err := image.DecodeConfig(t)
		if err != nil {
			logger.Error(c, err)
			return err
		}
		p.Image.Width = conf.Width
		p.Image.Height = conf.Height

		// Create a new blob.
		b, err := blobstore.Create(c, "image/"+ext)
		if err != nil {
			logger.Error(c, err)
			return err
		}

		// Copy the image into it.
		// Prepend what we read to decode the image config.
		r := io.MultiReader(buf, res.Body)

		_, err = io.Copy(b, r)
		if err != nil {
			logger.Error(c, err)
			return err
		}
		err = b.Close()
		if err != nil {
			logger.Error(c, err)
			return err
		}

		// Add the image to the painting.
		p.Image.BlobKey, err = b.Key()
		if err != nil {
			logger.Error(c, err)
			return err
		}
		u, err := aeimage.ServingURL(c, p.Image.BlobKey, nil)
		if err != nil {
			logger.Error(c, err)
			return err
		}
		p.Image.URL = u.String()
	}
	return p.Save(c)
}
