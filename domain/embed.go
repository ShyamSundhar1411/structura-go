package domain

import "embed"

//go:embed  templates/*
var TemplateFS embed.FS
