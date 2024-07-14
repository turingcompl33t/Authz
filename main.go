package main

import (
	"refl/internal/authz"
)

func main() {
	p := authz.ExprParser{}
	p.Parse("[]uint{123")
}
