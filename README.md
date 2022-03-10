# Goquette

With **Goquette** you can simplify your packaging experience with NuGet. \
Set up the required files - Steps are listed below - and just run **Goquette**, and you're done.

![Image alt text](images/goquette.png)

[![Go Reference](https://pkg.go.dev/badge/github.com/PatrickLaabs/goquette.svg)](https://pkg.go.dev/github.com/PatrickLaabs/goquette)
[![Go Report Card](https://goreportcard.com/badge/github.com/PatrickLaabs/goquette)](https://goreportcard.com/report/github.com/PatrickLaabs/goquette)

## Usage:

* Create a 'goquette.yaml' File inside your root project folder
* Create a 'tools' directory inside your root project folder

Inside the `tools`-Directory, put your PowerShell scripts, which are consumed by chocolatey, \
and your zipped program/binary.
Take a look inside the tools-Directory at this **[Tools-Dir of Goquette](https://github.com/PatrickLaabs/goquette/tree/main/tools)** \
for a better understanding on how to use **Goquette**.

## Installation

### Install via Go
```
go install github.com/PatrickLaabs/goquette@latest
```
Make sure you have set your GOBIN Path correctly. \
If not:
```
export GOBIN="$GOPATH/bin"
export PATH="$PATH:$GOBIN"
```
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

##  Example goquette.yaml

```
id: "<your_project_name>"
version: "<your_software_version>"
title: "<set_a_title>"
authors: "<who's_the_author>"
owners: "<who's_the_owner>"
requireLicenseAcceptance: "<choose_true_or_false>"
description: "<set_a_description>"
summary: "what_does_the_program_do"
tags: "<define_some_tags_for_chocolatey>"
zipPath : "<name_of_your_zipped_file_inside_tools_dir>"
```

---

## Example Tools Folder structure

```
tools/<your_zipped_binary>.zip
tools/<chocolateyinstall.ps1>
tools/<chocolateyuninstall.ps1>`
```

---

## Goquette and Jenkins

**Goquette** really shines, when used within a pipeline, e.g. Jenkins. \
The point here is, that you only need to configure your `goquette.yaml`, prepare your _powershell scripts_ \
and the rest is handled for you.

```
pipeline {
  agent any
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        GOBIN = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/bin"
        PATH = "${PATH}:${GOBIN}"
    }
    tools {
        go 'go-1.17.7'
    }

  stages {
    stage('Cleaning workspace'){
        steps {
            sh 'rm -rf $JENKINS_HOME/workspace/$JOB_NAME/content'
            sh 'rm -rf $JENKINS_HOME/workspace/$JOB_NAME/templates'
            sh 'rm -rf $JENKINS_HOME/workspace/$JOB_NAME/*.nupkg'
            sh 'rm -rf $JENKINS_HOME/workspace/$JOB_NAME/*.zip'
            sh 'rm -rf $JENKINS_HOME/workspace/$JOB_NAME/*.exe'
        }
    }

    stage('Build') {
      steps {
        sh 'GOOS=linux go build'
        sh 'GOOS=windows go build'
        sh 'tar cf $JENKINS_HOME/workspace/$JOB_NAME/tools/<your_application_name>.zip *.exe'
      }
    }

    stage('Goquette - NuGet Packaging') {
        steps {
            sh 'go install github.com/PatrickLaabs/goquette@latest'
            sh 'cd $JENKINS_HOME/workspace/$JOB_NAME && $JENKINS_HOME/jobs/$JOB_NAME/builds/$BUILD_ID/bin/goquette'
            }
        }

    stage('nFPM - rpm Packaging') {
        steps {
            sh 'go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest'
            sh '$JENKINS_HOME/jobs/$JOB_NAME/builds/$BUILD_ID/bin/nfpm pkg --packager rpm --target $JENKINS_HOME/workspace/$JOB_NAME/'
        }
    }

    stage('Deploy .nupkg to Nexus') {
        steps {
            echo 'deploying to nexus..'
            nexusArtifactUploader(
                nexusVersion: 'nexus2',
                protocol: 'http',
                nexusUrl: '<ip>:<port>/nexus',
                groupId: 'com.example',
                version: '<version>',
                repository: 'nuget',
                credentialsId: '<creds>',
                artifacts: [
                    [artifactId: '<project_name>',
                    classifier: 'release',
                     file: '$JENKINS_HOME/workspace/$JOB_NAME/<your_application_name>.nupkg',
                     type: 'nuget']
                ]
             )
        }
    }

    stage('Deploy .rpm to Nexus') {
        steps {
            nexusArtifactUploader(
                nexusVersion: 'nexus2',
                protocol: 'http',
                nexusUrl: '<ip>:<port>/nexus',
                groupId: 'com.example',
                version: '<version>',
                repository: 'rpm',
                credentialsId: 'nexus-user-credentials',
                artifacts: [
                    [artifactId: '<project_name>',
                     classifier: 'release',
                     file: '$JENKINS_HOME/workspace/$JOB_NAME/<your_application_name>.rpm',
                     type: 'rpm']
                ]
             )
        }
    }
  }
}
```

or a more simplistic approach:

```
pipeline {
    agent any
    
    stages {
        stage('Build and Packaging') {
            steps {
                script {
                    def root = tool type: 'go', name: 'go-1.17.7'
                    withEnv(["GOPATH=${root}", "PATH=${PATH}:${root}/bin"]) {
                       sh 'GOOS=linux go build'
                       sh 'GOOS=windows go build'
                       sh 'tar cf ./tools/<.zipName> *.exe'
                       sh 'go install github.com/PatrickLaabs/goquette@latest'
                       sh 'goquette'
                       sh 'go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest'
                       sh 'nfpm pkg --packager rpm --target ./'
                    }
                }
            }
        }
    }
}
```