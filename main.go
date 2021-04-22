package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/mmbros/chordpro/cmd"
	"github.com/mmbros/simpleflag"
)

const (
	defaultOverwrite   = cmd.OverwriteNone
	defaultFrontmatter = cmd.FrontmatterPreserve

	cmdnameTranformFolder      = "transform"
	cmdnameTranformFolderAlias = "folder, dir"

	cmdnameTranformFile      = "transform-file"
	cmdnameTranformFileAlias = "file"

	cmdnameClearFolder = "clear"
)

var appname string

func usageApp() {
	const msg = `%[1]s : utility to converts chordpro files to html format

Usage: %[1]s <command> [options] [args]  

Command:
  %-24[2]s transform all the chordpro files in the source folder
  %-24[3]s transform a single chordpro file
  %-24[4]s clear
`

	fmt.Fprintf(flag.CommandLine.Output(), msg, appname,
		fmt.Sprintf("%s (%s)", cmdnameTranformFolder, cmdnameTranformFolderAlias),
		fmt.Sprintf("%s (%s)", cmdnameTranformFile, cmdnameTranformFileAlias),
		cmdnameClearFolder,
	)
}

func usageTransformFolder() {
	const msg = `%[1]s %[10]s
    recursively transform all the chordpro files in the source folder
    and saves them to the corrisponding location in the dest folder.

Usage: %[1]s %[10]s [options] <source-folder> <dest-folder> 

Options:
  -o, --overwrite <overwrite-mode>
        how to handle existing output file (default %[2]q)
          %-11[3]q: never overwrite existing files
          %-11[4]q: overwrite older files
          %-11[5]q: overwrite all files
  -f, --frontmatter string
        how to handle frontmatter (default %[6]q)
          %-11[7]q: don't print frontmatter
          %-11[8]q: overwrite existing frontmatter
          %-11[9]q: preserve existing frontmatter
  -i, --index
        recursively creates "_index.md" files for folders
  -h, --help
        print this help message
`

	fmt.Fprintf(flag.CommandLine.Output(), msg, appname,
		defaultOverwrite, cmd.OverwriteNone, cmd.OverwriteOld, cmd.OverwriteAll,
		defaultFrontmatter, cmd.FrontmatterNone, cmd.FrontmatterOverwrite, cmd.FrontmatterPreserve,
		cmdnameTranformFolder,
	)
}
func usageTransformFile() {
	const msg = `%[1]s %[10]s
    transform a chordpro file to html format.

Usage: %[1]s %[10]s [options] <source-file> [<dest-file>=StdOut] 

Options:
  -o, --overwrite <overwrite-mode>
        how to handle existing output file (default %[2]q)
          %-11[3]q: never overwrite existing files
          %-11[4]q: overwrite older files
          %-11[5]q: overwrite all files
  -f, --frontmatter string
        how to handle frontmatter (default %[6]q)
          %-11[7]q: don't print frontmatter
          %-11[8]q: overwrite existing frontmatter
          %-11[9]q: preserve existing frontmatter
  -h, --help
        print this help message
`

	fmt.Fprintf(flag.CommandLine.Output(), msg, appname,
		defaultOverwrite, cmd.OverwriteNone, cmd.OverwriteOld, cmd.OverwriteAll,
		defaultFrontmatter, cmd.FrontmatterNone, cmd.FrontmatterOverwrite, cmd.FrontmatterPreserve,
		cmdnameTranformFile,
	)
}

func cmdApp(name string, arguments []string) error {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = usageApp

	if len(arguments) == 0 {
		fs.Usage()
		return nil
	}
	err := fs.Parse(arguments)
	return err
}

func cmdTransformFolder(name string, arguments []string) error {

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	var opts cmd.Options

	fs.Usage = usageTransformFolder

	opts.Recursive = true
	simpleflag.AliasedStringVar(fs, &opts.Overwrite, "overwrite,o", defaultOverwrite, "")
	simpleflag.AliasedStringVar(fs, &opts.Frontmatter, "frontmatter,f", defaultFrontmatter, "")
	simpleflag.AliasedBoolVar(fs, &opts.Index, "index,i", false, "")

	err := fs.Parse(arguments)
	if err != nil {
		return err
	}
	opts.Input = fs.Arg(0)
	opts.Output = fs.Arg(1)

	err = cmd.Run(&opts)
	if err != nil {
		err = simpleflag.WrapError(err, name)
	}
	return err
}

func cmdTransformFile(name string, arguments []string) error {

	fs := flag.NewFlagSet(cmdnameTranformFile, flag.ContinueOnError)
	var opts cmd.Options

	fs.Usage = usageTransformFile
	simpleflag.AliasedStringVar(fs, &opts.Overwrite, "overwrite,o", defaultOverwrite, "")
	simpleflag.AliasedStringVar(fs, &opts.Frontmatter, "frontmatter,f", defaultFrontmatter, "")
	simpleflag.AliasedBoolVar(fs, &opts.Index, "index,i", false, "")

	err := fs.Parse(arguments)
	if err != nil {
		return err
	}
	opts.Input = fs.Arg(0)
	opts.Output = fs.Arg(1)

	err = cmd.Run(&opts)
	if err != nil {
		err = simpleflag.WrapError(err, name)
	}
	return err
}

func main() {
	appname = path.Base(os.Args[0])

	var app = &simpleflag.Command{
		ParseExec: cmdApp,
		SubCmd: map[string]*simpleflag.Command{
			cmdnameTranformFolder + "," + cmdnameTranformFolderAlias: {
				ParseExec: cmdTransformFolder,
			},
			cmdnameTranformFile + "," + cmdnameTranformFileAlias: {
				ParseExec: cmdTransformFile,
			},
		},
	}

	err := simpleflag.ParseExec(app)
	if err == flag.ErrHelp {
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err.Error())
		os.Exit(1)
	}
}
