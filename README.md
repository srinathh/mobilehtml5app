# mobilehtml5app
mobilehtml5app is simple framework and a tool for creating mobile apps with the Go language
using an HTML5 supporting WebView as the frontend and a Go language HTTP server as the
backend. This allows us to create front-ends for mobile apps using standard HTML5, CSS
and Javascript and make use of the plethora of standard and very mature Web and Web App
development frameworks available today rather than using OpenGL or the Native UI system.
It currently supports only the Android platform.

## Motivation 
Go Language has introduced support for mobile application development with version 1.5
and the mobile toolchain is rapidly developing. However, the current UI focus is primarily
on OpenGL.

On the other hand, both Android and iOS support WebView based applications that deliver
their frontend using HTML5. This allows us to build rich user interfaces taking advantage
of Web and Web App technologies and frameworks. On Android specifically, the Apache
CrossWalk and recent versions of Android system WebView (Kit Kat onwards) provide
this functionality. Also, the most frequent domain of usage of the Go Language is 
in writing Web Servers.

This project attempts to make it easy to use the mobile support in Go and its HTTP
server strengths in conjunction with HTML5 supporting WebViews to create mobile apps.

## Usage & Reference
- First use `go get github.com/srinathh/mobilehtml5app/...` to get the packages and the command
- Refer to [mobilehtml5app command](http://godoc.org/github.com/srinathh/mobilehtml5app/cmd/mobilehtml5app) for documentation on how to generate a mobile app project with a go HTTP server backend and HTML5 frontend.
- Refer to [server package](http://godoc.org/github.com/srinathh/mobilehtml5app/server) for documentation on the server used in the webapp that supports graceful restarts and parameterized routing

More documentation to come.
