

package zscratchpad


import "bytes"
import "encoding/json"
import "fmt"
import "net"
import "os"
import "sort"
import "strings"


import "github.com/jessevdk/go-flags"
import "github.com/pelletier/go-toml"




type GlobalFlags struct {
	Help *bool `long:"help" short:"h"`
	ConfigurationPath *string `long:"configuration" short:"c" value-name:"{configuration-path}"`
	WorkingDirectory *string `long:"chdir" short:"C" value-name:"{working-directory-path}"`
}

type GlobalConfiguration struct {
	WorkingDirectory *string `toml:"working_directory"`
}

type LibraryFlags struct {
	Paths []string `long:"library-path" value-name:"{library-path}"`
}

type ServerFlags struct {
	EndpointIp *string `long:"server-ip" value-name:"{ip}" toml:"endpoint_ip"`
	EndpointPort *uint16 `long:"server-port" value-name:"{port}" toml:"endpoint_port"`
}

type ListFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Type *string `long:"type" short:"t" choice:"library" choice:"document"`
	What *string `long:"what" short:"w" choice:"identifier" choice:"title" choice:"name" choice:"path"`
	Format *string `long:"format" short:"f" choice:"text" choice:"text-0" choice:"json"`
}

type SelectFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Type *string `long:"type" short:"t" choice:"library" choice:"document"`
	What *string `long:"what" short:"w" choice:"identifier" choice:"title" choice:"name" choice:"path"`
	How *string `long:"how" short:"W" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"body"`
	Format *string `long:"format" short:"f" choice:"text" choice:"text-0" choice:"json"`
}

type GrepFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	What *string `long:"what" short:"w" choice:"identifier" choice:"title" choice:"name" choice:"path"`
	Where *string `long:"where" short:"W" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"body"`
	Format *string `long:"format" short:"f" choice:"text" choice:"text-0" choice:"json" choice:"context"`
	Terms []string `long:"term" short:"t" value-name:"{term}"`
}

type ExportFlags struct {
	Document *string `long:"document" short:"d" required:"-" value-name:"{identifier}"`
	Format *string `long:"format" short:"f" choice:"source" choice:"text" choice:"html"`
}

type EditFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Document *string `long:"document" short:"d" value-name:"{identifier}"`
	Select *bool `long:"select" short:"s"`
}

type CreateFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Document *string `long:"document" short:"d" value-name:"{identifier}"`
	Select *bool `long:"select" short:"s"`
}

type DumpFlags struct {}

type MainFlags struct {
	Global *GlobalFlags `group:"Global options"`
	Library *LibraryFlags `group:"Library options"`
	List *ListFlags `command:"list"`
	Select *SelectFlags `command:"select"`
	Grep *GrepFlags `command:"grep"`
	Export *ExportFlags `command:"export"`
	Edit *EditFlags `command:"edit"`
	Create *CreateFlags `command:"create"`
	Server *ServerFlags `command:"server"`
	Dump *DumpFlags `command:"dump"`
}


type MainConfiguration struct {
	Global *GlobalConfiguration `toml:"globals"`
	Libraries []Library `toml:"library"`
	Server *ServerFlags `toml:"server"`
}




