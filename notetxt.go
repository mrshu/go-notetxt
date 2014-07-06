package notetxt

import (
        "regexp"
        "strings"
        "os"
        "bufio"
        "errors"
        "path"
        "io/ioutil"
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

func readFilesInDir(dir string, subdir string) ([]string, []string) {
        var symlinks []string
        var files []string
        contents, _ := ioutil.ReadDir(dir + "/" + subdir + "/")
        for _, f := range contents {
                if f.IsDir() {
                        t_files, t_syms := readFilesInDir(dir, subdir + "/" + f.Name())
                        files = append(files,t_files...)
                        symlinks = append(symlinks, t_syms...)
                } else {
                        if f.Mode() & os.ModeSymlink != 0 {
                                symlinks = append(symlinks, dir + "/" + subdir + "/" + f.Name())
                        } else {
                                files = append(files, dir + "/" + subdir + "/" + f.Name())
                        }
                }
        }
        return files, symlinks
}

func findCategories(notedir string, subdir string, name string) []string {
        var out []string
        files, _ := ioutil.ReadDir(notedir + "/" + subdir + "/")
        for _, f := range files {
                if f.IsDir() {
                        out = append(out, findCategories(notedir, subdir + "/" + f.Name(), name)...)
                } else {
                        if f.Name() == name && subdir != "" {
                                out = append(out, subdir)
                        }
                }
        }
        return out
}

func ParseNote(notedir string, filename string) (Note, error) {
        var note = Note{}
        note.filename = filename

        f, err := os.Open(filename)
        if err != nil {
                return note, err
        }

        defer f.Close()
        reader := bufio.NewReaderSize(f, 4*1024)

        line, prefix, err := reader.ReadLine()
        if err != nil {
                return note, err
        }

        if prefix {
                return note, errors.New("Buffer reader too small for the name of the note.")
        }

        note.name = string(line)

        note.categories = findCategories(notedir, "", path.Base(filename))

        return note, nil
}

