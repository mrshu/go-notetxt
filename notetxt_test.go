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

func TestNoteParsing(t *testing.T) {
        _, symlinks := readFilesInDir("./test", "")
        var note, err = ParseNote("./test", "./test/some-nice-title.rst", symlinks)

        categories := make([]string, 2)
        categories[0] = "/tag/general"
        categories[1] = "/tag/title"

        assert.Equal(t, err, nil)
        assert.Equal(t, note.name, "Some nice title")
        assert.Equal(t, note.filename, "./test/some-nice-title.rst")
        assert.Equal(t, note.categories, categories)
}

func TestFindCategories(t *testing.T) {
        _, symlinks := readFilesInDir("./test", "")
        var note = findCategories("./test/some-nice-title.rst", "./test", symlinks)

        categories := make([]string, 2)
        categories[0] = "/tag/general"
        categories[1] = "/tag/title"

        assert.Equal(t, note, categories)

}

func TestDirListing(t *testing.T) {
        tfiles, tsymlinks := readFilesInDir("./test", "")

        files := make([]string, 3)
        files[0] = "./test/myproject-a-more-complicated-title.rst"
        files[1] = "./test/some-nice-title.rst"
        files[2] = "./test/tag/general/just-a-tag.rst"

        symlinks := make([]string, 4)
        symlinks[0] = "./test/tag/general/some-nice-title.rst"
        symlinks[1] = "./test/tag/project/myproject-a-more-complicated-title.rst"
        symlinks[2] = "./test/tag/title/myproject-a-more-complicated-title.rst"
        symlinks[3] = "./test/tag/title/some-nice-title.rst"

        assert.Equal(t, tfiles, files)
        assert.Equal(t, tsymlinks, symlinks)
}

func TestDirNoteParsing(t *testing.T) {
        notes, _ := ParseDir("./test")

        assert.Equal(t, len(notes), 3)

        assert.Equal(t, notes[0].name, "MyProject: A more complicated title")
        assert.Equal(t, notes[1].name, "Some nice title")
        assert.Equal(t, notes[2].name, "Just a tag")

}
