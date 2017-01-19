[![Build Status](https://travis-ci.org/Noah-Huppert/squad-up.svg?branch=master)](https://travis-ci.org/Noah-Huppert/squad-up)
[![Coverage Status](https://coveralls.io/repos/github/Noah-Huppert/squad-up/badge.svg?branch=master)](https://coveralls.io/github/Noah-Huppert/squad-up?branch=master)
# Squad Up
Helps plan events with no worry

# Contributing
The entire point of making this open source was so other people could contribute 
if they wanted. If you think of an idea for a feature or find a problem don't 
hesitate to open a new Github issue.  

## If you are looking for things to do
Check out the Github issue tracker and see 
if there are any issues labeled "help wanted" or "good for beginners". Issues with 
the "help wanted" label are tasks that need to be completed but don't have anyone 
working on them. Issues with the "good for beginners" label are tasks that are simple 
enough and do not require an in depth knowledge of systems architecture to complete.

Check out the Development section below to get started.

# Development
Requirements:
- [Docker](https://docker.com)
- [GNU Make](https://www.gnu.org/software/make/)
- [Golang](https://golang.org)
- [NodeJs (For Bower)](https://nodejs.org/en/)
    - [Bower](https://bower.io)
- [Glide](http://glide.sh/)

I am trying to move everything to Docker so you only need Docker and Make 
to develop Squad Up ([GH issue on move](https://github.com/Noah-Huppert/squad-up/issues/8)). 
But low priority until I start deploying Squad Up (Because "it works on my system" right now).

## Architecture
Squad Up's technical architecture is pretty simple:

- Postgres database stores information
- Golang server provides API and serves static web page
- Static web page consumes API to access information

Squad Up's directory architecture is a bit more confusing:

- `/bower_components`
    - Frontend assets managed by [Bower](https://bower.io/).
    - Served under the `/lib` path by server.
- `/client` - Frontend resources.
    - `/components` 
        - Web components source.
        - Served under the `/components` path by server.
    - `/css`
        - CSS resources.
        - Served under the `/css` path by server.
    - `/js`
        - JavaScript resources.
        - Served under the `/js` path by server.
    - `/views`
        - Client side views.
    - `manifest.json` - Web App manifest.
- `/server` - Server source.
    - `/handlers` - HTTP route handlers.
    - `/models` - "Models" / Any sort of data structure.
        - `/db` - Models for Database.
        - `/utils` - Helper utilities.
    - `main.go` - Server entry point.
- `/vendor`
    - Server libraries managed by [Glide](https://glide.sh/).
    
## Downloading
If you are familiar with Golang and the `GOPATH` then:  
Clone this repository down in `$GOPATH/src/github.com/Noah-Huppert/squad-up` and 
skip to the Setup section.  

If this is your first project using Golang and have no idea what the `GOPATH` is
then don't worry, I'll do my best to explain all that:  

Go requires that you put all your Go related projects in one directory 
following a specific pattern, like a workspace (Why it requires this is way beyond the 
scope of these docs). So lets make that directory:

- 1. Create a directory named `go` on your computer.
    - Place this directory where you normally place your projects.
- 2. In your terminal profile set the environment variable `GOPATH` to be the 
     path of the `go` folder.
- 3. This `go` directory will now be referred to as the `GOPATH`.

Further inside of the `GOPATH` Go requires that you put all your projects 
in the `$GOPATH/src` directory. Inside of this `$GOPATH/src` directory 
projects are stored similar to Java projects:  

- Each project has a collision resistant package name.
    - In Go these package names follow the scheme `github.com/Username/Project` 
      (Where as in Java they follow: `com.username.project`).
- Each "sub package" of the package name is a new directory.
    - In Go "sub packages" are separated by forward slashes (In Java they are 
      separated by dots).
    - So the package name `github.com/Username/Project` would yield the directory 
    structure of:
        - `/github.com`
            - `/Username`
                - `Project`

Squad Up's package name is `github.com/Noah-Huppert/squad-up` so it should be 
stored in the directory `$GOPATH/src/github.com/Noah-Huppert/squad-up`. So 
all you have to do now is clone this repo into that path.

## Setup 
- 1. Create Database Docker container.
    - If you have already completed this step once on your computer you 
      shouldn't have to do it again.
    - Run `make db-create`.
        - This command creates a Docker container running Postgres as our 
          database.
        - If the Docker container has already been created (aka., you already 
          did this step) you will get an error about the `The name 
          "/squad-up-postgres"` already being taken.
- 2. Install client side libraries with Bower.
    - You only have to run this command if the libraries listed in `bower.json` 
      change or on first time setup.
    - Run `bower install`.
        - This command installs the libraries listed in `bower.json` to `bower_components`.
- 3. Install Golang server libraries.
    - You only have to run this command if the libraries listed in `glide.yaml` change 
      or on first time setup.
    - Run `glide install`.
        - This installs all of our Golang libraries to the `vendor` directory.
        
## Running
- 1. Start Database Docker container.
    - If you just created the database container it should already be running 
      so you can skip this step.
    - Otherwise run `make db-start`.
        - This command starts the database Docker container we created previously 
           in Setup step 1.
- 2. Start Golang server.
    - You have to restart the server every single time you make a change to any 
      file that is not a static asset.
        - So changing pretty much anything in the `/client` directory doesn't 
          require a restart but anything in the `/server` directory does.
    - Run `make app-run`
        - This runs the Golang server from the entry point `/server/main.go`.
        
## Stopping
- 1. Stop Database Docker container.
    - Run `make db-stop`
        - This stops the Database Docker container we created in Setup step 1 
          and started in Running step 1.
    - If you would like to destroy the Database Docker container then run 
      `make db-destroy`.
        - This will delete the Database Docker container and all its data.
        - Note: To run Squad Up again you will have to complete Setup step 1 
          again.
- 2. Stopping Golang server
    - You can stop the server by sending the `SIGINT` signal.
    - This is done in most terminals by using the key combination `CTRL + C`.
