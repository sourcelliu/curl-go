package tool

import "strings"

// Basename is a translation of the C function `tool_basename` from
// curl-src/src/tool_bname.c.
//
// It returns the last component of a path, which is determined by the last
// occurrence of a forward or backward slash. If no slash is found, the
// original path is returned.
func Basename(path string) string {
	// Original C code logic from tool_bname.c, lines 36-37:
	//   s1 = strrchr(path, '/');
	//   s2 = strrchr(path, '\\');
	lastFwdSlash := strings.LastIndexByte(path, '/')
	lastBackSlash := strings.LastIndexByte(path, '\\')

	// Original C code logic from tool_bname.c, lines 39-46:
	//   if(s1 && s2) {
	//     path = (s1 > s2) ? s1 + 1 : s2 + 1;
	//   }
	//   else if(s1)
	//     path = s1 + 1;
	//   else if(s2)
	//     path = s2 + 1;
	lastSeparator := -1
	if lastFwdSlash > lastBackSlash {
		lastSeparator = lastFwdSlash
	} else {
		lastSeparator = lastBackSlash
	}

	if lastSeparator != -1 {
		// Return the substring after the last separator.
		return path[lastSeparator+1:]
	}

	// Original C code logic from tool_bname.c, line 48:
	//   return path;
	return path
}