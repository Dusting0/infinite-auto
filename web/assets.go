package web

import "embed"

// frontendFS contains the built React application served by RegisterRoutes.
//
//go:embed static/dist
var frontendFS embed.FS
