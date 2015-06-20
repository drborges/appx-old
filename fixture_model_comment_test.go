package appx_test

import "github.com/drborges/appx"

type Comment struct {
	appx.Model
	Content string
}

func (this Comment) KeyMetadata() *appx.KeyMetadata {
	return &appx.KeyMetadata{
		Kind:      "Comments",
		HasParent: true,
	}
}