func Main (_executable string, _arguments []string, _environment map[string]string) (*Error) {
	
	_globals, _error := GlobalsNew (_executable, _environment)
	if _error != nil {
		return _error
	}
	
	_flags := & MainFlags {
			Global : & GlobalFlags {},
			Library : & LibraryFlags {},
			List : & ListFlags {},
			Select : & SelectFlags {},
			Grep : & GrepFlags {},
			Export : & ExportFlags {},
			Edit : & EditFlags {},
			Create : & CreateFlags {},
			Server : & ServerFlags {},
			Dump : & DumpFlags {},
		}
	
	_parser := flags.NewNamedParser ("z-scratchpad", flags.PassDoubleDash)
	_parser.SubcommandsOptional = true
	if _, _error := _parser.AddGroup ("", "", _flags); _error != nil {
		return errorw (0x5b48e356, _error)
	}
	
	_help := func (_log bool, _error *Error) (*Error) {
		_buffer := bytes.NewBuffer (nil)
		_parser.WriteHelp (_buffer)
		if _log {
			if _globals.StdioIsTty && _globals.TerminalEnabled {
				logf ('`', 0xa725b4bc, "\n%s\n", _buffer.String ())
			}
		} else {
			if _, _error := _buffer.WriteTo (_globals.Stdout); _error != nil {
				return errorw (0xf4170873, _error)
			}
		}
		return _error
	}
	
	// FIXME:  The parser always uses the actual environment variables and not `_environment`!
	if _argumentsRest, _error := _parser.ParseArgs (_arguments); _error != nil {
		if flagBoolOrDefault (_flags.Global.Help, false) {
			return _help (false, nil)
		} else {
			return _help (true, errorw (0xa198fbfd, _error))
		}
	} else if len (_argumentsRest) != 0 {
		return _help (true, errorw (0x3c7b6224, nil))
	}
	
	if flagBoolOrDefault (_flags.Global.Help, false) {
		return _help (false, nil)
	}
	
	_configuration := & MainConfiguration {
			Global : & GlobalConfiguration {},
			Server : & ServerFlags {},
		}
	if _flags.Global.ConfigurationPath != nil {
		_path := *_flags.Global.ConfigurationPath
		if _path == "" {
			return errorw (0x9a6f64a7, nil)
		}
		_data, _error := os.ReadFile (_path)
		if _error != nil {
			return errorw (0xf2be5f5f, _error)
		}
		_buffer := bytes.NewBuffer (_data)
		_decoder := toml.NewDecoder (_buffer)
		_decoder.Strict (true)
		_error = _decoder.Decode (_configuration)
		if _error != nil {
			return errorw (0x93e9dab8, _error)
		}
	}
	
	if _parser.Active == nil {
		return _help (true, errorw (0x4cae2ee5, nil))
	}
	
	return MainWithFlags (_parser.Active.Name, _flags, _configuration, _globals)
}




func MainWithFlags (_command string, _flags *MainFlags, _configuration *MainConfiguration, _globals *Globals) (*Error) {
	
	if (_flags.Global.WorkingDirectory != nil) || (_configuration.Global.WorkingDirectory != nil) {
		_workingDirectory := flag2StringOrDefault (_flags.Global.WorkingDirectory, _configuration.Global.WorkingDirectory, "")
		if _workingDirectory == "" {
			return errorw (0xe7c58968, nil)
		}
		if _error := os.Chdir (_workingDirectory); _error != nil {
			return errorw (0x5aae8d30, _error)
		}
	}
	
	_index, _error := IndexNew (_globals)
	if _error != nil {
		return _error
	}
	
	_editor, _error := EditorNew (_globals, _index)
	if _error != nil {
		return _error
	}
	
	_error = MainLoadLibraries (_flags.Library, _configuration.Libraries, _globals, _index)
	if _error != nil {
		return _error
	}
	
	switch _command {
		
		case "list" :
			return MainList (_flags.List, _globals, _index)
		
		case "select" :
			return MainSelect (_flags.Select, _globals, _index, _editor)
		
		case "grep" :
			return MainGrep (_flags.Grep, _globals, _index, _editor)
		
		case "export" :
			return MainExport (_flags.Export, _globals, _index)
		
		case "edit" :
			return MainEdit (_flags.Edit, _globals, _index, _editor)
		
		case "create" :
			return MainCreate (_flags.Create, _globals, _index, _editor)
		
		case "server" :
			return MainServer (_flags.Server, _configuration.Server, _globals, _index, _editor)
		
		case "dump" :
			return MainDump (_flags.Dump, _globals, _index)
		
		default :
			return errorw (0xaca17bb9, nil)
	}
}




func MainExport (_flags *ExportFlags, _globals *Globals, _index *Index) (*Error) {
	
	if _flags.Document == nil {
		return errorw (0x1826914a, nil)
	}
	_document, _error := WorkflowDocumentResolve (*_flags.Document, _index)
	if _error != nil {
		return _error
	}
	
	_format := flagStringOrDefault (_flags.Format, "source")
	
	_buffer := (*bytes.Buffer) (nil)
	switch _format {
		
		case "source" :
			if _output, _error := DocumentRenderToSource (_document); _error == nil {
				_buffer = bytes.NewBufferString (_output)
			} else {
				return _error
			}
		
		case "html" :
			if _output, _error := DocumentRenderToHtml (_document); _error == nil {
				_buffer = bytes.NewBufferString (_output)
			} else {
				return _error
			}
		
		case "text" :
			if _output, _error := DocumentRenderToText (_document); _error == nil {
				_buffer = bytes.NewBufferString (_output)
			} else {
				return _error
			}
		
		default :
			return errorw (0x326240d3, nil)
	}
	
	if _, _error := _buffer.WriteTo (_globals.Stdout); _error != nil {
		return errorw (0xa797b17f, _error)
	}
	
	return nil
}




