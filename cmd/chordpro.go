package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmbros/chordpro/pkg/chordpro"
)

const (
	OverwriteNone = "none"
	OverwriteOld  = "old"
	OverwriteAll  = "all"
)

const (
	FrontmatterNone      = "none"
	FrontmatterPreserve  = "preserve"
	FrontmatterOverwrite = "overwrite"
)

type Options struct {
	Input       string // source file / folder
	Output      string // destination file / folder
	Overwrite   string // overwrite mode: "none", "old" or "all"
	Frontmatter string // front matter mode: "none", "preserve", "overwrite"
	Recursive   bool   // recursively transforms every chord file found in the input folder
	Index       bool   // recursively creates "_index.md" files for folders (only for recursive mode)
}

// internal overwrite values
type overwriteMode int

const (
	modeOverwriteNone overwriteMode = iota
	modeOverwriteOld
	modeOverwriteAll
)

// internal front matter values
type frontmatterMode int

const (
	modeFrontmatterNone frontmatterMode = iota
	modeFrontmatterPreserve
	modeFrontmatterOverwrite
)

// Error messages
var (
	// ErrInvalidOverwrite is returned when overwrite string is not valid.
	ErrInvalidOverwrite = errors.New("invalid overwrite")

	// ErrInvalidFrontmatter is returned when frontmatter string is not valid.
	ErrInvalidFrontmatter = errors.New("invalid frontmatter")

	// ErrMissingInput is returned when input file is not specified.
	ErrMissingInput = errors.New("missing input path")

	// ErrInputFileNotFound is returned when input file is not found.
	ErrInputFileNotFound = errors.New("input file not found")

	// ErrInputFileNotRegular is returned when input file is not a regular file.
	ErrInputFileNotRegular = errors.New("input is not a regular file")

	// ErrOutputFileExists is returned where output file already exists
	// and overwrite option was not given.
	ErrOutputFileExists = errors.New("output file already exists")

	// ErrOutputFileNewer is returned where output file already exists
	// and is newer then input file.
	ErrOutputFileNewer = errors.New("output file newer than chordpro input file")

	// ErrZeroSongs is returned when chordpro file does not contain songs.
	ErrZeroSongs = errors.New("no song found")

	// ErrMultipleSongs is returned when chordpro file contains two or more songs.
	ErrMultipleSongs = errors.New("multiple songs found")
)

// parseOverwrite function parses a string into overwriteMode.
// It returns an error in case of unknown input string.
func parseOverwrite(s string) (overwriteMode, error) {
	switch strings.ToLower(s) {
	case OverwriteNone:
		return modeOverwriteNone, nil
	case OverwriteAll:
		return modeOverwriteAll, nil
	case OverwriteOld:
		return modeOverwriteOld, nil
	}
	return modeOverwriteNone, ErrInvalidOverwrite
}

// parseFrontmatter function parses a string into frontmatterMode.
// It returns an error in case of unknown input string.
func parseFrontmatter(s string) (frontmatterMode, error) {
	switch strings.ToLower(s) {
	case FrontmatterNone:
		return modeFrontmatterNone, nil
	case FrontmatterPreserve:
		return modeFrontmatterPreserve, nil
	case FrontmatterOverwrite:
		return modeFrontmatterOverwrite, nil
	}
	return modeFrontmatterNone, ErrInvalidFrontmatter
}

// checkFiles function checks if input and output are valid files
// for the given overwrite mode.
func checkFiles(fin, fout string, overwrite overwriteMode) error {
	var inFileinfo os.FileInfo
	var outFileinfo os.FileInfo
	var err error

	// check input file
	if fin == "" {
		return ErrMissingInput
	}
	inFileinfo, err = os.Stat(fin)
	if os.IsNotExist(err) {
		return ErrInputFileNotFound
	}
	if !inFileinfo.Mode().IsRegular() {
		return ErrInputFileNotRegular
	}

	if (fout == "") || (overwrite == modeOverwriteAll) {
		// don't check output file
		return nil
	}

	// check if output file already exists
	outFileinfo, err = os.Stat(fout)
	if !os.IsNotExist(err) {
		// file esists
		if overwrite == modeOverwriteNone {
			return ErrOutputFileExists
		}

		// case modeOverwriteOld
		// check time
		if outFileinfo.ModTime().After(inFileinfo.ModTime()) {
			return ErrOutputFileNewer
		}
	}

	return nil
}

