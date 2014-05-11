package notetxt

import (
        "testing"
        "github.com/stretchr/testify/assert"
)

func TestTitleToFile (t *testing.T) {
        in := "Some interesting title"
        assert.Equal(t, TitleToFilename(in), "some-interesting-title")
}
