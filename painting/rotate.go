package painting

import (
	"appengine"
	"appengine/blobstore"
	"appengine/delay"
	aeimage "appengine/image"
	"code.google.com/p/graphics-go/graphics"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
)

// rotate rotates the Painting's image counter-clockwise by angle in degrees.
// angle modulo 360 must be one of 0, 90, 180, 270.
func (p *Painting) rotate(c appengine.Context, angle int) error {
	switch math.Abs(float64(angle % 360)) {
	case 0, 90, 180, 270:
		break
	default:
		return errors.New(fmt.Sprintf("painting: Unsupported angle %f.", angle))
	}

	if p.Image == (Image{}) {
		return nil
	}

	// Read the image from the blobstore.
	r := blobstore.NewReader(c, p.Image.BlobKey)
	src, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	// Create the rotated image.
	srcRect := src.Bounds()
	var dstRect image.Rectangle
	if angle == 0 || angle == 180 {
		dstRect = srcRect
	} else {
		dstRect = image.Rect(0, 0, srcRect.Dy(), srcRect.Dx())
	}

	dst := image.NewNRGBA(dstRect)
	err = graphics.Rotate(dst, src, &graphics.RotateOptions{
		Angle: float64(angle%360) * math.Pi / 180,
	})
	if err != nil {
		return err
	}

	// Create a new blob for the rotated image.
	w, err := blobstore.Create(c, "image/png")
	if err != nil {
		return err
	}

	err = png.Encode(w, dst)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	// Delete the old blob.
	deleteBlobLater.Call(c, p.Image.BlobKey)

	// Update the image metadata.
	p.Image.BlobKey, err = w.Key()
	if err != nil {
		return err
	}

	p.Image.Width = dstRect.Dx()
	p.Image.Height = dstRect.Dy()

	u, err := aeimage.ServingURL(c, p.Image.BlobKey, nil)
	if err != nil {
		return err
	}
	p.Image.URL = u.String()

	err = p.Save(c)
	if err != nil {
		return err
	}

	return nil
}

var deleteBlobLater *delay.Function

func init() {
	deleteBlobLater = delay.Func("deleteBlob", deleteBlob)
}

func deleteBlob(c appengine.Context, key appengine.BlobKey) error {
	return blobstore.Delete(c, key)
}
