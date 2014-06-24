package notetxt

import (
        "regexp"
        "strings"
        "os"
        "bufio"
        "errors"
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

func ParseNote(notedir string, filename string) (Note, error) {
        var note = Note{}
        note.filename = filename

        f, err := os.Open(filename)
        if err != nil {
                return note, err
        }

        defer f.Close()
        reader, err := bufio.NewReaderSize(f, 4*1024)
        if err != nil {
                return note, err
        }

        line, prefix, err := reader.ReadLine()
        if err != nil {
                return note, err
        }

        if prefix {
                return note, errors.New("Buffer reader too small for the name of the note.")
        }

        note.name = string(line)

        return note, nil
}

