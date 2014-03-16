package migrate

import (
	"appengine"
	"appengine/blobstore"
	"appengine/delay"
	"appengine/image"
	"appengine/urlfetch"
	"bufio"
	"github.com/danielphan/dandubois-net/painting"
	"github.com/danielphan/dandubois-net/painting/legacy"
)

import "net/http"
import "io"
import "encoding/json"

const paintingsUrl = baseUrl + "paintings.json"

type paintings []*legacy.Painting

func fetchPaintings(c *http.Client) (paintings, error) {
	res, err := c.Get(paintingsUrl)
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

func (p paintings) convert(c appengine.Context, cs categories, ms media) {
	for _, o := range p {
		old := o.Fields
		p := &painting.Painting{
			Number:      old.Number,
			Title:       old.Title,
			Image:       appengine.BlobKey(" "),
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
			p.ImageUrl = baseUrl + old.Image
		}
		p.Created(c)
		saveLater.Call(c, p)
	}
}

var saveLater *delay.Function

func init() {
	saveLater = delay.Func("save", save)
}

func save(c appengine.Context, p *painting.Painting) error {
	if p.ImageUrl != "" {
		// Fetch the image.
		res, err := urlfetch.Client(c).Get(p.ImageUrl)
		if err != nil {
			return err
		}

		// Get the content type.
		ct := res.Header.Get("Content-Type")
		if ct == "" {
			h, err := bufio.NewReader(res.Body).Peek(512)
			if err != nil {
				return err
			}
			ct = http.DetectContentType(h)
		}

		// Create a new blob and copy the image into it.
		b, err := blobstore.Create(c, ct)
		if err != nil {
			return err
		}
		_, err = io.Copy(b, res.Body)
		if err != nil {
			return err
		}
		err = b.Close()
		if err != nil {
			return err
		}

		// Add the image to the painting.
		p.Image, err = b.Key()
		if err != nil {
			return err
		}
		u, err := image.ServingURL(c, p.Image, nil)
		if err != nil {
			return err
		}
		p.ImageUrl = u.String()
	}
	return p.Save(c)
}