func MainEdit (_flags *EditFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_flagSelect := flagBoolOrDefault (_flags.Select, false)
	if _flagSelect && (_flags.Document != nil) {
		return errorw (0x17114913, nil)
	}
	
	_identifier := ""
	if _flagSelect {
		
		_libraryIdentifier := flagStringOrDefault (_flags.Library, "")
		_options, _error := mainListOptionsAndSelect (_libraryIdentifier, "document", "title", "identifier", _index, _editor)
		if _error != nil {
			return _error
		}
		switch len (_options) {
			case 0 :
				return errorw (0x2b83fe67, nil)
			case 1 :
				_identifier = _options[0][1]
			default :
				return errorw (0x979cfe4c, nil)
		}
		
	} else {
		
		if _library, _document, _error := mainMergeLibraryAndDocumentIdentifiers (_flags.Library, _flags.Document); _error != nil {
			return _error
		} else if _document != "" {
			_identifier = _document
		} else if _library != "" {
			return errorw (0xdbe83c6c, nil)
		} else {
			return errorw (0xa0fc749c, nil)
		}
	}
	
	return WorkflowDocumentEdit (_identifier, _index, _editor, true)
}




func MainCreate (_flags *CreateFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_flagSelect := flagBoolOrDefault (_flags.Select, false)
	if _flagSelect && (_flags.Document != nil) {
		return errorw (0x2a0a4328, nil)
	}
	if _flagSelect && (_flags.Library != nil) {
		return errorw (0x4d3444df, nil)
	}
	
	_identifier := ""
	if _flagSelect {
		
		_options, _error := mainListOptionsAndSelect ("", "library", "title", "identifier", _index, _editor)
		if _error != nil {
			return _error
		}
		switch len (_options) {
			case 0 :
				return errorw (0x29abcd02, nil)
			case 1 :
				_identifier = _options[0][1]
			default :
				return errorw (0x22d4ddbe, nil)
		}
		
	} else {
		
		if _library, _document, _error := mainMergeLibraryAndDocumentIdentifiers (_flags.Library, _flags.Document); _error != nil {
			return _error
		} else if _document != "" {
			_identifier = _document
		} else if _library != "" {
			_identifier = _library
		} else {
			return errorw (0x22cc7dea, nil)
		}
	}
	
	return WorkflowDocumentCreate (_identifier, _index, _editor, true)
}




func MainList (_flags *ListFlags, _globals *Globals, _index *Index) (*Error) {
	
	_libraryIdentifier := flagStringOrDefault (_flags.Library, "")
	_type := flagStringOrDefault (_flags.Type, "document")
	_what := flagStringOrDefault (_flags.What, "identifier")
	_format := flagStringOrDefault (_flags.Format, "text")
	
	_options, _error := mainListOptions (_libraryIdentifier, _type, "identifier", _what, _index)
	if _error != nil {
		return _error
	}
	
	return mainListOutput (_options, _format, _globals)
}


func MainSelect (_flags *SelectFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_libraryIdentifier := flagStringOrDefault (_flags.Library, "")
	_type := flagStringOrDefault (_flags.Type, "document")
	_what := flagStringOrDefault (_flags.What, "identifier")
	_how := flagStringOrDefault (_flags.How, "title")
	_format := flagStringOrDefault (_flags.Format, "text")
	
	_options, _error := mainListOptionsAndSelect (_libraryIdentifier, _type, _how, _what, _index, _editor)
	if _error != nil {
		return _error
	}
	
	return mainListOutput (_options, _format, _globals)
}


func mainListOptionsAndSelect (_libraryIdentifier string, _type string, _labelSource string, _valueSource string, _index *Index, _editor *Editor) ([][2]string, *Error) {
	
	_options, _error := mainListOptions (_libraryIdentifier, _type, _labelSource, _valueSource, _index)
	if _error != nil {
		return nil, _error
	}
	
	_selection, _error := mainListSelect (_options, _editor)
	if _error != nil {
		return nil, _error
	}
	
	return _selection, nil
}


