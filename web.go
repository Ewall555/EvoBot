package web

import "embed"

//go:embed   frontend/dist/*
var UIFs embed.FS
