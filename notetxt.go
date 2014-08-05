package notetxt

import (
        "regexp"
        "strings"
        "os"
        "bufio"
        "errors"
        "path"
        "path/filepath"
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

func FilenameMatches(filename string) bool {
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
        contents, _ := ioutil.ReadDir(dir + "/" + subdir)
        for _, f := range contents {
                if f.IsDir() {
                        t_files, t_syms := readFilesInDir(dir, subdir + "/" + f.Name())
                        files = append(files,t_files...)
                        symlinks = append(symlinks, t_syms...)
                } else {
                        if f.Mode() & os.ModeSymlink != 0 {
                                symlinks = append(symlinks, dir + subdir + "/" + f.Name())
                        } else {
                                files = append(files, dir + subdir + "/" + f.Name())
                        }
                }
        }
        return files, symlinks
}

func findCategories(filename string, notedir string, symlinks []string) []string {
        var out []string
        for _, f := range symlinks {
                p, err := filepath.EvalSymlinks(f)
                if err != nil {
                        panic(err);
                }

                if "./" + p == filename {
                        out = append(out, strings.Replace("./" + path.Dir(f), notedir, "", 1))
                }
        }

        return out
}

func ParseNote(notedir string, filename string, symlinks []string) (Note, error) {
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

        note.categories = findCategories(filename, notedir, symlinks)

        return note, nil
}

func ParseDir(notedir string) ([]Note, error) {
        var notes []Note

        files, symlinks := readFilesInDir(notedir, "")

        for _, f := range files {
                note, err := ParseNote(notedir, f, symlinks)
                if err != nil {
                        return nil, err
                }

                notes = append(notes, note)
        }

        return notes, nil
}

