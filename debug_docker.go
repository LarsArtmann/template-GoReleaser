package main

import (
	"fmt"

	"github.com/LarsArtmann/GoReleaser-Wizard/internal/validation"
)

func main() {
	te := validation.NewTemplateEscaper()
	result := te.EscapeDockerLabel("my@app$")
	fmt.Printf("Result: %q\n", result)
}
