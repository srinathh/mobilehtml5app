package main

import (
	"strings"
	"testing"
)

func TestMainDotJava(t *testing.T) {
	pkgName = "testapp"
	pkgPath = "com.example.testapp"
	ret, _ := mainDotJava(nil)
	if !strings.Contains(string(ret), "package com.example.testapp") {
		t.Error("Error in package statement")
	}

	if !strings.Contains(string(ret), "import go.testapp.Testapp") {
		t.Error("Error in import statement")
	}

	if !strings.Contains(string(ret), "Testapp.Start()") {
		t.Error("Error in webapp start")
	}

	if !strings.Contains(string(ret), "Testapp.Stop()") {
		t.Error("Error in webapp stop")
	}
}
