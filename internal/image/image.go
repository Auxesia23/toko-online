package image

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary(url string) *cloudinary.Cloudinary {
	
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		log.Fatal("Failed to initialize Cloudinary", err)
	}

	return cld
}
