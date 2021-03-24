package wrappers

import "github.com/rs/cors"

func NewCorsWrapper() *cors.Cors {
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})
	return corsWrapper
}