# godt
A small command-line tool for listing .odt document metadata and adding/removing file extensions.

**Note:** godt operates on all .odt documents in the current directory.

Command format:

`godt [argument]`

## All available options:

`help` - Show help (list options)

`rmex` - Remove .odt extensions

`adex` - Add .odt extensions

`lsdc` - List by date created

`lswd` - List by word count

`lsch` - List by character count

`lspg` - List by page count

`lspa` - List by paragraph count

`lsim` - List by image count

`lstb` - List by table count

`lsnw` - List by non-white-space character count

`lsob` - List by object count

## Example output:

    $godt lswd

    226            Document 12
    148            Document 4
    103            Document 10
    67             Document 8
    41             Document 13
    39             Document 7
    14             Document 2
    14             Document 3
    13             Document 11
    12             Document 6
    12             Document 9
    5              Document 5
    3              Document 1

