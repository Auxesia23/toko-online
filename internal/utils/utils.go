package utils

import "github.com/Auxesia23/toko-online/internal/env"

var (
	secretKey            string
	senderEmail          string
	emailPassword        string
	emailVerificationUrl string
)

func InitUtils() {
	secretKey = env.GetString("SECRET_KEY", "")
	senderEmail = env.GetString("EMAIL", "")
	emailPassword = env.GetString("EMAIL_PASSWORD", "")
	emailVerificationUrl = env.GetString("VERIFICATION_URL", "")
}
