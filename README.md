# ghupload

[![build](https://github.com/invit/ghupload/actions/workflows/build.yml/badge.svg)](https://github.com/invit/ghupload/actions/workflows/build.yml)

CLI to upload a file to a github repository through [github's REST API](https://docs.github.com/en/rest/reference/repos#create-or-update-file-contents). The file will be either created or overwritten.

## Installation

Downloadable binaries are available from the [releases page](https://github.com/invit/ghupload/releases/latest).

## Setup

* Create a [Github personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)
* Set _GITHUB_TOKEN_ environment variable 

```shell
$ export GITHUB_TOKEN="your-personal-access-token"
```

## Usage

```
Usage:
  ghupload upload <local-path> <remote-url> [flags]

Flags:
  -b, --branch string    Commit to branch (default branch if empty)
  -h, --help             help for upload
  -m, --message string   Commit message (required)
```
_local-path_ is either a path to a local file or "-" for STDIN. 

_remote-url_ can be one of the following formats and has to include the repository owner, the repository and the path to the file inside the repository:
* https://github.com/owner/repository/path/in/repo
* git@github.com:owner/repository.git/path/in/repo
* owner/repository/path/in/repo

### Examples

* Upload file

```shell
$ ghupload upload -m "commit msg" README.md invit/ghupload/README.md
```

* Upload data from STDIN

```shell
$ ghupload upload -m "commit msg" - invit/ghupload/README.md
this is the new 
content 
of the file
^D
```

## Build

On Linux:

```
$ git clone github.com/invit/ghupload 
$ cd ghupload
$ make 
```

## License

ghupload is licensed under the [MIT License](http://opensource.org/licenses/MIT).
