package filesystem

// Match checks if a file satisfies the given FileGuard by ensuring it is accepted and not rejected.
func FileGuardMatch(file File, gaurd FileGuard) bool {
	return gaurd.Accept(file) && !gaurd.Reject(file)
}

// MatchAll checks if a file satisfies all the provided FileGuards.
func FileGuardMatchAll(file File, gaurds []FileGuard) bool {
	for _, g := range gaurds {
		if !FileGuardMatch(file, g) {
			return false
		}
	}
	return true
}

// MatchAny checks if a file satisfies at least one of the provided FileGuards.
func FileGuardMatchAny(file File, gaurds []FileGuard) bool {
	for _, g := range gaurds {
		if FileGuardMatch(file, g) {
			return true
		}
	}
	return false
}

// AllMatch checks if all files in the list satisfy the given FileGuard.
func FileGuardAllMatch(files []File, gaurd FileGuard) bool {
	for _, file := range files {
		if !FileGuardMatch(file, gaurd) {
			return false
		}
	}
	return true
}

// FilterFiles filters a list of files based on the provided FileGuards.
// It returns a new slice containing only the files that match all the guards.
func FileGuardFilterFiles(files []File, guards []FileGuard) []File {
	var filteredFiles []File

	for _, file := range files {
		if FileGuardMatchAll(file, guards) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}
