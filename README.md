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

    226            DocumentA
    148            DocumentB
    103            DocumentL
    67             DocumentK
    41             DocumentM
    39             DocumentJ
    14             DocumentC
    14             DocumentI
    13             DocumentD
    12             DocumentH
    12             DocumentE
    5              DocumentF
    3              DocumentG