func mainListOptions (_libraryIdentifier string, _type string, _labelSource string, _valueSource string, _index *Index) ([][2]string, *Error) {
	
	_library := (*Library) (nil)
	if _libraryIdentifier != "" {
		if _library_0, _error := WorkflowLibraryResolve (_libraryIdentifier, _index); _error == nil {
			_library = _library_0
		} else {
			return nil, errorw (0x5a3e46e1, nil)
		}
	}
	
	_options := make ([][2]string, 0, 1024)
	
	switch _type {
		
		case "libraries", "library" :
			
			_libraries := []*Library (nil)
			if _library != nil {
				_libraries = []*Library { _library }
			} else {
				if _libraries_0, _error := IndexLibrariesSelectAll (_index); _error == nil {
					_libraries = _libraries_0
				} else {
					return nil, _error
				}
			}
			
			for _, _library := range _libraries {
				
				_label := ""
				_labels := make ([]string, 0, 16)
				switch _labelSource {
					case "identifier" :
						_label = _library.Identifier
					case "title", "name" :
						_label = _library.Name
						if _label == "" {
							_label = "[" + _library.Identifier + "]"
						}
					case "path" :
						_labels = _library.Paths
					case "body" :
						return nil, errorw (0x6aaf334b, nil)
					default :
						return nil, errorw (0xf0f17afb, nil)
				}
				if _label != "" {
					_labels = append (_labels, _label)
				}
				
				_value := ""
				_values := make ([]string, 0, 16)
				switch _valueSource {
					case "identifier" :
						_value = _library.Identifier
					case "title", "name" :
						_value = _library.Name
						if _value == "" {
							_value = "[" + _library.Identifier + "]"
						}
					case "path" :
						_values = _library.Paths
					case "body" :
						return nil, errorw (0xabd3314f, nil)
					default :
						return nil, errorw (0x4fab7acb, nil)
				}
				if _value != "" {
					_values = append (_values, _value)
				}
				
				for _, _label := range _labels {
					if _label == "" {
						continue
					}
					for _, _value := range _values {
						if _value == "" {
							continue
						}
						_options = append (_options, [2]string { _label, _value })
					}
				}
			}
		
		case "documents", "document" :
			
			_documents := []*Document (nil)
			if _library != nil {
				if _documents_0, _error := IndexDocumentsSelectInLibrary (_index, _library.Identifier); _error == nil {
					_documents = _documents_0
				} else {
					return nil, _error
				}
			} else {
				if _documents_0, _error := IndexDocumentsSelectAll (_index); _error == nil {
					_documents = _documents_0
				} else {
					return nil, _error
				}
			}
			
			for _, _document := range _documents {
				
				_label := ""
				_labels := make ([]string, 0, 16)
				switch _labelSource {
					case "identifier" :
						_label = _document.Identifier
					case "title", "name" :
						_label = _document.Title
						if _label == "" {
							_label = "[" + _document.Identifier + "]"
						}
						for _, _title := range _document.TitleAlternatives {
							if _title != _label {
								_labels = append (_labels, _title)
							}
						}
					case "path" :
						_label = _document.Path
					case "body" :
						_labels = make ([]string, 0, 1024)
						for _, _line := range _document.BodyLines {
							if stringTrimSpaces (_line) != "" {
								_labels = append (_labels, _line)
							}
						}
					default :
						return nil, errorw (0x9f3c1037, nil)
				}
				if _label != "" {
					_labels = append (_labels, _label)
				}
				
				_value := ""
				_values := make ([]string, 0, 16)
				switch _valueSource {
					case "identifier" :
						_value = _document.Identifier
					case "title", "name" :
						_value = _document.Title
						if _value == "" {
							_value = "[" + _document.Identifier + "]"
						}
						_values = make ([]string, 0, 16)
						for _, _title := range _document.TitleAlternatives {
							if _title != _value {
								_values = append (_values, _title)
							}
						}
					case "path" :
						_value = _document.Path
					case "body" :
						_values = make ([]string, 0, 1024)
						for _, _line := range _document.BodyLines {
							if stringTrimSpaces (_line) != "" {
								_values = append (_values, _line)
							}
						}
					default :
						return nil, errorw (0x2f341212, nil)
				}
				if _value != "" {
					_values = append (_values, _value)
				}
				
				for _, _label := range _labels {
					if _label == "" {
						continue
					}
					for _, _value := range _values {
						if _value == "" {
							continue
						}
						_options = append (_options, [2]string { _label, _value })
					}
				}
			}
		
		default :
			return nil, errorw (0x2c37fb9c, nil)
	}
	
	return _options, nil
}


