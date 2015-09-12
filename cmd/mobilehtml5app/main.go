package main /*
Command mobilehtml5app generates a simple framework to develop Go language
mobile applications with HTML5 based frontends using WebViews and a golang
backend. It currently supports Android only.

Usage

First create a folder within your GOPATH where you want your project to reside
and chdir into it. From within this folder run the mobilehtml5app command
as discussed in the platform specific sections below. This generates a
Go language HTTP server scaffolding for your app and a mobile platform specific
App project with a WebView that loads webpages from the server.

Under the hoods, the command will generate a file webapp.go that exports Start()
and Stop() functions to start and stop the backend server which are called by
lifecycle function hooks of the native portion of the App which houses the webview
The file webapp.go also will have two sample handlers to illustrate how to
create and register your HTTP handlers.

The webapp uses an server that integrates graceful shutdown and parameterized routing. It
requires handlers to satisfy the ContextHandler interface similar to http.Handler but
taking a context.Context as the first parameter. Server shutdown is signaled (when
the user closes the app etc.) by the Done() channel in the Context being closed
and handlers that spawn long-running processes should check for it. Named routing
parameters and any custom server instance specific settings are also passed
as through the Context and can be accessed via Context.Value(). For more details
on the server see http://godoc.org/github.com/srinathh/mobilehtml5app/server

You may want to set the environment variable $GO15VENDOREXPERIMENT=1 to use
the vendored versions of the packages github.com/julienschmidt/httprouter and
github.com/tylerb/graceful which are used in the Server.

Android apps

To create an Android project run the following command in the project folder
you create for your mobile app under your GOPATH.

	mobilehtml5app -apitarget <Build API Target> -name <Project Name>

This will generate webapp.go and an Android gradle based project in a subfolder
called androidapp. You can build the Android project thorugh the command line. To
work with it it in Android Studio, make sure to select "Import Project" in the
first screen rather than "Open Project".

There are two options for the WebView - the Android System WebView or the
Apache CrossWalk project XWalkView. The Android System WebView is a reliable
HTML5 platform only if you are targeting Android Kit-Kat (4.4) or higher devices
in which the WebView is based on Chromium. The CrossWalk project XWalkView has
compatibility from Android Ice Cream Sandwitch (4.0) version onwards and is the
default version used.

The full set of command line options for building Android apps are:
	-apitarget string
		Required. Android build target. To list possible targets run
		$ANDROID_HOME/tools/android list targets
	-gradle string
		Optional. Gradle version. (default "2.4")
	-name string
		Required. Android project name composed of a-z A-Z 0-9 _
	-plugin string
		Optional. Android gradle plugin version. (default "1.3.0")
	-target string
		Optional. Supports only android for now. (default "android")
	-title string
		Optional. App Title defaults to -name if omitted.
	-webview string
		Optional. Can be either xwalk (to use CrossWalk XWalkView) or system
		(to use Android WebView). Note that only KitKat (API19) and above
		have a system WebView supporting modern HTML5 capabilities based on
		Chromium and we set 19 as the minSdkVersion. (default "xwalk")
*/
import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var target, androidToolPath, name, apitarget, gradle, plugin, title, webview, pkgPath, pkgName, outPath string

func init() {
	flag.StringVar(&target, "target", "android", "Optional. Supports only android for now.")
	flag.StringVar(&name, "name", "", "Required. Android project name composed of a-z A-Z 0-9 _")
	flag.StringVar(&apitarget, "apitarget", "", "Required. Android build target. To list possible targets run $ANDROID_HOME/tools/android list targets")
	flag.StringVar(&gradle, "gradle", "2.4", "Optional. Gradle version.")
	flag.StringVar(&plugin, "plugin", "1.3.0", "Optional. Android gradle plugin version.")
	flag.StringVar(&title, "title", "", "Optional. App Title defaults to -name if omitted.")
	flag.StringVar(&webview, "webview", "xwalk", "Optional. Can be either xwalk (to use CrossWalk Webview) or system (to use Android WebView). Note that only KitKat (API19) and above have a system WebView supporting modern HTML5 capabilities based on Chromium and we set 19 as the minSdkVersion.")
}

func exitError(err error) {
	log.Printf("Error: %s\n", err)
	log.Fatalln("For usage details, run mobilehtml5app -help")
	log.Fatalln()
}

type modBytes func([]byte) ([]byte, error)

func modFile(fpath string, modfn modBytes) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		exitError(fmt.Errorf("unable to read %s: %s", fpath, err))
	}
	w, err := modfn(b)
	if err != nil {
		exitError(fmt.Errorf("error processing %s: %s", fpath, err))
	}
	if err := ioutil.WriteFile(fpath, w, 0644); err != nil {
		exitError(fmt.Errorf("umable to write changes to build.gradle: %s", err))
	}
}

