# Squad Up
Helps plan events with no worry

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
- Golang server provides API and servers static web page
- Static web page consumes API to access information

Squad Up's directory architecture is a bit more confusing:

- `/bower_components`
    - Frontend assets managed by [Bower](https://bower.io/)
    - Served under the `/lib` path by server
- `/client` - Frontend resources
    - `/components` 
        - Web components
        - Server under the `/components` path by server
    - `/css`
        - CSS resources
        - Served under the `/css` path by server
    - `/js`
        - JavaScript resources
        - Served under the `/js` path by server
    - `/views`
        - Client side views
    - `manifest.json` - Web App manifest
- `/server` - Server source
    - `/handlers` - HTTP route handlers
    - `/models` - "Models" / Any sort of data structure
        - `/db` - Models for Database
        - `/utils` - Helper utilities
    - `main.go` - Server entry point
- `/vendor`
    - Server libraries managed by [Glide](https://glide.sh/)
    
## Setup 
- 1. Create Database Docker container
    - If you have already completed this step once on your computer you 
      shouldn't have to run it again.
    - Run `make db-create`.
        - This command creates a Docker container running Postgres as our 
          database.
        - If the Docker container has already been created you will get an 
          error about the `The name "/squad-up-postgres"` already being taken.
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
        - This command starts the database Docker container we created previously.
- 2. Start Golang server.
    - You have to restart the server every single time you make a change to any 
      file that is not a static asset.
        - So changing pretty much anything in the `/client` directory doesn't 
          require a restart but anything is the `/server` directory does.
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
    - This is done in most terminals by using the key combo `CTRL + C`.