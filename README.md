# Goquette

With **Goquette** you can simplify your packaging experience with NuGet. \
Set up the required files - Steps are listed below - and just run **Goquette**, and you're done.

## Usage:

* Create a 'config.yaml' File inside your root project folder
* Create a 'tools' directory inside your root projekt folder

Inside the `tools`-Directory, put your PowerShell scripts, which are consumed by chocolatey, \
and your zipped program/binary.
Take a look at the tools-Directory inside this **[Tools-Dir of Goquette](https://github.com/PatrickLaabs/goquette/tree/main/tools)** \
for a better understanding.

## Installation

### Install via Go
`go install github.com/PatrickLaabs/goquette@latest` \
Make sure you have set your GOBIN Path correctly. \
If not: \
> export GOBIN="$GOPATH/bin" \
> export PATH="$PATH:$GOBIN"

### Build Goquette-Binary with Go
Make sure you have a working installation of Go. Its easy to set up - just follow the official documentations. \
Fork this repository and run `go build .` inside this project folder.

### Pre-Compiled Binary
Download the pre-compiled binary of **Goquette** from the 'Release'-Page on this **[GitHub Page](https://github.com/PatrickLaabs/goquette/releases)**.

* On Windows: \
Move the extraced binary to a folder of your choice and put the path to **Goquette** into the Machine's PATH.
* On Linux & Darwin(macOS): \
Move the extraced binary to `/usr/local/bin`, check for permissions,
and add the path inside your $PATH \
`export PATH=$HOME/bin:/usr/local/bin:$PATH`

## Contribution

Hope you like this project. \
Every contribution is appreciated - feel free to use it in your project, fork it, modify it. Whatever you like :) 

If you experience any issues during setup / running **Goquette**: \
Open an Issue and let me know what's not working for you.

##  Example config.yaml

> id: "<your_project_name>" \
> version: "<your_software_version>" \
> title: "<set_a_title>" \
> authors: "<who's_the_author>" \
> owners: "<who's_the_owner>" \
> requireLicenseAcceptance: "<choose_true_or_false>" \
> description: "<set_a_description>" \
> summary: "what_does_the_program_do" \
> tags: "<define_some_tags_for_chocolatey>" \
> zipPath : "<name_of_your_zipped_file_inside_tools_dir>"

## Example Tools Folder structure

> tools/<your_zipped_binary>.zip \
> tools/<chocolateyinstall.ps1> \
> tools/<chocolateyuninstall.ps1>