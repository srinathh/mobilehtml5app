# Golang Todo App

This app demonstrates how to use the mobilehtml5app framework to generate
a rich HTML5 based user interface on mobile for a Go Language webapp backend.
It also demonstrates how to use the private app storage space on Android
to persist data and how to link up all the build all components using Gradle.

<img src="https://github.com/srinathh/mobilehtml5app/raw/master/example/todoapp/screenshot.jpg" width="300">

# Web Frameworks Used
The app uses [React](https://facebook.github.io/react/) and [Bootstrap](http://getbootstrap.com/)
to build the user interface and interacts using AJAX calls with the Go Language
backend. It uses [BoltDB](https://github.com/boltdb/bolt) to persist the data
in the app's private persistant data folder which is obtained using
Activity.getFilesDir() and provided to the Start() function of the Go app.


## Requirements to Build
- [Node.js](https://nodejs.org/) and [Babel](https://babeljs.io/) installed
  to compile JSX to Javascript
- [Go-Bindata](https://github.com/jteeuwen/go-bindata) installed to compile assets into bindata.go

## How to Build
- In the Android Studio's opening dialog, select "Import Project"
- Import the androidapp folder

## Photo Credits
Markus Spiske, www.markusspiske.com