func genAndroid() {
	// Check for existance of the android tool in the sdk path
	toolPath := filepath.Join(os.Getenv("ANDROID_HOME"), "tools", "android")
	if runtime.GOOS == "windows" {
		toolPath = toolPath + ".bat"
	}
	if _, err := os.Stat(toolPath); err != nil {
		exitError(fmt.Errorf("couldn't find %s. Environment variable ANDROID_HOME must be defined and point to a valid Android SDK folder", toolPath))
	}

	// -name and -apitarget must be provided
	if name == "" {
		exitError(fmt.Errorf("-name must be specified"))
	}

	if apitarget == "" {
		exitError(fmt.Errorf("-target must be specified"))
	}

	if !(webview == "xwalk" || webview == "system") {
		exitError(fmt.Errorf("-webview must be either xwalk or system. Got %s", webview))
	}

	// figure out the package name from the current working directory
	// and set the import path for the java app
	wd, err := os.Getwd()
	if err != nil {
		exitError(fmt.Errorf("couldn't fetch working dir path: %s", err))
	}

	pkgName = filepath.Base(wd)
	pkg, err := build.ImportDir(wd, build.FindOnly)
	if err != nil || pkg.ImportPath == "" {
		exitError(fmt.Errorf("could not derive package path from import path"))
	}
	pkgPath = javaImportPath(pkg.ImportPath) + ".androidapp"

	if title == "" {
		title = name
	}
	outPath = "./androidapp"

	// Project Generation: Run the android tool and then generate the go files
	{
		if out, err := exec.Command(toolPath, "create", "project", "--name", name, "--package", pkgPath, "--activity", "Main", "--target", apitarget, "--gradle", "--gradle-version", plugin, "--path", outPath).CombinedOutput(); err != nil {
			exitError(fmt.Errorf("error in android tool: %s\n%s", err, out))
		}
		srccode, err := format.Source([]byte(fmt.Sprintf("package %s\n%s", pkgName, webapp)))
		if err != nil {
			exitError(fmt.Errorf("gofmt error : %s", err))
		}

		if err := ioutil.WriteFile("webapp.go", srccode, 0644); err != nil {
			exitError(fmt.Errorf("error writing webapp.go: %s", err))
		}
	}

	// build.gradle: add dependencies for gomobile bind and XWalkView; fix minifyEnabled
	modFile(filepath.Join(outPath, "build.gradle"), buildDotGradle)

	// libs: create the folder
	if err := os.Mkdir(filepath.Join(outPath, "libs"), os.ModeDir|0775); err != nil {
		exitError(fmt.Errorf("unable to create libs folder: %s", err))
	}

	// gradle-wrapper.properties: update to a more modern gradle version
	modFile(filepath.Join(outPath, "gradle", "wrapper", "gradle-wrapper.properties"), gradlewrapperDotProperties)

	// AndroidManifest.xml: add permissions for INTERNET and NETWORK_STATE
	modFile(filepath.Join(outPath, "src", "main", "AndroidManifest.xml"), androidManifestDotXML)

	// strings.xml: modify the value of app_name to title
	modFile(filepath.Join(outPath, "src", "main", "res", "values", "strings.xml"), func(b []byte) ([]byte, error) {
		return bytes.Replace(b, []byte("Main"), []byte(title), 1), nil
	})

	// main.xml: remove the current layout file
	if err := os.Remove(filepath.Join(outPath, "src", "main", "res", "layout", "main.xml")); err != nil {
		exitError(fmt.Errorf("unable to remove main.xml: %s", err))
	}

	// finally write the new contents of Main.java
	fpath := filepath.Join(outPath, "src", "main", "java")
	for _, s := range strings.Split(pkgPath, ".") {
		fpath = filepath.Join(fpath, s)
	}
	fpath = filepath.Join(fpath, "Main.java")
	modFile(fpath, mainDotJava)

	// and generate backend.aar
	if out, err := exec.Command("gomobile", "bind", "-o", filepath.Join(outPath, "libs", "backend.aar"), ".").CombinedOutput(); err != nil {
		exitError(fmt.Errorf("could not generate libs/backend.aar: %s: %s", err, out))
	}

	// write a gitignore so we don't checkin local properties or build files
	if err := ioutil.WriteFile(filepath.Join(outPath, ".gitignore"), []byte(gitignore), 0644); err != nil {
		exitError(fmt.Errorf("error writing .gitignore:%s", err))
	}
}

func javaImportPath(goImportPath string) string {
	parts := strings.Split(goImportPath, "/")
	domainparts := strings.Split(parts[0], ".")
	ret := ""
	for _, domainpart := range domainparts {
		ret = domainpart + "." + ret
	}
	ret = ret[:len(ret)-1]

	for _, part := range parts[1:] {
		ret = ret + "." + part
	}
	return ret
}

func main() {
	flag.Parse()

	switch target {
	case "android":
		genAndroid()
	default:
		exitError(fmt.Errorf("-target must be android or ios. Received %s", target))
	}
}
