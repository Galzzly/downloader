# downloader

Take in a list of web addresses, and download the files to a target location. Location can be either a local filesystem, or in HDFS.

## Installation

There will be a compiled binary included in the repository, `downloader`, that can be used out of the box. However, if you wish to compile the package yourself, you will need to use the following steps.

1. Download [**Go**](https://golang.org/dl/) and follow the [**installation instructions**](https://go.dev/doc/install)
2. [Set up your install](https://go.dev/doc/gopath_code)
3. Clone the `downloader` repository, and change into the directory
4. Compile the binary by running the following command:

    `go build`

5. The `downloader` binary will have been generated

## Usage

If ran without any arguments, the `downloader` command will display the help message. The same as if it was ran with `downloader -h`, `downloader --help`, and `downloader help`.

Running `downloader -v` or `downloader --version` will display the version.

There are two subcommands:
- `local` - for when you want to download files into a local directory 
- `hdfs` - for when you want to download files directly into HDFS

Both of the subcommands make use of the same flags:
|Short Flag|Long Flag|Description|
|---|---|---|
|`-c`|`--concurrent`| The number of concurrent downloads (default 1)|
|`-d`|`--destination`|The desired destination to save the downloads|
|`-f`|`--file`| The file containing list of URLs|
|`-h`|`--help`| The help message|

## Known Limitations

There may be issues downloading files directly to HDFS in a kerberized environment. This is something that will be worked on, and updated accordingly. 

## Issues

For any issues that are found, please raise an issue on Github.