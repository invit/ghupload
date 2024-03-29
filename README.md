# ghupload

[![build](https://github.com/invit/ghupload/actions/workflows/build.yml/badge.svg)](https://github.com/invit/ghupload/actions/workflows/build.yml)

CLI to upload a file to a github repository through [github's REST API](https://docs.github.com/en/rest/reference/repos#create-or-update-file-contents). The file will be either created or overwritten.

## Installation

Downloadable binaries are available from the [releases page](https://github.com/invit/ghupload/releases/latest).

## Setup

Create a [Github personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) and store it in one of the following ways:

*  _GITHUB_TOKEN_ environment variable 

```shell
$ export GITHUB_TOKEN="your-personal-access-token"
```

* In a _github-token_ file. See help message for the upload command for the location on your system. On Linux this is usally _$HOME/.config/ghupload/github-token_, on MacOS _$HOME/Library/Application/ghupload/github-token_ and _%AppData%/ghupload/github-token_ on Windows. 

* Provide a _token_ parameter to the upload command, pointing to a file containing the access token. 

## Usage

```
Usage:
  ghupload upload -m <commit-msg> [-b <branch>] [-t <token-file>] <local-path> <remote-url>

Flags:
  -b, --branch string    Commit to branch (default branch if empty)
  -m, --message string   Commit message (required)
  -t, --token string     File to read token from (default /home/USER/.config/ghupload/github-token)
```

_local-path_ is either a path to a local file or "-" for STDIN. 

_remote-url_ can be one of the following formats and has to include the repository owner, the repository and the path to the file inside the repository:
* https://github.com/owner/repository/path/in/repo
* git@github.com:owner/repository.git/path/in/repo
* owner/repository/path/in/repo

Command prints the commit SHA on success.

### Examples

* Upload file

```shell
$ ghupload upload -m "commit msg" README.md invit/ghupload/README.md
b6cbb5b2ea041956c4ac8da17007f95d2312a461
```

* Upload data from STDIN

```shell
$ ghupload upload -m "commit msg" - invit/ghupload/README.md
this is the new 
content 
of the file
^D
3be39e60c3ae44faa40f4efc31241f3564c396f1
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
