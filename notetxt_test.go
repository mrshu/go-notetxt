package notetxt

import (
        "testing"
        "path/filepath"
        "github.com/stretchr/testify/assert"
)

var testdir string

func init() {
        testdir, _ = filepath.Abs("./test")
}

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
        assert.Equal(t, FilenameMatches(in1), true)

        in2 := "filewithoutextension"
        assert.Equal(t, FilenameMatches(in2), true)

        in3 := "non matching filename"
        assert.Equal(t, FilenameMatches(in3), false)

}

func TestNoteParsing(t *testing.T) {
        _, symlinks := readFilesInDir(testdir, "")
        file, _ := filepath.Abs("./test/some-nice-title.rst")

        var note, err = ParseNote(testdir, file, symlinks)

        tags := make([]string, 2)
        tags[0] = "/tag/general"
        tags[1] = "/tag/title"

        assert.Equal(t, err, nil)
        assert.Equal(t, note.Name, "Some nice title")
        assert.Equal(t, note.Filename, file)
        assert.Equal(t, note.Tags, tags)
}

func TestFindTags(t *testing.T) {
        _, symlinks := readFilesInDir(testdir, "")
        file, _ := filepath.Abs("./test/some-nice-title.rst")

        var note = findTags(file, testdir, symlinks)

        tags := make([]string, 2)
        tags[0] = "/tag/general"
        tags[1] = "/tag/title"

        assert.Equal(t, note, tags)

        another_file, _ := filepath.Abs("./test/tag/general/just-a-tag.rst")

        var another_note = findTags(another_file, testdir, symlinks)

        assert.Equal(t, len(another_note), 1)

}

func TestDirListing(t *testing.T) {
        tfiles, tsymlinks := readFilesInDir(testdir, "")

        files := make([]string, 3)
        files[0], _ = filepath.Abs("./test/myproject-a-more-complicated-title.rst")
        files[1], _ = filepath.Abs("./test/some-nice-title.rst")
        files[2], _ = filepath.Abs("./test/tag/general/just-a-tag.rst")

        symlinks := make([]string, 4)
        symlinks[0], _ = filepath.Abs("./test/tag/general/some-nice-title.rst")
        symlinks[1], _ = filepath.Abs("./test/tag/project/myproject-a-more-complicated-title.rst")
        symlinks[2], _ = filepath.Abs("./test/tag/title/myproject-a-more-complicated-title.rst")
        symlinks[3], _ = filepath.Abs("./test/tag/title/some-nice-title.rst")

        assert.Equal(t, tfiles, files)
        assert.Equal(t, tsymlinks, symlinks)
}

func TestDirNoteParsing(t *testing.T) {
        notes, _ := ParseDir(testdir)

        assert.Equal(t, len(notes), 3)

        assert.Equal(t, notes[0].Name, "MyProject: A more complicated title")
        assert.Equal(t, notes[1].Name, "Some nice title")
        assert.Equal(t, notes[2].Name, "Just a tag")

        assert.Equal(t, len(notes[2].Tags), 1)
        assert.Equal(t, notes[2].Tags[0], "/tag/general")
}
