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

- [X] load TTS SaveData
- [X] view and search Decks, Bags, and Card Objects
- [X] edit/save Cards
  - [X] edit an existing Card
  - [X] create a new Card / remove a Card
- [ ] Cache currently open Save File
- [ ] script hot reload / external editor api
  - [ ] sync / notify of divergent state in game/ide/savefile
- [ ] script bundling / packaging tool
  - [ ] standard lua import semantics
- [ ] "Framework": ObjectPatcher / "Classloader" for reducing size of tabletop savefiles

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
rm -r ../res
cp -r dist/* ../res
cd ..
# build Go backend
# make sure you're building for the right platform (Windows/Linux/MacOS)
# this assumes you are not cross-compiling
go build -o app.exe ./backend/main/main.go
```
To run, simply execute the compiled binary `app.exe` in the root directory of the project.
```shell
@echo off
.\app.exe --browser=true # opens your default browser to the app's url
pause
```


## Justification
- Editing cards / managing layout can be best leveraged by the powerful layout engine of the web browser's html renderer
- Last but not least, _I am a web developer, so it takes me less time to implement a functioning UI_
