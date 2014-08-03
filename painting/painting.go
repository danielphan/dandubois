package painting

import (
	"appengine"
	"appengine/datastore"
	"github.com/danielphan/ae/object"
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

func GetAll(c appengine.Context) ([]*Painting, error) {
	q := datastore.NewQuery(Kind).
		Order("-Year").
		Order("-ID")

	var ps []*Painting
	_, err := q.GetAll(c, &ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func Get(c appengine.Context, id ID) (*Painting, error) {
	p := Painting{
		Object: object.New(Kind, string(id)),
	}
	err := object.Get(c, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Painting) Save(c appengine.Context) error {
	return object.Save(c, p)
}

func GetAllCategories(c appengine.Context) ([]string, error) {
	q := datastore.NewQuery(Kind).
		Project("Categories").
		Distinct().
		Order("Categories")

	var paintings []*struct{ Categories []string }
	_, err := q.GetAll(c, &paintings)
	if err != nil {
		return nil, err
	}

	var cs []string
	for _, p := range paintings {
		cs = append(cs, p.Categories[0])
	}
	return cs, nil
}

func GetAllMedia(c appengine.Context) ([]string, error) {
	q := datastore.NewQuery(Kind).
		Project("Media").
		Distinct().
		Order("Media")

	var paintings []*struct{ Media []string }
	_, err := q.GetAll(c, &paintings)
	if err != nil {
		return nil, err
	}

	var ms []string
	for _, p := range paintings {
		ms = append(ms, p.Media[0])
	}

	return ms, nil
}
