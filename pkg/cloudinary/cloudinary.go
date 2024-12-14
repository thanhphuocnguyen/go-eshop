package cloudinary

import "github.com/cloudinary/cloudinary-go/v2"

func NewCloudinary(url string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(url)
	return cld, err
	// Cloudinary configuration
}
