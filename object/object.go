package object

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strconv"
	"time"
)

type Object struct {
	Id         string
	State      string
	ModifiedAt time.Time
	CreatedAt  time.Time
}

func (o *Object) Created(id int64) {
	o.Id = strconv.FormatInt(id, 16)
	o.CreatedAt = time.Now()
}

func (o *Object) Modified(i interface{}) (bool, error) {
	h := sha256.New()
	b, err := json.Marshal(i)
	_, err = io.Copy(h, bytes.NewReader(b))
	if err != nil {
		return false, err
	}
	s := hex.EncodeToString(h.Sum(nil))
	if o.State != s {
		o.State = s
		o.ModifiedAt = time.Now()
		return true, nil
	}
	return false, nil
}
