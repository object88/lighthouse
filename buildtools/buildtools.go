// +build buildtools

package buildtools

// The purpose of this package and file is to build our tools in the vendor
// directory.
// See:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md

import (
	// The go-bindata program is used in the prebuild step for shipyard.
	_ "k8s.io/code-generator"
)