func mainListSelect (_options [][2]string, _editor *Editor) ([][2]string, *Error) {
	
	_labels := make ([]string, 0, len (_options))
	_values := make (map[string]map[string]bool, len (_options))
	for _, _option := range _options {
		_label := _option[0]
		_value := _option[1]
		_label = stringTrimSpaces (_label)
		_values_1 := map[string]bool (nil)
		if _values_0, _exists := _values[_label]; _exists {
			_values_1 = _values_0
		} else {
			_labels = append (_labels, _label)
			_values_1 = make (map[string]bool, 16)
			_values[_label] = _values_1
		}
		_values_1[_value] = true
	}
	
	sort.Strings (_labels)
	
	_selection_0, _error := EditorSelect (_editor, _labels)
	if _error != nil {
		return nil, _error
	}
	
	_selection := make ([][2]string, 0, 16)
	for _, _label := range _selection_0 {
		if _values_0, _exists := _values[_label]; _exists {
			for _value, _ := range _values_0 {
				_selection = append (_selection, [2]string { _label, _value })
			}
		} else {
			return nil, errorw (0xdbff774c, nil)
		}
	}
	
	return _selection, nil
}


func mainListOutput (_options [][2]string, _format string, _globals *Globals) (*Error) {
	
	_list := make ([]string, 0, len (_options))
	_listSet := make (map[string]bool, len (_options))
	for _, _option := range _options {
		_value := _option[1]
		if _, _exists := _listSet[_value]; _exists {
			continue
		}
		_list = append (_list, _value)
		_listSet[_value] = true
	}
	
	sort.Strings (_list)
	
	_buffer := bytes.NewBuffer (nil)
	
	switch _format {
		
		case "text", "text-0" :
			_separator := byte ('\n')
			if _format == "text-0" {
				_separator = 0
			}
			for _, _value := range _list {
				_buffer.WriteString (_value)
				_buffer.WriteByte (_separator)
			}
		
		case "json" :
			_encoder := json.NewEncoder (_buffer)
			if _error := _encoder.Encode (_list); _error != nil {
				return errorw (0xc65a050c, _error)
			}
		
		default :
			return errorw (0x4def007c, nil)
	}
	
	if _, _error := _buffer.WriteTo (_globals.Stdout); _error != nil {
		return errorw (0xcf76965f, _error)
	}
	
	return nil
}




func MainGrep (_flags *GrepFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_libraryIdentifier := flagStringOrDefault (_flags.Library, "")
	_what := flagStringOrDefault (_flags.What, "identifier")
	_where := flagStringOrDefault (_flags.Where, "title")
	_format := flagStringOrDefault (_flags.Format, "text")
	
	_terms := make ([]string, 0, len (_flags.Terms))
	for _, _term := range _flags.Terms {
		if _term == "" {
			continue
		}
		_terms = append (_terms, _term)
	}
	if len (_terms) == 0 {
		return errorw (0xa95cd520, nil)
	}
	
	_options, _error := mainListOptions (_libraryIdentifier, "document", _where, _what, _index)
	if _error != nil {
		return _error
	}
	
	_selection := make ([][2]string, 0, len (_options) / 2)
	for _, _option := range _options {
		_contents := _option[0]
		_matched := false
		if !_matched {
			for _, _term := range _terms {
				if strings.Index (_contents, _term) != -1 {
					_matched = true
					break
				}
			}
		}
		if _matched {
			_selection = append (_selection, _option)
		}
	}
	
	return mainListOutput (_selection, _format, _globals)
}




func MainServer (_flags *ServerFlags, _configuration *ServerFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_endpointIp := flag2StringOrDefault (_flags.EndpointIp, _configuration.EndpointIp, "127.13.160.195")
	_endpointPort := flag2Uint16OrDefault (_flags.EndpointPort, _configuration.EndpointPort, 8080)
	
	_endpoint := fmt.Sprintf ("%s:%d", _endpointIp, _endpointPort)
	
	logf ('i', 0x210494be, "[server]  listening on `%s`...", _endpoint)
	
	_listener, _error_0 := net.Listen ("tcp", _endpoint)
	if _error_0 != nil {
		return errorw (0xedeea766, _error_0)
	}
	
	_globals.TerminalEnabled = false
	
	_server, _error := ServerNew (_globals, _index, _editor, _listener)
	if _error != nil {
		return _error
	}
	
	_error = ServerRun (_server)
	if _error != nil {
		return _error
	}
	
	return nil
}




