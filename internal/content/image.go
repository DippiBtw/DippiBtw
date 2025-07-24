package content

type Image struct {
	image string
}

func (i *Image) String() string {
	return i.image
}
