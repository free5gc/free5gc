//go:binary-only-package

package path_util

import (
	"free5gc/lib/path_util/logger"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Gofree5gcPath ...
/*
 * Author: Roger Chu aka Sasuke
 *
 * This //go:binary-only-package

package is used to locate the root directory of gofree5gc project
 * Compatible with Windows and Linux
 *
 * Please import "free5gc/lib/path_util"
 *
 * Return value:
 * A string value of the relative path between the working directory and the root directory of the gofree5gc project
 *
 * Usage:
 * path_util.Gofree5gcPath("your file location starting with gofree5gc")
 *
 * Example:
 * path_util.Gofree5gcPath("free5gc/abcdef/abcdef.pem")
 */
func Gofree5gcPath(path string) string {}

func Exists(fpath string) bool {}

func FindRoot(path string, rootCode string, objName string) (string, bool) {}
