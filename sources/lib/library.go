

package zscratchpad


import "os"
import "path/filepath"
import "regexp"
import "strings"


import "github.com/gobwas/glob"




type Library struct {
	
	Identifier string
	Name string
	
	Paths []string
	
	EditEnabled bool
	
	CreateEnabled bool
	CreatePath string
	CreateExtension string
	
	SnapshotEnabled bool
	SnapshotExtension string
	
	IncludeGlobPatterns []string
	ExcludeGlobPatterns []string
	
	IncludeRegexPatterns []string
	ExcludeRegexPatterns []string
	
	UseFileNameAsIdentifier bool
	UseFileExtensionAsFormat bool
	
	includeGlobMatchers []glob.Glob
	excludeGlobMatchers []glob.Glob
	
	includeRegexMatchers []*regexp.Regexp
	excludeRegexMatchers []*regexp.Regexp
}




func LibraryInitialize (_library *Library) (*Error) {
	
	if _library.Identifier == "" {
		return errorw (0x94465013, nil)
	}
	
	for _index, _path := range _library.Paths {
		if _path == "" {
			return errorw (0x8b174330, nil)
		}
		if _path_0, _error := filepath.Abs (_path); _error == nil {
			_path = _path_0
		} else {
			return errorw (0xe0ece239, _error)
		}
		if _stat, _error := os.Stat (_path); _error == nil {
			if ! _stat.IsDir () {
				return errorw (0x410a4abd, nil)
			}
		} else {
			return errorw (0x1513652d, _error)
		}
		_library.Paths[_index] = _path
	}
	
	if _library.CreateEnabled {
		if _library.CreatePath == "" {
			if len (_library.Paths) == 1 {
				_library.CreatePath = _library.Paths[0]
			} else {
				return errorw (0xd76cef62, nil)
			}
		}
		if _path_0, _error := filepath.Abs (_library.CreatePath); _error == nil {
			_library.CreatePath = _path_0
		} else {
			return errorw (0xea573d9c, _error)
		}
		if _stat, _error := os.Stat (_library.CreatePath); _error == nil {
			if ! _stat.IsDir () {
				return errorw (0x1ad922a3, nil)
			}
		} else {
			return errorw (0x98ade3fc, _error)
		}
		if _library.CreateExtension == "" {
			_library.CreateExtension = "txt"
		}
		_library.CreateExtension = strings.TrimLeft (_library.CreateExtension, ".")
	} else {
		if _library.CreatePath != "" {
			return errorw (0x5b55e852, nil)
		}
		if _library.CreateExtension != "" {
			return errorw (0x2ffc3bf4, nil)
		}
	}
	
	if _library.SnapshotEnabled {
		if _library.SnapshotExtension == "" {
			_library.SnapshotExtension = "snapshot"
		}
		_library.SnapshotExtension = strings.TrimLeft (_library.SnapshotExtension, ".")
	} else {
		if _library.SnapshotExtension != "" {
			return errorw (0x3ede0dc5, nil)
		}
	}
	
	_library.includeGlobMatchers = make ([]glob.Glob, 0, len (_library.IncludeGlobPatterns))
	for _, _pattern := range _library.IncludeGlobPatterns {
		if _matcher, _error := glob.Compile (_pattern); _error == nil {
			_library.includeGlobMatchers = append (_library.includeGlobMatchers, _matcher)
		} else {
			return errorw (0x674d8ba9, _error)
		}
	}
	
	_library.excludeGlobMatchers = make ([]glob.Glob, 0, len (_library.IncludeGlobPatterns))
	for _, _pattern := range _library.ExcludeGlobPatterns {
		if _matcher, _error := glob.Compile (_pattern); _error == nil {
			_library.excludeGlobMatchers = append (_library.excludeGlobMatchers, _matcher)
		} else {
			return errorw (0x5d547147, _error)
		}
	}
	
	_library.includeRegexMatchers = make ([]*regexp.Regexp, 0, len (_library.IncludeRegexPatterns))
	for _, _pattern := range _library.IncludeRegexPatterns {
		if _matcher, _error := regexp.Compile (_pattern); _error == nil {
			_library.includeRegexMatchers = append (_library.includeRegexMatchers, _matcher)
		} else {
			return errorw (0x3515908f, _error)
		}
	}
	
	_library.excludeRegexMatchers = make ([]*regexp.Regexp, 0, len (_library.IncludeRegexPatterns))
	for _, _pattern := range _library.ExcludeRegexPatterns {
		if _matcher, _error := regexp.Compile (_pattern); _error == nil {
			_library.excludeRegexMatchers = append (_library.excludeRegexMatchers, _matcher)
		} else {
			return errorw (0xe3938785, _error)
		}
	}
	
	return nil
}




func libraryDocumentsLoad (_library *Library, _documentPaths []string) ([]*Document, *Error) {
	
	_documents := make ([]*Document, 0, len (_documentPaths))
	
	for _, _documentPath := range _documentPaths {
		if _document, _error := DocumentLoadFromPath (_documentPath); _error == nil {
			_documents = append (_documents, _document)
		} else {
			return nil, _error
		}
	}
	
	return _documents, nil
}




