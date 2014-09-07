package notetxt

import (
        "regexp"
        "strings"
        "os"
        "fmt"
        "bufio"
        "errors"
        "path"
        "path/filepath"
        "io/ioutil"
)
var title_clearer = regexp.MustCompile("[^a-zA-Z0-9\\s\\.\\-_]+")
var whitespace_clearer = regexp.MustCompile("\\s+")

func TitleToFilename (title string) string {

        // strip all non-conforming characters
        out := title_clearer.ReplaceAllString(title, "")

        // title shall be lowercase
        out = strings.ToLower(out)

        // every whitespace should become a space
        // if there are multiple whitespace insert only one space as a result
        out = whitespace_clearer.ReplaceAllString(out, " ")

        // every white space should become a dash (because they look nice)
        out = strings.Replace(out, " ", "-", -1)
        return out
}

var filename_regex = regexp.MustCompile("^[a-zA-Z0-9\\-\\.]+$")

func FilenameMatches(filename string) bool {
        return filename_regex.MatchString(filename)
}

type Note struct {
        Name string
        Filename string
        Tags []string
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

func findTags(filename string, notedir string, symlinks []string) []string {
        var out []string

        plain_tag := strings.Replace(path.Dir(filename), notedir, "", 1)
        if len(plain_tag) != 0 {
                out = append(out, plain_tag)
        }

        for _, f := range symlinks {
                p, err := filepath.EvalSymlinks(f)
                if err != nil {
                        panic(err);
                }

                if p == filename {
                        out = append(out, strings.Replace(path.Dir(f), notedir, "", 1))
                }
        }

        return out
}

func ParseNote(notedir string, filename string, symlinks []string) (Note, error) {
        var note = Note{}
        note.Filename = filename

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

        note.Name = string(line)

        note.Tags = findTags(filename, notedir, symlinks)

        return note, nil
}

func ParseDir(notedir string) ([]Note, error) {
        var notes []Note

        notedir, _ = filepath.Abs(notedir)
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


func CreateNote(title string, tag string, dir string) (string, error) {
        spacer := "\n" + strings.Repeat("=", len(title))
        text := title + spacer

        directory := fmt.Sprintf("%s/%s", dir, tag)
        err := os.MkdirAll(directory, 755)
        if err != nil {
                return "", err
        }

        file := fmt.Sprintf("%s/%s.rst", directory, TitleToFilename(title))

        if _, err := os.Stat(file); err == nil {
                return "", errors.New("Notefile '" + file + "' already exists. " +
                                        "You can still edit it if you want.")
        }

        e := ioutil.WriteFile(file, []byte(text), 0644)
        if e != nil {
                return "", e
        }

        return file, nil
}

func TagNote(file string, tag string, dir string) error {
        directory := fmt.Sprintf("%s/%s", dir, tag)
        err := os.MkdirAll(directory, 755)
        if err != nil {
                return err
        }

        filename := filepath.Base(file)
        newpath := fmt.Sprintf("%s/%s/%s", dir, tag, filename)

        err = os.Symlink(file, newpath)
        if err != nil {
                return err
        }

        return nil
}
