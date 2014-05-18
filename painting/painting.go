package painting

import (
	"appengine"
	"github.com/danielphan/object"
)

const Kind = "Painting"

type ID string

type Image struct {
	BlobKey       appengine.BlobKey `json:"-"`
	URL           string
	Width, Height int
}

type Painting struct {
	object.Object
	Title       string
	Image       Image
	Description string `datastore:",noindex"`
	Year        int
	Categories  []string
	Media       []string
	Price       int
	Width       int
	Height      int
	ForSale     bool
}

func Get(c appengine.Context, id ID) (*Painting, error) {
	var p Painting
	err := object.Get(c, Kind, string(id), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetAll(c appengine.Context) ([]*Painting, error) {
	return nil, nil
}

func (p *Painting) Save(c appengine.Context) error {
	return object.Save(c, Kind, p)
}
