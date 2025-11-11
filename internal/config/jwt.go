package config

import "os"

var JWTSecret = []byte(getSecret())

func getSecret() string {
    s := os.Getenv("JWT_SECRET")
    if s == "" {
        s = "supersecretkey" 
    }
    return s
}
