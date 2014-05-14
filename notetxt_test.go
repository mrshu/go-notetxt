package notetxt

import (
        "testing"
        "github.com/stretchr/testify/assert"
)

func TestTitleToFile (t *testing.T) {

        // a normal title
        in1 := "Some interesting title"
        assert.Equal(t, TitleToFilename(in1), "some-interesting-title")

        // a title with some funny characters
        in2 := "SomeName: a nice Name"
        assert.Equal(t, TitleToFilename(in2), "somename-a-nice-name")

        // a name full of funny characters
        in3 := "HaHaHa !! funny characters ,, .af425-1.q@22@#%^^&**@!"
        assert.Equal(t, TitleToFilename(in3), "hahaha-funny-characters-.af425-1.q22")

        // a name consisting of funny characters only
        in4 := "!@#$%^&*())!@#$%^^&!$@#%@^$%@$#$@#$@"
        assert.Equal(t, TitleToFilename(in4), "")
}

func TestFilenameMatch(t *testing.T) {
        in1 := "some-interesting-title.rst"
        assert.Equal(t, MatchesNoteFilename(in1), true)

        in2 := "filewithoutextension"
        assert.Equal(t, MatchesNoteFilename(in2), true)

        in3 := "non matching filename"
        assert.Equal(t, MatchesNoteFilename(in3), false)

}
