package painting

import (
	"appengine"
	"appengine/datastore"
	"github.com/danielphan/dandubois-net/object"
)

const kind = "Painting"

type Painting struct {
	object.Object

	Number      int
	Title       string
	Image       appengine.BlobKey `json:"-"`
	ImageUrl    string
	Description string `datastore:",noindex"`
	Year        int
	Categories  []string
	Media       []string
	Price       int
	Width       int
	Height      int
	ForSale     bool
}

func (p *Painting) Created(c appengine.Context) error {
	id, _, err := datastore.AllocateIDs(c, kind, nil, 1)
	if err != nil {
		return err
	}
	p.Object.Created(id)
	return nil
}

func (p *Painting) Save(c appengine.Context) error {
	modified, err := p.Modified(p)
	if !modified || err != nil {
		return err
	}

	key := datastore.NewKey(c, kind, p.Id, 0, nil)
	_, err = datastore.Put(c, key, p)
	return err
}