func libraryDocumentsWalk (_library *Library) ([]string, *Error) {
	
	_documentPaths := []string (nil)
	for _, _libraryPath := range _library.Paths {
		if _documentPaths_0, _error := libraryDocumentsWalkPath (_library, _libraryPath); _error == nil {
			if _documentPaths == nil {
				_documentPaths = _documentPaths_0
			} else {
				_documentPaths = append (_documentPaths, _documentPaths_0 ...)
			}
		} else {
			return nil, _error
		}
	}
	
	return _documentPaths, nil
}


func libraryDocumentsWalkPath (_library *Library, _libraryPath string) ([]string, *Error) {
	
	if _libraryPath == "" {
		return nil, errorw (0x83afc399, nil)
	}
	
	_snapshotSuffix := ""
	if _library.SnapshotEnabled && (_library.SnapshotExtension != "") {
		_snapshotSuffix = "." + _library.SnapshotExtension
	}
	
	_documentPaths := make ([]string, 0, 1024)
	_folderPaths := make ([]string, 0, 128)
	
	_walkFunc := func (_pathEntry string, _entry os.DirEntry) (*Error) {
		
//		logf ('d', 0x18d84756, "%s", _pathEntry)
		
		_name := _entry.Name ()
		
		if strings.HasPrefix (_name, ".") {
			if _entry.IsDir () {
//				logf ('d', 0xb53c5778, "%s", _pathEntry)
				return nil
			} else {
//				logf ('d', 0x63546fb2, "%s", _pathEntry)
				return nil
			}
		}
		
		_stat := os.FileInfo (nil)
		if _stat_0, _error := os.Stat (_pathEntry); _error == nil {
			_stat = _stat_0
		} else {
			return errorw (0xb00f4f21, _error)
		}
		
		_mode := _stat.Mode ()
		if _mode.IsRegular () {
			// NOP
		} else if _mode.IsDir () {
//			logf ('d', 0x47608981, "%s", _pathEntry)
			_folderPaths = append (_folderPaths, _pathEntry)
			return nil
		} else {
			return errorf (0xb0cc4319, "invalid entry `%s`", _pathEntry)
		}
		
		_pathRelative := ""
		if _pathRelative_0, _error := filepath.Rel (_libraryPath, _pathEntry); _error == nil {
			_pathRelative = "/" + _pathRelative_0
		} else {
			return errorw (0xacc84f2b, _error)
		}
		
		if _snapshotSuffix != "" {
			if strings.HasSuffix (_name, _snapshotSuffix) {
//				logf ('d', 0xeed5814c, "%s", _pathEntry)
				return nil
			}
		}
		
		_exclude := false
		if !_exclude {
			for _, _matcher := range _library.excludeGlobMatchers {
				if _matcher.Match (_pathRelative) {
					_exclude = true
					break
				}
			}
		}
		if !_exclude {
			for _, _matcher := range _library.excludeRegexMatchers {
				if _matcher.MatchString (_pathRelative) {
					_exclude = true
					break
				}
			}
		}
		if _exclude {
//			logf ('d', 0x71694f7f, "%s", _pathEntry)
			return nil
		}
		
		_include := false
		if !_include {
			for _, _matcher := range _library.includeGlobMatchers {
				if _matcher.Match (_pathRelative) {
					_include = true
					break
				}
			}
		}
		if !_include {
			for _, _matcher := range _library.includeRegexMatchers {
				if _matcher.MatchString (_pathRelative) {
					_include = true
					break
				}
			}
		}
		if !_include {
			if (len (_library.includeGlobMatchers) == 0) && (len (_library.includeRegexMatchers) == 0) {
				_include = true
			}
		}
		if !_include {
//			logf ('d', 0x3da79eb9, "%s", _pathEntry)
			return nil
		}
		
//		logf ('d', 0xaa73f1ac, "%s", _pathEntry)
		
		_documentPath := filepath.Join (_libraryPath, _pathRelative[1:])
		_documentPaths = append (_documentPaths, _documentPath)
		
		return nil
	}
	
	_folderPaths = append (_folderPaths, _libraryPath)
	
	for _folderIndex := 0; _folderIndex < len (_folderPaths); _folderIndex += 1 {
		_folderPath := _folderPaths[_folderIndex]
		_folderEntries, _error := os.ReadDir (_folderPath)
		if _error != nil {
			return nil, errorw (0x28422546, _error)
		}
		for _, _folderEntry := range _folderEntries {
			_folderEntryPath := filepath.Join (_folderPath, _folderEntry.Name ())
			if _error := _walkFunc (_folderEntryPath, _folderEntry); _error != nil {
				return nil, _error
			}
		}
	}
	
	return _documentPaths, nil
}




func LibraryValidateIdentifier (_identifier string) (*Error) {
	if ! LibraryIdentifierRegex.MatchString (_identifier) {
		return errorw (0x2d8a1040, nil)
	}
	return nil
}

func LibraryParseIdentifier (_identifier string) (string, *Error) {
	if _error := LibraryValidateIdentifier (_identifier); _error != nil {
		return "", _error
	}
	return _identifier, nil
}

var LibraryIdentifierRegexToken string = `(?:(?:[a-z0-9]+)(?:[_-]+[a-z0-9]+)*)`
var LibraryIdentifierRegex *regexp.Regexp = regexp.MustCompile (`^` + LibraryIdentifierRegexToken + `$`)

