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
Uploads (commits) a local file to a github repository

<local-path> is either a path to a local file or - for STDIN.
<remote-url> can be one of the following formats and has to include the repository owner, the repository and the path to the file inside the repository:
* https://github.com/owner/repository/path/in/repo
* git@github.com:owner/repository.git/path/in/repo
* owner/repository/path/in/repo

Command prints the commit SHA on success.

Usage:
  ghupload upload -m <commit-msg> [-b <branch>] <local-path> <remote-url>

Examples:
* Upload local file
  $ ghupload upload -m "commit msg" README.md owner/repository/README.md
  b6cbb5b2ea041956c4ac8da17007f95d2312a461
* Upload data from STDIN
  $ ghupload upload -m "commit msg" - owner/repository/README.md
  this is the new
  content
  of the file
  ^D
  3be39e60c3ae44faa40f4efc31241f3564c396f1

Flags:
  -b, --branch string    Commit to branch (default branch if empty)
  -h, --help             help for upload
  -m, --message string   Commit message (required)
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
