package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mmbros/chordpro/pkg/chordpro"
)

// getFrontMatterFromReader function returns the front matter of the io.Reader stream.
// Only TOML and YAML front matter are handled.
// - TOML: identified by opening and closing +++.
// - YAML: identified by opening and closing ---.
// It returns an empty string if no front matter is found.
// NOTE: can't handle JSON and ORG front matter
func getFrontMatterFromReader(r io.Reader) string {

	var isTOML, isYAML, inFrontMatter bool

	builder := strings.Builder{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {

		switch txt := scanner.Text(); txt {
		case "+++":
			// TOML

			fmt.Fprintln(&builder, txt) // append the line
			if inFrontMatter {
				if isTOML {
					// end of TOML front matter
					return builder.String()
				}
			} else {
				// start of TOML front matter
				inFrontMatter = true
				isTOML = true
			}
		case "---":
			// YAML

			fmt.Fprintln(&builder, txt) // append the line
			if inFrontMatter {
				if isYAML {
					// end of YML front matter
					return builder.String()
				}
			} else {
				// start of YAML front matter
				inFrontMatter = true
				isYAML = true
			}
		default:
			if inFrontMatter {
				fmt.Fprintln(&builder, txt) // append the line
			} else {
				txt = strings.TrimSpace(txt)
				if len(txt) > 0 {
					// exit on the first not empty line
					return ""
				}
			}
		}
	}

	// NOTE: in case of unclosed FrontMatter, returns all the lines found
	return builder.String()
}

// getFrontMatter function returns the front matter of file.
// See getFrontMatterFromReader for more information.
func getFrontMatter(filepath string) string {

	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer f.Close()
	return getFrontMatterFromReader(f)
}

// appendFrontMatter function prints to the writer a
// YAML front matter based on the song.
func appendFrontMatter(w io.Writer, song *chordpro.Song) {
	var s string
	fmt.Fprintln(w, `---`)

	// title
	s = song.Title()
	fmt.Fprintf(w, "title: %q\n", s)

	// artist
	if s = song.Artist(); s == "" {
		s = song.SubTitle()
	}
	if s != "" {
		fmt.Fprintf(w, "artist: %q\n", s)
	}

	// album
	if s = song.Album(); s != "" {
		fmt.Fprintf(w, "album: %q\n", s)
	}

	// year
	if s = song.Year(); s != "" {
		fmt.Fprintf(w, "year: %q\n", s)
	}

	fmt.Fprintln(w, `---`)
}
