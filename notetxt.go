package notetxt

import (
        "regexp"
        "strings"
)
var title_clearer = regexp.MustCompile("[^a-zA-Z0-9\\s\\.\\-_]+")
var whitespace_clearer = regexp.MustCompile("\\s+")

func TitleToFilename (title string) string {

        out := title_clearer.ReplaceAllString(title, "")
        out = strings.ToLower(out)
        out = whitespace_clearer.ReplaceAllString(out, " ")
        out = strings.Replace(out, " ", "-", -1)
        return out
}

var filename_regex = regexp.MustCompile("^[a-zA-Z0-9\\-\\.]+$")

func MatchesNoteFilename(filename string) bool {
        return filename_regex.MatchString(filename)
}

type Note struct {
        name string
        filename string
        categories []string
}

func ParseNote(notedir string, filename string) (Note) {
        var note = Note{}

        return note
}

