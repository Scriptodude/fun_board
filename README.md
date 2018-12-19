# FUN BOARD

A simple client-server implementation of multiple board games made from scratch with as few libraries as possible - if not none - using GoLang as the backend server, Javascript for the frontend and HTML5's iframe as the viewport.

I define a Library as :

- Something not included de-facto in the language
- Some unofficial patch you can download to upgrade the language
- Any Framework is a set of libraries
- Any extension that requires download from a remote source (e.g JQuery or Bootstrap)

## How to execute
### Linux only at the moment

**Tested on Linux Mint 17**

1. [Install GoLang >= 1.11](https://golang.org/doc/install)
1. clone the repo on your $GOPATH _**Do NOT use `go install` !**_
	1. I haven't tested with `go install`, but i have a guess it may not work properly.
2. go to the root folder of fun_board (the one with .gitignore and stuff)
3. type in `chmod +x build.sh && ./build.sh`
4. `cd bin && ./server`
5. go to `localhost:8080`
6. Enjoy !