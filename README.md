# Goquette

With this binary, you are able to package a NuGet Package without windows.
All you need is Goquette, and you are ready to take off.

## Usage:

* Create a 'config.yaml' File inside your root project folder
* Create a 'tools' directory inside your root projekt folder

###  Example config.yaml

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

### Example Tools Folder structure

> tools/<your_zipped_binary>.zip \
> tools/<chocolateyinstall.ps1> \
> tools/<chocolateyuninstall.ps1>