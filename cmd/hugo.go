package cmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmbros/chordpro/pkg/chordpro"
)

func copyFile(srcFile, dstFile string, overwrite overwriteMode) error {

	if overwrite != modeOverwriteAll {
		// check if output file already exists
		outFileinfo, err := os.Stat(dstFile)
		if !os.IsNotExist(err) {
			// file esists
			if overwrite == modeOverwriteNone {
				// skip
				// return ErrOutputFileExists
				return nil
			}

			// case modeOverwriteOld
			// check time
			inFileinfo, err2 := os.Stat(srcFile)
			if err2 != nil {
				return err2
			}

			if outFileinfo.ModTime().After(inFileinfo.ModTime()) {
				// skip
				// return ErrOutputFileNewer
				return nil
			}
		}
		// else file not exists
	}

	input, err := ioutil.ReadFile(srcFile)
	if err == nil {
		err = ioutil.WriteFile(dstFile, input, 0644)
	}

	return err
}

// transform function parse a chordpro.Songs object from io.Reader and output the resut to io.Writer.
// It returns an error if the number of songs is not exactly one.
// If the flag songFrontmatter is true, the first part of the result is the front matter created from the song metadata.
// Then it prints the given prefix.
// At last it prints the formatted song.
func transform2(r io.Reader, w io.Writer) error {

	// retrieve from reader
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// parse string
	file := chordpro.ParseText(string(data))
	songs := file.Songs

	// check number of songs
	if totSongs := len(songs); totSongs == 0 {
		return ErrZeroSongs
	} else if totSongs > 1 {
		return ErrMultipleSongs
	}

	// format the first song, discard the others
	formatter := chordpro.NewHtmlDivFormatter(w)
	song := songs[0]

	if file.Frontmatter != "" {
		w.Write([]byte(file.Frontmatter))
	} else {
		appendFrontMatter(w, song)
	}

	formatter.FormatBody(song)

	return nil
}

// trasformFile dunction transforms the ChordPro source file
// into the HTML destination file.
func trasformFileHugo(srcFile, dstFile string, overwrite overwriteMode) error {

	err := checkFiles(srcFile, dstFile, overwrite)
	if err != nil {
		return err
	}

	// source file
	fin, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer fin.Close()

	// writer
	fout := os.Stdout
	if dstFile != "" {
		fout, err = createFileAll(dstFile)
		if err != nil {
			return err
		}
		defer fout.Close()
	}
	writer := bufio.NewWriter(fout)

	reader := bufio.NewReader(fin)

	err = transform2(reader, writer)
	writer.Flush()

	return err
}

// // trasformFile dunction transforms the ChordPro source file
// // into the HTML destination file.
// func trasformFileHugoOLD(srcFile, dstFile string, overwrite overwriteMode) error {

// 	err := checkFiles(srcFile, dstFile, overwrite)
// 	if err != nil {
// 		return err
// 	}

// 	// retrieve from reader
// 	data, err := ioutil.ReadFile(srcFile)
// 	if err != nil {
// 		return err
// 	}

// 	// s := toUtf8(data)
// 	s := string(data)

// 	// writer
// 	fout := os.Stdout
// 	if dstFile != "" {
// 		fout, err = createFileAll(dstFile)
// 		if err != nil {
// 			return err
// 		}
// 		defer fout.Close()
// 	}
// 	writer := bufio.NewWriter(fout)

// 	saveFrontMatter := getFrontMatterFromReader(strings.NewReader(s))

// 	L := len(saveFrontMatter)
// 	j := 0

// 	if L > 0 {
// 		j = strings.Index(s, saveFrontMatter)
// 		s = s[j+L:]
// 	}
// 	s = strings.TrimSpace(s)

// 	songFrontmatter := (saveFrontMatter == "")

// 	// fmt.Println(saveFrontMatter)
// 	// fmt.Println("###############################################################################")

// 	// fmt.Println(s)

// 	err = transform(strings.NewReader(s), writer, saveFrontMatter, songFrontmatter)
// 	writer.Flush()

// 	return err
// }

// Run executes the transformation from chordpro files to html files
// based on the given options.
func runHugo(opts *Options) error {

	overwrite, err := parseOverwrite(opts.Overwrite)
	if err != nil {
		return err
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

			if !info.IsDir() {
				ext := filepath.Ext(path)
				switch strings.ToLower(ext) {
				case ".cho", ".chopro", ".chordpro":
					dstpath := filepath.Join(opts.Output, relpath) + ".html"
					fmt.Println(relpath)
					trasformFileHugo(path, dstpath, overwrite)
					// if err != nil && err != ErrOutputFileNewer {
					// 	return err
					// }
				default:
					// copy
					dstpath := filepath.Join(opts.Output, relpath)
					err = copyFile(path, dstpath, overwrite)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
					return err

				}
			}
			return nil
		})

	return err
}
