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

    226            DocumentOne
    148            DocumentTwo
    103            DocumentThree
    67             DocumentFour
    41             DocumentFive
    39             DocumentSix
    14             DocumentSeven
    14             DocumentEight
    13             DocumentNine
    12             DocumentTen
    12             DocumentEleven
    5              DocumentTwelve
    3              DocumentThirteen