func MainDump (_flags *DumpFlags, _globals *Globals, _index *Index) (*Error) {
	
	_documents, _error := IndexDocumentsSelectAll (_index)
	if _error != nil {
		return _error
	}
	
	_buffer := bytes.NewBuffer (nil)
	for _, _document := range _documents {
		_buffer.WriteString ("\n")
		_error = DocumentDump (_buffer, _document, true, false, false)
		if _error != nil {
			return _error
		}
		_buffer.WriteString ("\n")
	}
	
	if _, _error := _buffer.WriteTo (_globals.Stdout); _error != nil {
		return errorw (0xbf6a449c, _error)
	}
	
	return nil
}




func MainLoadLibraries (_flags *LibraryFlags, _configuration []Library, _globals *Globals, _index *Index) (*Error) {
	
	if (len (_flags.Paths) > 0) && (len (_configuration) > 0) {
		return errorw (0x374ece0f, nil)
	}
	
	_libraries := make ([]*Library, 0, 16)
	
	if len (_flags.Paths) > 0 {
		_library := & Library {
				Identifier : "library",
				Name : "Library",
				Paths : _flags.Paths,
				UseFileNameAsIdentifier : false,
				UseFileExtensionAsFormat : true,
				IncludeGlobPatterns : []string { "**/*.{txt,md}" },
				EditEnabled : true,
				CreateEnabled : true,
				CreatePath : _flags.Paths[0],
				SnapshotEnabled : true,
			}
		_libraries = append (_libraries, _library)
	}
	
	if len (_configuration) > 0 {
		for _, _library_0 := range _configuration {
			_library := & Library {}
			*_library = _library_0
			_libraries = append (_libraries, _library)
		}
	}
	
	if len (_libraries) == 0 {
		return errorw (0x00ea182b, nil)
	}
	
	for _, _library := range _libraries {
		if _error := LibraryInitialize (_library); _error != nil {
			return _error
		}
	}
	
	for _, _library := range _libraries {
		
		_error := IndexLibraryInclude (_index, _library)
		if _error != nil {
			return _error
		}
		
		_documentPaths, _error := libraryDocumentsWalk (_library)
		if _error != nil {
			return _error
		}
		
		_documents, _error := libraryDocumentsLoad (_library, _documentPaths)
		if _error != nil {
			return _error
		}
		
		for _, _document := range _documents {
			
			if _document.Library == "" {
				_document.Library = _library.Identifier
			}
			
			_document.EditEnabled = _library.EditEnabled
			
			_error = DocumentInitializeIdentifier (_document, _library)
			if _error != nil {
				return _error
			}
			
			_error = DocumentInitializeFormat (_document, _library)
			if _error != nil {
				return _error
			}
			
			_error = IndexDocumentInclude (_index, _document)
			if _error != nil {
				return _error
			}
		}
	}
	
	if true {
		_documents, _error := IndexDocumentsSelectAll (_index)
		if _error != nil {
			return _error
		}
		for _, _document := range _documents {
			_, _error = DocumentRenderToText (_document)
			if _error != nil {
				return _error
			}
			_, _error = DocumentRenderToHtml (_document)
			if _error != nil {
				return _error
			}
		}
	}
	
	return nil
}




func mainMergeLibraryAndDocumentIdentifiers (_library *string, _document *string) (string, string, *Error) {
	
	if _library != nil {
		
		if _document != nil {
			if _identifier, _error := DocumentFormatIdentifier (*_library, *_document); _error == nil {
				return "", _identifier, nil
			} else {
				return "", "", _error
			}
		} else {
			if _identifier, _error := LibraryParseIdentifier (*_library); _error == nil {
				return _identifier, "", nil
			} else {
				return "", "", _error
			}
		}
		
	} else if _document != nil {
		
		if _identifier, _, _, _error := DocumentParseIdentifier (*_document); _error == nil {
			return "", _identifier, nil
		} else {
			return "", "", _error
		}
		
	} else {
		
		return "", "", nil
	}
}




func flagBoolOrDefault (_value *bool, _default bool) (bool) {
	if _value != nil {
		return *_value
	}
	return _default
}

func flagUint16OrDefault (_value *uint16, _default uint16) (uint16) {
	if _value != nil {
		return *_value
	}
	return _default
}

func flagStringOrDefault (_value *string, _default string) (string) {
	if _value != nil {
		return *_value
	}
	return _default
}


func flag2BoolOrDefault (_value_1 *bool, _value_2 *bool, _default bool) (bool) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

func flag2Uint16OrDefault (_value_1 *uint16, _value_2 *uint16, _default uint16) (uint16) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

func flag2StringOrDefault (_value_1 *string, _value_2 *string, _default string) (string) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

