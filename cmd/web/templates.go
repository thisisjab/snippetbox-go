package main

import "web-dev-journey/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
