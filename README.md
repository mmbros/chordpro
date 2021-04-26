# chordpro

Utility to converts `chordpro` files to `html` format.

    chordpro <command> [options] [args]  

Command:

    transform (folder, dir)  transform all the chordpro files in the source folder
    transform-file (file)    transform a single chordpro file
    clear                    clear

## transform

Recursively transform all the `chordpro` files in the source folder and saves them to the corrisponding location in the dest folder.

    chordpro transform [options] <source-folder> <dest-folder> 

Options:

    -o, --overwrite <overwrite-mode>
          how to handle existing output file (default "none")
            "none"     : never overwrite existing files
            "old"      : overwrite older files
            "all"      : overwrite all files
    -f, --frontmatter string
          how to handle frontmatter (default "preserve")
            "none"     : don't print frontmatter
            "overwrite": overwrite existing frontmatter
            "preserve" : preserve existing frontmatter
    -i, --index
          recursively creates "_index.md" files for folders
    -h, --help
          print this help message


## transform-file


Transform a `chordpro` file to `html` format.

    chordpro transform-file [options] <source-file> [<dest-file>=StdOut] 

Options:

    -o, --overwrite <overwrite-mode>
          how to handle existing output file (default "none")
            "none"     : never overwrite existing files
            "old"      : overwrite older files
            "all"      : overwrite all files
    -f, --frontmatter string
          how to handle frontmatter (default "preserve")
            "none"     : don't print frontmatter
            "overwrite": overwrite existing frontmatter
            "preserve" : preserve existing frontmatter
    -h, --help
          print this help message