// checkDirs function checks if input and output are valid directory.
func checkDirs(fin, fout string) error {
	var info os.FileInfo
	var err error

	// check input is a dir
	if fin == "" {
		return ErrMissingInput
	}
	info, err = os.Stat(fin)
	if os.IsNotExist(err) {
		return errors.New("input path not exists")
	}
	if !info.Mode().IsDir() {
		return errors.New("input path not a directory")
	}

	// checks output is a dir or not exists
	if fout == "" {
		return errors.New("missing output dir")
	}
	info, err = os.Stat(fout)
	if os.IsNotExist(err) {
		return nil
	}
	if !info.Mode().IsDir() {
		return errors.New("output path not a directory")
	}

	// TODO: test output is not inside input

	return nil
}

// createFileAll creates or truncates the named file
// along with any necessary parents directory.
func createFileAll(pathname string) (*os.File, error) {

	file, err := os.Create(pathname)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(pathname), 0700); err != nil {
				return nil, err
			}
			return os.Create(pathname)
		}
		return nil, err
	}
	return file, err
}

// transform function parse a chordpro.Songs object from io.Reader and output the resut to io.Writer.
// It returns an error if the number of songs is not exactly one.
// If the flag songFrontmatter is true, the first part of the result is the front matter created from the song metadata.
// Then it prints the given prefix.
// At last it prints the formatted song.
func transform(r io.Reader, w io.Writer, prefix string, songFrontmatter bool) error {

	// retrieve from reader
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// parse string
	songs := chordpro.ParseText(string(data))

	// check number of songs
	if totSongs := len(songs); totSongs == 0 {
		return ErrZeroSongs
	} else if totSongs > 1 {
		return ErrMultipleSongs
	}

	// format the first song, discard the others
	formatter := chordpro.NewHtmlDivFormatter(w)
	song := songs[0]

	if songFrontmatter {
		appendFrontMatter(w, song)
	}
	w.Write([]byte(prefix))
	formatter.FormatBody(song)

	return nil
}

// trasformFile dunction transforms the ChordPro source file
// into the HTML destination file.
func trasformFile(srcFile, dstFile string, overwrite overwriteMode, frontmatter frontmatterMode) error {
	var saveFrontMatter string

	err := checkFiles(srcFile, dstFile, overwrite)
	if err != nil {
		return err
	}

	// reader
	fin, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	// writer
	fout := os.Stdout
	if dstFile != "" {
		if frontmatter == modeFrontmatterPreserve {
			saveFrontMatter = getFrontMatter(dstFile)
		}

		fout, err = createFileAll(dstFile)
		if err != nil {
			return err
		}
		defer fout.Close()
	}
	writer := bufio.NewWriter(fout)

	songFrontmatter := (frontmatter != modeFrontmatterNone) && (saveFrontMatter == "")
	err = transform(fin, writer, saveFrontMatter, songFrontmatter)
	writer.Flush()

	return err
}

// createIndexMD function recursively creates an "_index.md" file
// in each sub folder of root folder, if not already exists.
// The frontmatter of the created "_index.md" file contains only the "title" key
// whose value is the name of the folder.
func createIndexMD(root string) error {

	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if (err != nil) || !info.IsDir() || (path == root) {
				return err
			}

			pathIndex := filepath.Join(path, "_index.md")
			_, err = os.Stat(pathIndex)
			if os.IsNotExist(err) {
				// create index file

				file, err := os.Create(pathIndex)
				if err != nil {
					return err
				}
				fmt.Fprintf(file, "---\ntitle: %q\n---\n", info.Name())
				file.Close()
			}
			return nil
		})
}

// Run executes the transformation from chordpro files to html files
// based on the given options.
func Run(opts *Options) error {

	overwrite, err := parseOverwrite(opts.Overwrite)
	if err != nil {
		return err
	}

	frontmatter, err := parseFrontmatter(opts.Frontmatter)
	if err != nil {
		return err
	}

	if !opts.Recursive {
		// single file mode
		return trasformFile(opts.Input, opts.Output, overwrite, frontmatter)
	}

	err = checkDirs(opts.Input, opts.Output)
	if err != nil {
		return err
	}

	// recursively transforms all chordpro files under the input dir
	err = filepath.Walk(opts.Input,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relpath, _ := filepath.Rel(opts.Input, path)

			dstpath := filepath.Join(opts.Output, relpath) + ".html"

			if !info.IsDir() {
				ext := filepath.Ext(path)
				switch strings.ToLower(ext) {
				case ".cho", ".chopro", ".chordpro":
					err = trasformFile(path, dstpath, overwrite, frontmatter)
					if err == nil {
						fmt.Println(dstpath)
					}
				}
			}
			return nil
		})

	if (err == nil) && opts.Index {
		err = createIndexMD(opts.Output)
	}

	return err
}
