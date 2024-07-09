# TTS Savefile Bundler
Edit your Tabletop Simulator save files in your web browser. This _should_ be platform-agnostic, but I have only tested 
on my Windows 10 machine.

## Motivation
The purpose of this application is to solve 2 nagging problems:

1. Script management is cumbersome to maintain and present-day tools are not free of bugs
2. Similarly to the point above, Card/Deck management offered by TTS is lackluster, and the Java-based Deck editor tool
   needs a feature uplift as well as stability improvements

## How it works

A web server is launched listening to port `3000`. You can open `http://localhost:3000` using the browser of your choice
and load/edit your save files from there.

## Features / Implementation roadmap

- [X] load tabletop savefile
- [ ] edit objects of the tabletop savefile
  - [ ] basic json editing
  - [ ] simplified version
- [ ] save tabletop files
  - [ ] edit the savefile json
- [ ] script hot reload / external editor api
  - [ ] synchronisation / notification of divergent states in your game/ide/savefile
- [ ] script bundling / packaging tool
  - [ ] enable creation of resuable APIs
  - [ ] classloader for reducing size of tabletop savefiles

## Build steps
**Prerequsistes**
- `npm` Node Package Manager ^20.2.0
- `go` Go Programming language v1.22 or later

**Install Script (basic, I recommend you to adjust it for your needs)**
```shell
#!/bin/sh
# run in your project root

# build Angular frontend
cd frontend
npm i
npm run build
# move compiled data to backend's static resources
cp -r dist/* ../res
cd ../backend
# build Go backend
# make sure you're building for the right platform (Windows/Linux/MacOS)
# this assumes you are not cross-compiling
go build -o app.exe ./backend/main/main.go
```


## Justification
- Editing cards / managing layout can be best leveraged by the powerful layout engine of the web browser's html renderer
- Last but not least, _I am a web developer, so it takes me less time to implement a functioning UI_
