package notetxt

import (
        "regexp"
        "strings"
)

title_clearer = regexp.MustCompile("[^a-zA-Z0-9\\s\\.-_]+")
whitespace_clearer = regexp.MustCompile("\\s+")

func TitleToFilename (title string) string {
        out := title_clearer.ReplaceAllString(title, "")
        out = strings.ToLower(out)
        out = whitespace_clearer.ReplaceAllString(out, " ")
        return out
}
