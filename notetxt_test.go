package notetxt

import (
        "testing"
        "github.com/stretchr/testify/assert"
)

func TestTitleToFile (t *testing.T) {
        in1 := "Some interesting title"
        assert.Equal(t, TitleToFilename(in1), "some-interesting-title")

        in2 := "SomeName: a nice Name"
        assert.Equal(t, TitleToFilename(in2), "somename-a-nice-name")

        in3 := "HaHaHa !! funny characters ,, .af425-1.q@22@#%^^&**@!"
        assert.Equal(t, TitleToFilename(in3), "hahaha-funny-characters-.af425-1.q22")
}
