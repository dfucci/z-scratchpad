

package zscratchpad


import "bytes"
import "encoding/json"
import "fmt"
import "net"
import "os"
import "path"
import "path/filepath"
import "runtime"
import "runtime/debug"
import "runtime/pprof"
import "strings"
import "time"


import "github.com/jessevdk/go-flags"
import "github.com/pelletier/go-toml"
import "github.com/akutz/sortfold"


import . "github.com/volution/z-scratchpad/embedded"




type GlobalFlags struct {
	Help *bool `long:"help" short:"h"`
	ConfigurationPath *string `long:"configuration" short:"c" value-name:"{configuration-path}"`
	WorkingDirectory *string `long:"chdir" short:"C" value-name:"{working-directory-path}"`
}

type GlobalConfiguration struct {
	UniqueIdentifier *string `toml:"unique_identifier"`
	WorkingDirectory *string `toml:"working_directory"`
	TerminalEnabled *bool `toml:"terminal_enabled"`
	XorgEnabled *bool `toml:"xorg_enabled"`
}

type IndexFlags struct {
	WalkDisabled *bool `long:"index-disable-walk"`
	DatabaseDisabled *bool `long:"index-disable-database"`
	LoadDisabled *bool `long:"index-disable-load"`
	StoreDisabled *bool `long:"index-disable-store"`
	DirtyDisabled *bool `long:"index-disable-dirty"`
	RefreshDisabled *bool `long:"index-disable-refresh"`
}

type IndexConfiguration struct {
	DatabaseEnabled *bool `toml:"database_enabled"`
	DatabasePath *string `toml:"database_path"`
	LibrariesRefreshEnabled *bool `toml:"libraries_refresh_enabled"`
	DocumentsRefreshEnabled *bool `toml:"documents_refresh_enabled"`
}

type LibraryFlags struct {
	Paths []string `long:"library-path" value-name:"{library-path}"`
}

type DebuggingFlags struct {
	ProfileCpuPath *string `long:"profile-cpu-path" value-name:"{path}"`
	ProfileMemoryPath *string `long:"profile-memory-path" value-name:"{path}"`
	ProfileLoops *uint `long:"profile-loops" value-name:"{count}"`
}

type EditorConfiguration struct {
	DefaultCreateLibrary *string `toml:"default_create_library"`
	TerminalEditCommand *[]string `toml:"terminal_edit_command"`
	XorgEditCommand *[]string `toml:"xorg_edit_command"`
	TerminalSelectCommand *[]string `toml:"terminal_select_command"`
	XorgSelectCommand *[]string `toml:"xorg_select_command"`
	TerminalClipboardStoreCommand *[]string `toml:"terminal_clipboard_store_command"`
	XorgClipboardStoreCommand *[]string `toml:"xorg_clipboard_store_command"`
}


type ListFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Type *string `long:"type" short:"t" choice:"library" choice:"document"`
	What *string `long:"what" short:"w" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"commonmark-link"`
	Format *string `long:"format" short:"f" choice:"text" choice:"text-sp" choice:"text-0" choice:"json"`
}

type SearchFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Type *string `long:"type" short:"t" choice:"library" choice:"document"`
	What *string `long:"what" short:"w" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"commonmark-link"`
	How *string `long:"how" short:"W" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"body"`
	Format *string `long:"format" short:"f" choice:"text" choice:"text-sp" choice:"text-0" choice:"json"`
	Action *string `long:"action" short:"a" choice:"output" choice:"edit" choice:"export" choice:"browse"`
	MultipleAllowed *bool `long:"multiple" short:"m"`
}

type GrepFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	What *string `long:"what" short:"w" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"commonmark-link"`
	Where *string `long:"where" short:"W" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"body"`
	Format *string `long:"format" short:"f" choice:"text" choice:"text-sp" choice:"text-0" choice:"json"`
	Terms []string `long:"term" short:"t" value-name:"{term}"`
	Action *string `long:"action" short:"a" choice:"output" choice:"edit" choice:"export" choice:"browse"`
	MultipleAllowed *bool `long:"multiple" short:"m"`
	MatchAny *bool `long:"match-any"`
}


type CreateFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Document *string `long:"document" short:"d" value-name:"{identifier}"`
	Select *bool `long:"select" short:"s"`
}

type EditFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Document *string `long:"document" short:"d" value-name:"{identifier}"`
	Select *bool `long:"select" short:"s"`
}

type ExportFlags struct {
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Document *string `long:"document" short:"d" value-name:"{identifier}"`
	Format *string `long:"format" short:"f" choice:"identifier" choice:"title" choice:"name" choice:"path" choice:"source" choice:"text" choice:"html" choice:"html-plain" choice:"html-github" choice:"html-github-auto" choice:"html-github-light" choice:"html-github-dark" choice:"html-modest" choice:"html-tufte" choice:"html-body"`
	Select *bool `long:"select" short:"s"`
}

type DumpFlags struct {}


type ServerFlags struct {
	UrlBase *string `long:"server-url" value-name:"{url}"`
	EndpointIp *string `long:"server-ip" value-name:"{ip}"`
	EndpointPort *uint16 `long:"server-port" value-name:"{port}"`
	EditEnabled *bool `long:"server-edit-enabled"`
	CreateEnabled *bool `long:"server-create-enabled"`
	BrowseEnabled *bool `long:"server-browse-enabled"`
	ClipboardEnabled *bool `long:"server-clipboard-enabled"`
}

type ServerConfiguration struct {
	UrlBase *string `toml:"url_base"`
	EndpointIp *string `toml:"endpoint_ip"`
	EndpointPort *uint16 `toml:"endpoint_port"`
	EditEnabled *bool `toml:"edit_enabled"`
	CreateEnabled *bool `toml:"create_enabled"`
	BrowseEnabled *bool `toml:"browse_enabled"`
	ClipboardEnabled *bool `toml:"clipboard_enabled"`
	OpenExternalConfirm *bool `toml:"open_external_confirm"`
	OpenExternalConfirmSkipForSchemas *[]string `toml:"open_external_confirm_skip_for_schemas"`
	AuthenticationCookieName *string `toml:"authentication_cookie_name"`
	AuthenticationCookieTimeout *uint `toml:"authentication_cookie_timeout"`
	AuthenticationCookieSecret *string `toml:"authentication_cookie_secret"`
}


type BrowseFlags struct {
	UrlBase *string `long:"server-url" value-name:"{url}"`
	Library *string `long:"library" short:"l" value-name:"{identifier}"`
	Document *string `long:"document" short:"d" value-name:"{identifier}"`
	SelectLibrary *bool `long:"select-library" short:"S"`
	SelectDocument *bool `long:"select" short:"s"`
	Authenticate *bool `long:"authenticate" short:"a"`
	UrlDisplay *bool `long:"print"`
	QrcodeDisplay *bool `long:"qrcode"`
}

type BrowserConfiguration struct {
	UrlBase *string `toml:"url_base"`
	AuthenticationSecret *string `toml:"authentication_secret"`
	TerminalOpenInternalCommand *[]string `toml:"terminal_open_internal_command"`
	XorgOpenInternalCommand *[]string `toml:"xorg_open_internal_command"`
	TerminalOpenExternalCommand *[]string `toml:"terminal_open_external_command"`
	XorgOpenExternalCommand *[]string `toml:"xorg_open_external_command"`
}


type HelpFlags struct {}


type MainFlags struct {
	
	Global *GlobalFlags `group:"Global options"`
	Library *LibraryFlags `group:"Library options"`
	Index *IndexFlags `group:"Index options"`
	Debugging *DebuggingFlags `group:"Debugging options"`
	
	List *ListFlags `command:"list"`
	Search *SearchFlags `command:"search"`
	Grep *GrepFlags `command:"grep"`
	
	Create *CreateFlags `command:"create"`
	Edit *EditFlags `command:"edit"`
	Export *ExportFlags `command:"export"`
	Dump *DumpFlags `command:"dump"`
	
	Server *ServerFlags `command:"server"`
	Browse *BrowseFlags `command:"browse"`
	
	Menu *MenuFlags `command:"menu"`
	
	Help *HelpFlags `command:"help"`
}

type MainConfiguration struct {
	
	Global *GlobalConfiguration `toml:"globals"`
	Index *IndexConfiguration `toml:"index"`
	Editor *EditorConfiguration `toml:"editor"`
	Libraries []*Library `toml:"library"`
	
	Server *ServerConfiguration `toml:"server"`
	Browser *BrowserConfiguration `toml:"browser"`
	
	Menus []*Menu `toml:"menu"`
}


type MenuFlags struct {
	Menu *string `long:"menu" short:"m" value-name:"{identifier}"`
	Loop *bool `long:"loop" short:"L"`
}

type Menu struct {
	Identifier string `toml:"identifier"`
	Label string `toml:"label"`
	Commands []*MenuCommand `toml:"commands"`
	Default bool `toml:"default"`
	Loop bool `toml:"loop"`
}

type MenuCommand struct {
	Label string `toml:"label"`
	Command string `toml:"command"`
	Arguments []string `toml:"arguments"`
}




func Main (_executable string, _arguments []string, _environment map[string]string) (*Error) {
	
	_globals, _error := GlobalsNew (_executable, _environment)
	if _error != nil {
		return _error
	}
	defer triggerAtExit (_globals)
	
	_flags := & MainFlags {
			
			Global : & GlobalFlags {},
			Library : & LibraryFlags {},
			Index : & IndexFlags {},
			
			List : & ListFlags {},
			Search : & SearchFlags {},
			Grep : & GrepFlags {},
			
			Create : & CreateFlags {},
			Edit : & EditFlags {},
			Export : & ExportFlags {},
			Dump : & DumpFlags {},
			
			Server : & ServerFlags {},
			Browse : & BrowseFlags {},
			
			Menu : & MenuFlags {},
		}
	
	_configuration := & MainConfiguration {
			Global : & GlobalConfiguration {},
			Index : & IndexConfiguration {},
			Editor : & EditorConfiguration {},
			Server : & ServerConfiguration {},
			Browser : & BrowserConfiguration {},
		}
	
	_parser, _error := mainParserNew (_flags)
	if _error != nil {
		return _error
	}
	
	_help := func (_log bool, _error *Error) (*Error) {
		_buffer := BytesBufferNewSize (4 * 1024)
		defer BytesBufferRelease (_buffer)
		_buffer.WriteByte ('\n')
		_parser.WriteHelp (_buffer)
		_buffer.WriteByte ('\n')
		if _log {
			if _globals.StdioIsTty && _globals.TerminalEnabled {
				logf ('`', 0xa725b4bc, "%s", _buffer.String ())
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
	
	if _flags.Debugging.ProfileCpuPath != nil {
		_close, _error := mainProfileCpu (*_flags.Debugging.ProfileCpuPath)
		if _error != nil {
			return _error
		}
		defer func () () {
				if _error := _close (); _error != nil {
					logError ('w', _error)
				}
			} ()
	}
	if _flags.Debugging.ProfileMemoryPath != nil {
		_close, _error := mainProfileMemory (*_flags.Debugging.ProfileMemoryPath)
		if _error != nil {
			return _error
		}
		defer func () () {
				if _error := _close (); _error != nil {
					logError ('w', _error)
				}
			} ()
	}
	
	if _flags.Global.WorkingDirectory != nil {
		_path := *_flags.Global.WorkingDirectory
		if _path == "" {
			return errorw (0x2289141b, nil)
		}
		if _error := os.Chdir (_path); _error != nil {
			return errorw (0x6fe4c660, _error)
		}
	}
	
	_configurationPath := (*string) (nil)
	_configurationPathInStore := false
	if (_configurationPath == nil) && (_flags.Global.ConfigurationPath != nil) {
		_configurationPath = _flags.Global.ConfigurationPath
	}
	if (_configurationPath == nil) && (len (_flags.Library.Paths) == 0) {
		_homeStore, _ := os.UserHomeDir ()
		_configStore, _ := os.UserConfigDir ()
		for _, _storeAndFolderAndFile := range [][3]string {
				{ ".", "", ".scratchpad" },
				{ ".", "", ".scratchpad.toml" },
				{ ".", "", ".z-scratchpad" },
				{ ".", "", ".z-scratchpad.toml" },
				{ ".", "", "default.toml" },
				{ _homeStore, "", ".scratchpad" },
				{ _homeStore, "", ".scratchpad.toml" },
				{ _homeStore, ".scratchpad", "default.toml" },
				{ _homeStore, "", ".z-scratchpad" },
				{ _homeStore, "", ".z-scratchpad.toml" },
				{ _homeStore, ".z-scratchpad", "default.toml" },
				{ _configStore, "z-scratchpad", "default.toml" },
		} {
			if _storeAndFolderAndFile[0] == "" {
				continue
			}
			_path := path.Join (_storeAndFolderAndFile[0], _storeAndFolderAndFile[1], _storeAndFolderAndFile[2])
			if _stat, _error := os.Stat (_path); _error == nil {
				if _storeAndFolderAndFile[1] == "" {
					if _stat.IsDir () {
						continue
					}
				}
				_configurationPath = &_path
				if _storeAndFolderAndFile[0] == "." {
					_configurationPathInStore = true
				}
				break
			} else if ! os.IsNotExist (_error) {
				return errorw (0xbb4d9103, _error)
			}
		}
	}
	
	if _configurationPath != nil {
		_path := *_configurationPath
		if _path == "" {
			return errorw (0x9a6f64a7, nil)
		}
		_path, _error := filepath.Abs (_path)
		if _error != nil {
			return errorw (0xb9029faf, nil)
		}
		_path, _error = filepath.EvalSymlinks (_path)
		if _error != nil {
			return errorw (0x260be8f0, nil)
		}
		_dataBytes, _error := os.ReadFile (_path)
		if _error != nil {
			return errorw (0xf2be5f5f, _error)
		}
		_dataBuffer := bytes.NewBuffer (_dataBytes)
		defer BytesBufferRelease (_dataBuffer)
		if len (bytes.TrimSpace (_dataBytes)) > 0 {
			_decoder := toml.NewDecoder (_dataBuffer)
			_decoder.Strict (true)
			_error = _decoder.Decode (_configuration)
			if _error != nil {
				return errorw (0x93e9dab8, _error)
			}
		} else {
			if _configurationPathInStore {
				_library, _error := mainLibraryForPaths ([]string { path.Dir (_path) })
				if _error != nil {
					return _error
				}
				_configuration.Libraries = []*Library { _library }
			}
		}
		_globals.ConfigurationPath = _path
	}
	
	if _configuration.Global.UniqueIdentifier != nil {
		_globals.UniqueIdentifier = *_configuration.Global.UniqueIdentifier
	}
	if _globals.UniqueIdentifier == "" {
		if _globals.ConfigurationPath != "" {
			_token := fmt.Sprintf ("%s\000%s", UNAME_NODE, _globals.ConfigurationPath)
			_fingerprint := fingerprintString (_token)
			_globals.UniqueIdentifier = _fingerprint
		}
	}
	
	if _flags.Server.UrlBase != nil {
		_configuration.Server.UrlBase = _flags.Server.UrlBase
	}
	if _configuration.Server.UrlBase == nil {
		_endpointIp := flag2StringOrDefault (_flags.Server.EndpointIp, _configuration.Server.EndpointIp, "127.0.0.1")
		_endpointPort := flag2Uint16OrDefault (_flags.Server.EndpointPort, _configuration.Server.EndpointPort, 49894)
		if _endpointIp_0 := net.ParseIP (_endpointIp); _endpointIp_0 != nil {
			_endpointIp = _endpointIp_0.String ()
		} else {
			return errorw (0x1be6e804, nil)
		}
		_urlBase := ""
		if _endpointPort == 80 {
			_urlBase = fmt.Sprintf ("http://%s/", _endpointIp)
		} else {
			_urlBase = fmt.Sprintf ("http://%s:%d/", _endpointIp, _endpointPort)
		}
		_configuration.Server.UrlBase = &_urlBase
	}
	if _flags.Browse.UrlBase != nil {
		_configuration.Browser.UrlBase = _flags.Browse.UrlBase
	}
	if _configuration.Browser.UrlBase == nil {
		_configuration.Browser.UrlBase = _configuration.Server.UrlBase
	}
	
	if (_configuration.Browser.AuthenticationSecret == nil) && (_configuration.Server.AuthenticationCookieSecret != nil) {
		_configuration.Browser.AuthenticationSecret = _configuration.Server.AuthenticationCookieSecret
	}
	
	_command := ""
	if _parser.Active != nil {
		_command = _parser.Active.Name
	} else {
		if len (_configuration.Menus) > 0 {
			_command = "menu"
		} else {
			_command = "help-abort"
		}
	}
	
	_profilingLoops := uint (0)
	if _flags.Debugging.ProfileLoops != nil {
		_profilingLoops = *_flags.Debugging.ProfileLoops
	}
	if _profilingLoops > 0 {
		for _loop := uint (0); _loop < _profilingLoops; _loop += 1 {
			if _error := MainWithFlags (_command, _flags, _configuration, _globals); _error != nil {
				return _error
			}
			triggerAtExit (_globals)
			_globals.atExitSignal = make (chan struct {})
			debug.SetGCPercent (-1)
		}
		return nil
	}
	
	return MainWithFlags (_command, _flags, _configuration, _globals)
}


func mainParserNew (_flags *MainFlags) (*flags.Parser, *Error) {
	_parser := flags.NewNamedParser ("z-scratchpad", flags.PassDoubleDash)
	_parser.SubcommandsOptional = true
	if _, _error := _parser.AddGroup ("", "", _flags); _error != nil {
		return nil, errorw (0x5b48e356, _error)
	}
	return _parser, nil
}




func MainWithFlags (_command string, _flags *MainFlags, _configuration *MainConfiguration, _globals *Globals) (*Error) {
	
	if _configuration.Global.WorkingDirectory != nil {
		_path := *_configuration.Global.WorkingDirectory
		if _path == "" {
			return errorw (0xe7c58968, nil)
		}
		if _path == "{CONF}" {
			if _globals.ConfigurationPath != "" {
				_path = path.Dir (_globals.ConfigurationPath)
			} else {
				return errorw (0x7b7b780e, nil)
			}
		}
		if _error := os.Chdir (_path); _error != nil {
			return errorw (0x5aae8d30, _error)
		}
	}
	
	_globals.TerminalEnabled = _globals.TerminalEnabled && flagBoolOrDefault (_configuration.Global.TerminalEnabled, true)
	_globals.XorgEnabled = _globals.XorgEnabled && flagBoolOrDefault (_configuration.Global.XorgEnabled, true)
	
	_libraries, _error := mainLibrariesResolve (_flags.Library, _configuration.Libraries)
	if _error != nil {
		return _error
	}
	
	_index, _error := mainIndexNew (_flags.Index, _configuration.Index, _libraries, _globals)
	if _error != nil {
		return _error
	}
	
	_editor, _error := EditorNew (_globals, _index)
	if _error != nil {
		return _error
	}
	
	if _configuration.Editor.DefaultCreateLibrary == nil {
		_libraries, _error := IndexLibrariesSelectAll (_index)
		if _error != nil {
			return _error
		}
		_defaultLibrary := ""
		_multipleLibraries := 0
		for _, _library := range _libraries {
			if _library.CreateEnabled {
				_defaultLibrary = _library.Identifier
				_multipleLibraries += 1
			}
		}
		if _multipleLibraries == 1 {
			_configuration.Editor.DefaultCreateLibrary = &_defaultLibrary
		}
	}
	if _configuration.Editor.DefaultCreateLibrary != nil {
		_library := *_configuration.Editor.DefaultCreateLibrary
		if _library == "" {
			return errorw (0xd3b3131d, nil)
		}
		_editor.DefaultCreateLibrary = _library
	}
	
	if _configuration.Editor.TerminalEditCommand != nil {
		_command := *_configuration.Editor.TerminalEditCommand
		if len (_command) == 0 {
			return errorw (0x28e59c3d, nil)
		}
		_editor.TerminalEditCommand = _command
	}
	if _configuration.Editor.XorgEditCommand != nil {
		_command := *_configuration.Editor.XorgEditCommand
		if len (_command) == 0 {
			return errorw (0x7fd5d86e, nil)
		}
		_editor.XorgEditCommand = _command
	}
	
	if _configuration.Editor.TerminalSelectCommand != nil {
		_command := *_configuration.Editor.TerminalSelectCommand
		if len (_command) == 0 {
			return errorw (0xe9ff3646, nil)
		}
		_editor.TerminalSelectCommand = _command
	}
	if _configuration.Editor.XorgSelectCommand != nil {
		_command := *_configuration.Editor.XorgSelectCommand
		if len (_command) == 0 {
			return errorw (0x8b6b008b, nil)
		}
		_editor.XorgSelectCommand = _command
	}
	
	if _configuration.Editor.TerminalClipboardStoreCommand != nil {
		_command := *_configuration.Editor.TerminalClipboardStoreCommand
		if len (_command) == 0 {
			return errorw (0x957bc35d, nil)
		}
		_editor.TerminalClipboardStoreCommand = _command
	}
	if _configuration.Editor.XorgClipboardStoreCommand != nil {
		_command := *_configuration.Editor.XorgClipboardStoreCommand
		if len (_command) == 0 {
			return errorw (0xd4b2e03d, nil)
		}
		_editor.XorgClipboardStoreCommand = _command
	}
	
	_browser, _error := mainBrowserNew (_configuration.Browser, _globals, _index)
	
	return MainWithFlagsAndContext (_command, _flags, _configuration, _globals, _index, _editor, _browser)
}




func MainWithFlagsAndContext (_command string, _flags *MainFlags, _configuration *MainConfiguration, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	debug.SetGCPercent (200)
	
	switch _command {
		
		
		case "list" :
			return MainList (_flags.List, _globals, _index)
		
		case "search" :
			return MainSearch (_flags.Search, _globals, _index, _editor, _browser)
		
		case "grep" :
			return MainGrep (_flags.Grep, _globals, _index, _editor, _browser)
		
		
		case "create" :
			return MainCreate (_flags.Create, _globals, _index, _editor)
		
		case "edit" :
			return MainEdit (_flags.Edit, _globals, _index, _editor)
		
		case "export" :
			return MainExport (_flags.Export, _globals, _index, _editor)
		
		case "dump" :
			return MainDump (_flags.Dump, _globals, _index)
		
		
		case "server" :
			return MainServer (_flags.Server, _configuration.Server, _globals, _index, _editor, _browser)
		
		case "browse" :
			return MainBrowse (_flags.Browse, _globals, _index, _editor, _browser)
		
		
		case "menu" :
			return MainMenu (_flags.Menu, _configuration.Menus, _configuration, _globals, _index, _editor, _browser)
		
		
		case "help", "help-abort" :
			if _error := MainHelp (_flags.Help, _globals, _editor); _error != nil {
				return _error
			}
			if _command == "help-abort" {
				return errorf (0x7f3fc1a7, "no command was specified, and no default menu was configured!")
			} else {
				return nil
			}
		
		
		default :
			return errorw (0xaca17bb9, nil)
	}
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




func MainSearch (_flags *SearchFlags, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	_libraryIdentifier := flagStringOrDefault (_flags.Library, "")
	_type := flagStringOrDefault (_flags.Type, "document")
	_what := flagStringOrDefault (_flags.What, "identifier")
	_how := flagStringOrDefault (_flags.How, "title")
	_format := flagStringOrDefault (_flags.Format, "text")
	_action := flagStringOrDefault (_flags.Action, "output")
	
	switch _action {
		case "output" :
			// NOP
		case "edit", "export", "browse" :
			if _flags.Type != nil {
				return errorw (0x8133f4ab, nil)
			}
			if _flags.What != nil {
				return errorw (0xf998d0d9, nil)
			}
			if _flags.Format != nil {
				return errorw (0x304ff173, nil)
			}
		default :
			return errorw (0x332d42c3, nil)
	}
	
	_selection, _error := mainListOptionsAndSelect (_libraryIdentifier, _type, _how, _what, _index, _editor)
	if _error != nil {
		return _error
	}
	
	switch _action {
		
		case "output" :
			return mainListOutput (_selection, _format, _globals)
		
		case "edit", "export", "browse" :
			switch len (_selection) {
				case 0 :
					return nil
				case 1 :
					// NOP
				default :
					if ! flagBoolOrDefault (_flags.MultipleAllowed, false) {
						_options := make ([][2]string, 0, len (_selection))
						for _, _selection := range _selection {
							_identifier := _selection[1]
							_document, _error := WorkflowDocumentResolve (_identifier, _index)
							if _error != nil {
								return _error
							}
							_label := _document.Title
							_options = append (_options, [2]string { _label, _identifier })
						}
						_selection, _error = mainListSelect (_options, _editor)
						if _error != nil {
							return _error
						}
					}
			}
			return mainListAction (_selection, _action, _globals, _index, _editor, _browser)
		
		default :
			return errorw (0xe611caea, nil)
	}
}




func MainGrep (_flags *GrepFlags, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	_libraryIdentifier := flagStringOrDefault (_flags.Library, "")
	_what := flagStringOrDefault (_flags.What, "identifier")
	_where := flagStringOrDefault (_flags.Where, "title")
	_format := flagStringOrDefault (_flags.Format, "text")
	_action := flagStringOrDefault (_flags.Action, "output")
	_matchAny := flagBoolOrDefault (_flags.MatchAny, false)
	
	switch _action {
		case "output" :
			// NOP
		case "edit", "export", "browse" :
			if _flags.What != nil {
				return errorw (0x966bbfc4, nil)
			}
			if _flags.Format != nil {
				return errorw (0x92252a21, nil)
			}
		default :
			return errorw (0x4b4f9c3b, nil)
	}
	
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
	
	_matchExpected := len (_terms)
	if _matchAny {
		_matchExpected = 1
	}
	
	_selection := make ([][2]string, 0, len (_options) / 2)
	for _, _option := range _options {
		_contents := _option[0]
		_matchCount := 0
		for _, _term := range _terms {
			if strings.Index (_contents, _term) != -1 {
				_matchCount += 1
				if _matchCount == _matchExpected {
					break
				}
			}
		}
		if _matchCount == _matchExpected {
			_selection = append (_selection, _option)
		}
	}
	
	switch _action {
		
		case "output" :
			return mainListOutput (_selection, _format, _globals)
		
		case "edit", "export", "browse" :
			switch len (_selection) {
				case 0 :
					return nil
				case 1 :
					// NOP
				default :
					if ! flagBoolOrDefault (_flags.MultipleAllowed, false) {
						return errorw (0x1e4d02e6, nil)
					}
			}
			return mainListAction (_selection, _action, _globals, _index, _editor, _browser)
		
		default :
			return errorw (0x1217cd0b, nil)
	}
}




func MainCreate (_flags *CreateFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_identifier := ""
	_exitOnEmpty := false
	_error := (*Error) (nil)
	if _flags.Document != nil {
		if _flags.Select != nil {
			return errorw (0x6c15ba50, nil)
		}
		// FIXME: _identifier, _error = mainResolveDocumentIdentifier (_flags.Library, _flags.Document, _flags.Select, _index, _editor)
		if _flags.Library != nil {
			_identifier = fmt.Sprintf ("%s:%s", *_flags.Library, *_flags.Document)
		} else if _editor.DefaultCreateLibrary != "" {
			_identifier = fmt.Sprintf ("%s:%s", _editor.DefaultCreateLibrary, *_flags.Document)
		} else {
			_identifier = *_flags.Document
		}
	} else if (_flags.Library != nil) || (_flags.Select != nil) {
		_identifier, _error = mainResolveLibraryIdentifier (_flags.Library, _flags.Select, _index, _editor)
		_exitOnEmpty = true
	}
	if _error != nil {
		return _error
	}
	if _exitOnEmpty && (_identifier == "") {
		return nil
	}
	
	return WorkflowDocumentCreate (_identifier, _index, _editor, true)
}




func MainEdit (_flags *EditFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_identifier, _error := mainResolveDocumentIdentifier (_flags.Library, _flags.Document, _flags.Select, _index, _editor)
	if _error != nil {
		return _error
	}
	if _identifier == "" {
		return nil
	}
	
	return WorkflowDocumentEdit (_identifier, _index, _editor, true)
}




func MainExport (_flags *ExportFlags, _globals *Globals, _index *Index, _editor *Editor) (*Error) {
	
	_identifier, _error := mainResolveDocumentIdentifier (_flags.Library, _flags.Document, _flags.Select, _index, _editor)
	if _error != nil {
		return _error
	}
	if _identifier == "" {
		return nil
	}
	
	_format := flagStringOrDefault (_flags.Format, "source")
	
	return mainExportOutput (_identifier, _format, _globals, _index)
}


func mainExportOutput (_identifier string, _format string, _globals *Globals, _index *Index) (*Error) {
	
	_document, _error := WorkflowDocumentResolve (_identifier, _index)
	if _error != nil {
		return _error
	}
	
	_buffer := (*bytes.Buffer) (nil)
	switch _format {
		
		case "identifier" :
			_buffer = bytes.NewBufferString (_document.Identifier)
		
		case "title", "name" :
			_buffer = bytes.NewBufferString (_document.Title)
		
		case "path" :
			_buffer = bytes.NewBufferString (_document.Path)
		
		case "source" :
			if _output, _error := DocumentRenderToSource (_document); _error == nil {
				_buffer = bytes.NewBufferString (_output)
			} else {
				return _error
			}
		
		case "html-body" :
			if _output, _error := DocumentRenderToHtml (_document, true); _error == nil {
				_buffer = bytes.NewBufferString (_output)
			} else {
				return _error
			}
		
		case "html", "html-plain", "html-github", "html-github-auto", "html-github-light", "html-github-dark", "html-modest", "html-tufte" :
			if _format == "html" {
				_format = "html-default"
			}
			_theme := _format[5:]
			_buffer = BytesBufferNewSize (128 * 1024)
			defer BytesBufferRelease (_buffer)
			_templates, _error := TemplatesNew ()
			if _error != nil {
				return _error
			}
			if _error := DocumentRenderToHtmlDocument (_document, true, _theme, _templates, _buffer); _error != nil {
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




func MainDump (_flags *DumpFlags, _globals *Globals, _index *Index) (*Error) {
	
	_documents, _error := IndexDocumentsSelectAll (_index)
	if _error != nil {
		return _error
	}
	
	_buffer := BytesBufferNewSize (128 * 1024)
	defer BytesBufferRelease (_buffer)
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




func mainResolveLibraryIdentifier (_libraryFlag *string, _selectFlag *bool, _index *Index, _editor *Editor) (string, *Error) {
	
	_select := flagBoolOrDefault (_selectFlag, false)
	if _select && (_libraryFlag != nil) {
		return "", errorw (0x4d3444df, nil)
	}
	
	_identifier := ""
	
	if _select {
		
		_options, _error := mainListOptionsAndSelect ("", "library", "title", "identifier", _index, _editor)
		if _error != nil {
			return "", _error
		}
		
		switch len (_options) {
			case 0 :
				return "", nil
			case 1 :
				_identifier = _options[0][1]
			default :
				return "", errorw (0x22d4ddbe, nil)
		}
		
	} else {
		
		if _libraryFlag == nil {
			return "", errorw (0x302d616d, nil)
		}
		
		if _library, _error := LibraryParseIdentifier (*_libraryFlag); _error == nil {
			_identifier = _library
		} else {
			return "", _error
		}
	}
	
	return _identifier, nil
}


func mainResolveDocumentIdentifier (_libraryFlag *string, _documentFlag *string, _selectFlag *bool, _index *Index, _editor *Editor) (string, *Error) {
	
	_select := flagBoolOrDefault (_selectFlag, false)
	if _select && (_documentFlag != nil) {
		return "", errorw (0xaf2210a5, nil)
	}
	
	_identifier := ""
	
	if _select {
		
		_libraryIdentifier := flagStringOrDefault (_libraryFlag, "")
		_options, _error := mainListOptionsAndSelect (_libraryIdentifier, "document", "title", "identifier", _index, _editor)
		if _error != nil {
			return "", _error
		}
		
		switch len (_options) {
			case 0 :
				return "", nil
			case 1 :
				_identifier = _options[0][1]
			default :
				return "", errorw (0x43982abc, nil)
		}
		
	} else {
		
		if _library, _document, _error := mainMergeLibraryAndDocumentIdentifiers (_libraryFlag, _documentFlag); _error != nil {
			return "", _error
		} else if _document != "" {
			_identifier = _document
		} else if _library != "" {
			return "", errorw (0x20457235, nil)
		} else {
			return "", errorw (0x0684c89e, nil)
		}
	}
	
	return _identifier, nil
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
					case "commonmark-link" :
						_valueEscaped := _library.Name
						_valueEscaped = strings.ReplaceAll (_valueEscaped, "\\", "\\\\")
						_valueEscaped = strings.ReplaceAll (_valueEscaped, "[", "\\[")
						_valueEscaped = strings.ReplaceAll (_valueEscaped, "]", "\\]")
						_value = fmt.Sprintf ("[%s](sl:%s)", _valueEscaped, _library.Identifier)
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
					case "commonmark-link" :
						_valueEscaped := _document.TitleOriginal
						_valueEscaped = strings.ReplaceAll (_valueEscaped, "\\", "\\\\")
						_valueEscaped = strings.ReplaceAll (_valueEscaped, "[", "\\[")
						_valueEscaped = strings.ReplaceAll (_valueEscaped, "]", "\\]")
						_value = fmt.Sprintf ("[%s](sd:%s)", _valueEscaped, _document.Identifier)
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
	
	_values := make (map[string]map[string]bool, len (_options))
	_valuesDuplicate := false
	for _, _option := range _options {
		_label := _option[0]
		_value := _option[1]
		_label = stringTrimSpaces (_label)
		_values_1 := map[string]bool (nil)
		if _values_0, _exists := _values[_label]; _exists {
			_values_1 = _values_0
			_valuesDuplicate = true
		} else {
			_values_1 = make (map[string]bool, 16)
			_values[_label] = _values_1
		}
		_values_1[_value] = true
	}
	
	_labels := make ([]string, 0, len (_values))
	_labelsMap := make (map[string]string, len (_values))
	for _label, _values := range _values {
		if !_valuesDuplicate {
			_labels = append (_labels, _label)
			_labelsMap[_label] = _label
		} else {
			_labelWithCount := ""
			if len (_values) > 1 {
				_labelWithCount = fmt.Sprintf ("%3d | %s", len (_values), _label)
			} else {
				_labelWithCount = fmt.Sprintf ("      %s", _label)
			}
			_labels = append (_labels, _labelWithCount)
			_labelsMap[_labelWithCount] = _label
		}
	}
	
	sortfold.Strings (_labels)
	
	_selection_0, _error := EditorSelect (_editor, _labels)
	if _error != nil {
		return nil, _error
	}
	
	_selection := make ([][2]string, 0, 16)
	for _, _label := range _selection_0 {
		if _label_0, _exists := _labelsMap[_label]; _exists {
			_label = _label_0
		} else {
			return nil, errorw (0xa37f357b, nil)
		}
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
	
	sortfold.Strings (_list)
	
	_buffer := BytesBufferNewSize (128 * 1024)
	defer BytesBufferRelease (_buffer)
	
	switch _format {
		
		case "text", "text-sp", "text-0" :
			_separator := byte ('\n')
			_always := true
			if _format == "text-0" {
				_separator = 0
			} else if _format == "text-sp" {
				_separator = byte (' ')
				_always = false
			}
			_first := true
			for _, _value := range _list {
				if !_first {
					_buffer.WriteByte (_separator)
				}
				_buffer.WriteString (_value)
				_first = false
			}
			if _always {
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


func mainListAction (_options [][2]string, _action string, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
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
	
	for _, _identifier := range _list {
		_error := (*Error) (nil)
		switch _action {
			case "edit" :
				_error = WorkflowDocumentEdit (_identifier, _index, _editor, true)
			case "browse" :
				_error = WorkflowDocumentBrowse (_identifier, _index, _browser, true)
			case "export" :
				// FIXME:  Add support for other formats!
				_error = mainExportOutput (_identifier, "source", _globals, _index)
			default :
				return errorw (0xb5fa0b59, nil)
		}
		if _error != nil {
			return _error
		}
	}
	
	return nil
}




func MainServer (_flags *ServerFlags, _configuration *ServerConfiguration, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	_endpointIp := flag2StringOrDefault (_flags.EndpointIp, _configuration.EndpointIp, "127.0.0.1")
	_endpointPort := flag2Uint16OrDefault (_flags.EndpointPort, _configuration.EndpointPort, 49894)
	if _endpointIp_0 := net.ParseIP (_endpointIp); _endpointIp_0 != nil {
		_endpointIp = _endpointIp_0.String ()
	} else {
		return errorw (0xfb27d8b0, nil)
	}
	
	_editEnabled := flag2BoolOrDefault (_flags.EditEnabled, _configuration.EditEnabled, false)
	_createEnabled := flag2BoolOrDefault (_flags.CreateEnabled, _configuration.CreateEnabled, false)
	_browseEnabled := flag2BoolOrDefault (_flags.BrowseEnabled, _configuration.BrowseEnabled, false)
	_clipboardEnabled := flag2BoolOrDefault (_flags.ClipboardEnabled, _configuration.ClipboardEnabled, false)
	_openExternalConfirm := flagBoolOrDefault (_configuration.OpenExternalConfirm, true)
	_openExternalConfirmSkipForSchemas := flagStringsOrDefault (_configuration.OpenExternalConfirmSkipForSchemas, nil)
	
	_endpoint := fmt.Sprintf ("%s:%d", _endpointIp, _endpointPort)
	
	_listener, _error_0 := net.Listen ("tcp", _endpoint)
	if _error_0 != nil {
		return errorw (0xedeea766, _error_0)
	}
	
	_globals.TerminalEnabled = false
	
	_server, _error := ServerNew (_globals, _index, _editor, _browser, _listener)
	if _error != nil {
		return _error
	}
	
	_server.EditEnabled = _server.EditEnabled && _editEnabled
	_server.CreateEnabled = _server.CreateEnabled && _createEnabled
	_server.BrowseEnabled = _server.BrowseEnabled && _browseEnabled
	_server.ClipboardEnabled = _server.ClipboardEnabled && _clipboardEnabled
	_server.OpenExternalConfirm = _server.OpenExternalConfirm || _openExternalConfirm
	_server.OpenExternalConfirmSkipForSchemas = append (_server.OpenExternalConfirmSkipForSchemas, _openExternalConfirmSkipForSchemas ...)
	
	if _configuration.UrlBase != nil {
		_server.UrlBase = *_configuration.UrlBase
	} else {
		return errorw (0xf9bef449, nil)
	}
	
	if _configuration.AuthenticationCookieName != nil {
		_server.AuthenticationCookieName = *_configuration.AuthenticationCookieName
	}
	if _configuration.AuthenticationCookieTimeout != nil {
		_server.AuthenticationCookieTimeout = *_configuration.AuthenticationCookieTimeout
	}
	if _configuration.AuthenticationCookieSecret != nil {
		_server.AuthenticationCookieSecret = *_configuration.AuthenticationCookieSecret
	}
	
	logf ('i', 0x210494be, "[server]  access URL `%s`;  listening on `%s`;", _server.UrlBase, _endpoint)
	
	_error = ServerRun (_server)
	if _error != nil {
		return _error
	}
	
	return nil
}




func MainBrowse (_flags *BrowseFlags, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	if (_flags.SelectLibrary != nil) && (_flags.SelectDocument != nil) {
		return errorw (0x8dbd7a13, nil)
	}
	if (_flags.SelectLibrary != nil) && (_flags.Document != nil) {
		return errorw (0xdb7b52f6, nil)
	}
	
	_libraryIdentifier := ""
	_documentIdentifier := ""
	_exitOnEmpty := false
	_error := (*Error) (nil)
	if (_flags.Document != nil) || (_flags.SelectDocument != nil) {
		_documentIdentifier, _error = mainResolveDocumentIdentifier (_flags.Library, _flags.Document, _flags.SelectDocument, _index, _editor)
		_exitOnEmpty = true
	} else if (_flags.Library != nil) || (_flags.SelectLibrary != nil) {
		_libraryIdentifier, _error = mainResolveLibraryIdentifier (_flags.Library, _flags.SelectLibrary, _index, _editor)
		_exitOnEmpty = true
	}
	if _error != nil {
		return _error
	}
	if _exitOnEmpty && (_libraryIdentifier == "") && (_documentIdentifier == "") {
		return nil
	}
	
	_authenticate := flagBoolOrDefault (_flags.Authenticate, false)
	if _authenticate {
		if _browser.ServerAuthenticationSecret == "" {
			return errorw (0x8602b60d, nil)
		}
	} else {
		_browser.ServerAuthenticationSecret = ""
	}
	
	_urlDisplay := flagBoolOrDefault (_flags.UrlDisplay, false)
	_browser.UrlDisplay = _urlDisplay
	
	_qrcodeDisplay := flagBoolOrDefault (_flags.QrcodeDisplay, false)
	_browser.QrcodeDisplay = _qrcodeDisplay
	
	if _documentIdentifier != "" {
		return WorkflowDocumentBrowse (_documentIdentifier, _index, _browser, true)
	}
	if _libraryIdentifier != "" {
		return WorkflowLibraryBrowse (_libraryIdentifier, _index, _browser, true)
	}
	if _authenticate {
		return BrowserHomeOpen (_browser, true)
	} else {
		return BrowserIndexOpen (_browser, true)
	}
}


func mainBrowserNew (_configuration *BrowserConfiguration, _globals *Globals, _index *Index) (*Browser, *Error) {
	
	_browser, _error := BrowserNew (_globals, _index)
	if _error != nil {
		return nil, _error
	}
	
	if _configuration.TerminalOpenInternalCommand != nil {
		_command := *_configuration.TerminalOpenInternalCommand
		if len (_command) == 0 {
			return nil, errorw (0xd23959ac, nil)
		}
		_browser.TerminalOpenInternalCommand = _command
	}
	if _configuration.XorgOpenInternalCommand != nil {
		_command := *_configuration.XorgOpenInternalCommand
		if len (_command) == 0 {
			return nil, errorw (0x045b13e4, nil)
		}
		_browser.XorgOpenInternalCommand = _command
	}
	
	if _configuration.TerminalOpenExternalCommand != nil {
		_command := *_configuration.TerminalOpenExternalCommand
		if len (_command) == 0 {
			return nil, errorw (0xeaf6799d, nil)
		}
		_browser.TerminalOpenExternalCommand = _command
	}
	if _configuration.XorgOpenExternalCommand != nil {
		_command := *_configuration.XorgOpenExternalCommand
		if len (_command) == 0 {
			return nil, errorw (0x7bf7b1d7, nil)
		}
		_browser.XorgOpenExternalCommand = _command
	}
	
	if _configuration.UrlBase != nil {
		_browser.ServerUrlBase = *_configuration.UrlBase
	} else {
		return nil, errorw (0xa88827e6, nil)
	}
	if _configuration.AuthenticationSecret != nil {
		_browser.ServerAuthenticationSecret = *_configuration.AuthenticationSecret
	}
	
	return _browser, nil
}




func MainMenu (_flags *MenuFlags, _menus []*Menu, _configuration *MainConfiguration, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	_menuIdentifier := flagStringOrDefault (_flags.Menu, "")
	if _menuIdentifier == "" {
		for _, _menu := range _menus {
			// NOTE:  We select the first default menu...
			if (_menu.Identifier != "") && _menu.Default {
				_menuIdentifier = _menu.Identifier
				break
			}
		}
	}
	if _menuIdentifier == "" {
		return errorw (0x876f0980, nil)
	}
	
	_menu := (*Menu) (nil)
	for _, _menu_0 := range _menus {
		if _menu_0.Identifier == _menuIdentifier {
			_menu = _menu_0
			break
		}
	}
	
	_options := make ([][2]string, 0, len (_menu.Commands))
	_commands := make (map[string]*MenuCommand, len (_menu.Commands))
	for _, _command := range _menu.Commands {
		if _command.Label == "" {
			return errorw (0x854ba0ab, nil)
		}
		if _command.Command == "" {
			return errorw (0xdd4d0687, nil)
		}
		if _, _exists := _commands[_command.Label]; _exists {
			return errorw (0x6c32847a, nil)
		}
		_options = append (_options, [2]string { _command.Label, _command.Label })
		_commands[_command.Label] = _command
	}
	
	_loop := flagBoolOrDefault (_flags.Loop, _menu.Loop)
	
	for {
		
		_selection, _error := mainListSelect (_options, _editor)
		if _error != nil {
			return _error
		}
		
		_selected := ""
		switch len (_selection) {
			case 0 :
				return nil
			case 1 :
				_selected = _selection[0][1]
			default :
				return errorw (0xde0c52f4, nil)
		}
		
		_command, _ := _commands[_selected]
		if _command == nil {
			return errorw (0x2f57b12e, nil)
		}
		
		if _error := MainCommand (_command.Command, _command.Arguments, _configuration, _globals, _index, _editor, _browser); _error != nil {
			return _error
		}
		
		if !_loop {
			break
		}
	}
	
	return nil
}



func MainCommand (_command string, _arguments []string, _configuration *MainConfiguration, _globals *Globals, _index *Index, _editor *Editor, _browser *Browser) (*Error) {
	
	_flags_0 := interface{} (nil)
	_execute := (func () (*Error)) (nil)
	
	switch _command {
		case "edit" :
			_flags := & EditFlags {}
			_flags_0 = _flags
			_execute = func () (*Error) {
					return MainEdit (_flags, _globals, _index, _editor)
				}
		case "create" :
			_flags := & CreateFlags {}
			_flags_0 = _flags
			_execute = func () (*Error) {
					return MainCreate (_flags, _globals, _index, _editor)
				}
		case "search" :
			_flags := & SearchFlags {}
			_flags_0 = _flags
			_execute = func () (*Error) {
					return MainSearch (_flags, _globals, _index, _editor, _browser)
				}
		case "browse" :
			_flags := & BrowseFlags {}
			_flags_0 = _flags
			_execute = func () (*Error) {
					return MainBrowse (_flags, _globals, _index, _editor, _browser)
				}
		case "menu" :
			_flags := & MenuFlags {}
			_flags_0 = _flags
			_execute = func () (*Error) {
					return MainMenu (_flags, _configuration.Menus, _configuration, _globals, _index, _editor, _browser)
				}
		default :
			return errorw (0xbd997a82, nil)
	}
	
	_parser := flags.NewNamedParser (_command, flags.PassDoubleDash)
	if _, _error := _parser.AddGroup ("", "", _flags_0); _error != nil {
		return errorw (0x8d45cee0, _error)
	}
	
	if _argumentsRest, _error := _parser.ParseArgs (_arguments); _error != nil {
		return errorw (0x0ddaf31b, _error)
	} else if len (_argumentsRest) != 0 {
		return errorw (0xdc656ded, nil)
	}
	
	return _execute ()
}




func MainHelp (_flags *HelpFlags, _globals *Globals, _editor *Editor) (*Error) {
	_parser, _error := mainParserNew (& MainFlags {})
	if _error != nil {
		return _error
	}
	_buffer := BytesBufferNewSize (128 * 1024)
	defer BytesBufferRelease (_buffer)
	_buffer.WriteByte ('\n')
	_parser.WriteHelp (_buffer)
	_buffer.WriteByte ('\n')
	if _, _error := _buffer.WriteTo (_globals.Stdout); _error != nil {
		return errorw (0x7626a566, _error)
	}
	return nil
}




func mainIndexNew (_flags *IndexFlags, _configuration *IndexConfiguration, _libraries []*Library, _globals *Globals) (*Index, *Error) {
	
	_index, _error := IndexNew (_globals)
	if _error != nil {
		return nil, _error
	}
	
	_databaseCanWalk := ! flagBoolOrDefault (_flags.WalkDisabled, false)
	_databaseEnabled := ! flagBoolOrDefault (_flags.DatabaseDisabled, false) && flagBoolOrDefault (_configuration.DatabaseEnabled, true)
	_databaseCanLoad := ! flagBoolOrDefault (_flags.LoadDisabled, false)
	_databaseCanStore := ! flagBoolOrDefault (_flags.StoreDisabled, false)
	_databaseCanDirty := ! flagBoolOrDefault (_flags.DirtyDisabled, false)
	_databaseCanRefresh := ! flagBoolOrDefault (_flags.RefreshDisabled, false)
	
	_databasePath := ""
	_databaseDirtyPath := ""
	if _databaseEnabled {
		if _configuration.DatabasePath != nil {
			_databasePath = *_configuration.DatabasePath
		}
		if (_databasePath == "") && (_globals.UniqueIdentifier != "") {
			_databasePath = path.Join ("{CACHEDIR}", _globals.UniqueIdentifier + ".db")
		}
		if (_databasePath != "") && strings.HasPrefix (_databasePath, "{CACHEDIR}") {
			_cachePath, _error := os.UserCacheDir ()
			if _error != nil {
				return nil, errorw (0xf1aa16da, _error)
			}
			_cachePath = path.Join (_cachePath, "z-scratchpad")
			_error = os.MkdirAll (_cachePath, 0o750)
			if _error != nil {
				return nil, errorw (0xc65eca13, _error)
			}
			_databasePath = _databasePath[10:]
			_databasePath = path.Join (_cachePath, _databasePath)
		}
		if (_databasePath != "") && strings.HasPrefix (_databasePath, "{TMPDIR}") {
			_databasePath = _databasePath[8:]
			_databasePath = path.Join (_globals.TemporaryDirectory, _databasePath)
		}
	}
	if _databasePath != "" {
		_databaseDirtyPath = _databasePath + "-dirty"
	}
	
	if _databasePath == "" {
		_databaseEnabled = false
	}
	if _databaseDirtyPath == "" {
		_databaseCanDirty = false
	}
	
	_databaseCanLoad = _databaseCanLoad && _databaseEnabled
	_databaseCanStore = _databaseCanStore && _databaseEnabled
	_databaseCanDirty = _databaseCanDirty && _databaseEnabled
	_databaseCanRefresh = _databaseCanRefresh && _databaseEnabled
	
	if _error := mainIndexLoad (_index, _libraries, _databasePath, _databaseDirtyPath, _databaseCanWalk, _databaseCanLoad, _databaseCanStore, _databaseCanDirty); _error != nil {
		return nil, _error
	}
	
	_librariesRefresh := _databaseEnabled && _databaseCanWalk && _databaseCanRefresh && _index.librariesRefreshEnabled
	_librariesRefresh = _librariesRefresh && flagBoolOrDefault (_configuration.LibrariesRefreshEnabled, true)
	_index.librariesRefreshEnabled = _librariesRefresh
	if _librariesRefresh {
		_index.librariesRefreshCallback = func (_index_0 *Index) (*Error) {
					if _index_0 != _index {
						return errorw (0x7bfcc92c, nil)
					}
					if !_databaseCanLoad {
						return nil
					}
					return mainIndexLoad (_index, _libraries, _databasePath, _databaseDirtyPath, _databaseCanWalk, _databaseCanLoad, _databaseCanStore, _databaseCanDirty)
				}
	}
	
	_documentsRefresh := _databaseEnabled && _databaseCanWalk && _databaseCanRefresh && _index.documentRefreshEnabled
	_documentsRefresh = _documentsRefresh && flagBoolOrDefault (_configuration.DocumentsRefreshEnabled, true)
	_index.documentRefreshEnabled = _documentsRefresh
	if _documentsRefresh {
		_index.documentRefreshCallback = func (_index_0 *Index, _document *Document) (*Document, *Error) {
					if _index_0 != _index {
						return nil, errorw (0x435c1cc3, nil)
					}
					return WorkflowDocumentRefresh (_document, _index)
				}
	}
	
	_databaseCanDirty = _databaseCanDirty && _index.dirtyEnabled
	_index.dirtyEnabled = _databaseCanDirty
	if _databaseCanDirty {
		_markDirty := func () (*Error) {
//				logf ('d', 0xfcc0490e, "marking dirty database...")
				_timestamp := time.Now ()
				if _, _error := os.Stat (_databaseDirtyPath); _error == nil {
					if _error := os.Chtimes (_databaseDirtyPath, _timestamp, _timestamp); _error != nil {
						return errorw (0x392bb23b, _error)
					}
				} else if os.IsNotExist (_error) {
					if _file, _error := os.OpenFile (_databaseDirtyPath, os.O_WRONLY | os.O_CREATE, 0640); _error == nil {
						_file.Close ()
					} else {
						return errorw (0xb0f0682c, _error)
					}
					if _error := os.Chtimes (_databaseDirtyPath, _timestamp, _timestamp); _error != nil {
						return errorw (0xfe939aea, _error)
					}
				} else {
					return errorw (0xe29e9be2, _error)
				}
				_index.refreshTimestamp = _timestamp
				return nil
			}
		_index.dirtyCallback = func (_index_0 *Index) (*Error) {
				if _index_0 != _index {
					return errorw (0x086e68ac, nil)
				}
				return _markDirty ()
			}
	}
	
	return _index, nil
}


func mainIndexLoad (_index *Index, _libraries []*Library, _databasePath string, _databaseDirtyPath string, _databaseCanWalk bool, _databaseCanLoad bool, _databaseCanStore bool, _databaseCanDirty bool) (*Error) {
	
	_beginTimestamp := time.Now ()
	
	_databaseShouldWalk := _databaseCanWalk
	_databaseShouldLoad := _databaseCanLoad
	_databaseShouldStore := _databaseCanStore
	_databaseTimestamp := time.Time {}
	
	if _databaseShouldLoad {
		if _stat, _error := os.Stat (_databasePath); _error == nil {
			_databaseTimestamp = _stat.ModTime ()
		} else if ! os.IsNotExist (_error) {
			return errorw (0xc35078ab, _error)
		} else {
			_databaseShouldLoad = false
		}
	}
	if _databaseShouldLoad && _databaseCanDirty {
		if _stat, _error := os.Stat (_databaseDirtyPath); _error == nil {
			if _stat.ModTime () .After (_databaseTimestamp) {
//				logf ('d', 0xd23ee348, "index database not loaded (dirty);")
				_databaseShouldLoad = false
				_databaseTimestamp = _stat.ModTime ()
			}
		} else if ! os.IsNotExist (_error) {
			return errorw (0xfb9afa94, _error)
		}
	}
	
	if _databaseShouldLoad {
		if ! _databaseTimestamp.After (_index.refreshTimestamp) {
//			logf ('d', 0xe1b51623, "index database not loaded (unchanged);")
			return nil
		}
	}
	
	_databaseLoaded := false
	if _databaseShouldLoad {
		if _loaded, _error := IndexLoadFromPath (_index, _databasePath); _error != nil {
			return _error
		} else if _loaded {
//			logf ('d', 0xae0fba86, "index database loaded;")
			_databaseLoaded = true
			_databaseShouldWalk = false
			_databaseShouldStore = false
		} else {
//			logf ('i', 0x1dff7d9a, "index database not loaded due to incompatible versions;")
		}
	}
	
	if _databaseShouldWalk {
		_databaseTimestamp = time.Now ()
		if _error := mainIndexWalkAndLoad (_index, _libraries); _error != nil {
			return _error
		}
//		logf ('d', 0xb6d41aa3, "index database walked;")
		_databaseLoaded = true
	}
	
	if !_databaseLoaded {
		return errorw (0xacfdaefc, nil)
	}
	
	_index.refreshTimestamp = _databaseTimestamp
	
	if _databaseShouldStore {
		if _error := IndexStoreToPath (_index, _databasePath); _error != nil {
			return _error
		}
//		logf ('d', 0x3c58b633, "index database stored;")
	}
	
	_elapsed := time.Since (_beginTimestamp)
	_elapsedMilliseconds := _elapsed.Milliseconds ()
	if (_databaseShouldWalk && (_elapsedMilliseconds >= 200)) || (!_databaseShouldWalk && (_elapsedMilliseconds > 75)) {
		logf ('d', 0x67aaf8be, "index loading took %d milliseconds;", _elapsedMilliseconds)
	}
	
	return nil
}


func mainIndexWalkAndLoad (_index *Index, _libraries []*Library) (*Error) {
	
	_documentPaths, _error := mainLibrariesWalk (_libraries)
	if _error != nil {
		return _error
	}
	
	_documents, _error := mainLibrariesLoad (_libraries, _documentPaths)
	if _error != nil {
		return _error
	}
	
	_error = mainLibrariesInclude (_index, _libraries, _documents)
	if _error != nil {
		return _error
	}
	
	return nil
}


func mainLibrariesResolve (_flags *LibraryFlags, _configuration []*Library) ([]*Library, *Error) {
	
	if (len (_flags.Paths) > 0) && (len (_configuration) > 0) {
		return nil, errorw (0x374ece0f, nil)
	}
	
	_libraries := make ([]*Library, 0, 16)
	
	if len (_flags.Paths) > 0 {
		_library, _error := mainLibraryForPaths (_flags.Paths)
		if _error != nil {
			return nil, _error
		}
		_libraries = append (_libraries, _library)
	}
	
	if len (_configuration) > 0 {
		for _, _library_0 := range _configuration {
			_library := & Library {}
			*_library = *_library_0
			_libraries = append (_libraries, _library)
		}
	}
	
	if len (_libraries) == 0 {
		return nil, errorw (0x00ea182b, nil)
	}
	
	for _index := 0; _index < len (_libraries); _index += 1 {
		if ! _libraries[_index].Disabled {
			continue
		}
		if _index < len (_libraries) - 1 {
			copy (_libraries[_index :], _libraries[_index + 1 :])
			_index -= 1
		}
		_libraries[len (_libraries) - 1] = nil
		_libraries = _libraries[: len (_libraries) - 1]
	}
	
	for _, _library := range _libraries {
		if _error := LibraryInitialize (_library); _error != nil {
			return nil, _error
		}
	}
	
	return _libraries, nil
}


func mainLibrariesWalk (_libraries []*Library) ([][][2]string, *Error) {
	
	_documents := make ([][][2]string, 0, len (_libraries))
	
	for _, _library := range _libraries {
		
		_libraryDocuments, _error := libraryDocumentsWalk (_library)
		if _error != nil {
			return nil, _error
		}
		
		_documents = append (_documents, _libraryDocuments)
	}
	
	return _documents, nil
}


func mainLibrariesLoad (_libraries []*Library, _libraryDocuments [][][2]string) ([]*Document, *Error) {
	
	_documents := make ([]*Document, 0, 16 * 1024)
	
	for _libraryIndex := range _libraries {
		
		_library := _libraries[_libraryIndex]
		_libraryDocumentPaths := _libraryDocuments[_libraryIndex]
		
		_libraryDocuments, _error := libraryDocumentsLoad (_library, _libraryDocumentPaths)
		if _error != nil {
			return nil, _error
		}
		
		for _, _document := range _libraryDocuments {
			
			if _document.Library == "" {
				_document.Library = _library.Identifier
			}
			
			_document.EditEnabled = _library.EditEnabled
			
			_error = DocumentInitializeIdentifier (_document, _library)
			if _error != nil {
				return nil, _error
			}
			
			_error = DocumentInitializeFormat (_document, _library)
			if _error != nil {
				return nil, _error
			}
			
			_error = DocumentInitializeTitle (_document, _library)
			if _error != nil {
				return nil, _error
			}
			
			_documents = append (_documents, _document)
		}
	}
	
	return _documents, nil
}


func mainLibrariesInclude (_index *Index, _libraries []*Library, _documents []*Document) (*Error) {
	
	IndexClearData (_index)
	
	_dirtyEnabled := _index.dirtyEnabled
	_index.dirtyEnabled = false
	
	for _, _library := range _libraries {
		_error := IndexLibraryInclude (_index, _library)
		if _error != nil {
			return _error
		}
	}
	for _, _document := range _documents {
		
		_error := IndexDocumentInclude (_index, _document)
		if _error != nil {
			return _error
		}
	}
	
	_index.dirtyEnabled = _dirtyEnabled
	
	return nil
}


func mainLibraryForPaths (_paths []string) (*Library, *Error) {
	_library := & Library {
			Identifier : "library",
			Name : "Library",
			Paths : _paths,
			UsePathInLibraryAsIdentifier : true,
			UseFileExtensionAsFormat : true,
			IncludeGlobPatterns : []string { "**/*.{md,markdown,gmi,gemini,txt,text}" },
			EditEnabled : true,
			CreateEnabled : true,
			CreatePath : _paths[0],
		}
	return _library, nil
}




func mainProfileCpu (_path string) (func () (*Error), *Error) {
	if _path == "" {
		return nil, errorw (0x6a07cc92, nil)
	}
	logf ('i', 0x0b2ca2b0, "profiling CPU, writing output to `%s`...", _path)
	_output, _error := os.OpenFile (_path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0o640)
	if _error != nil {
		return nil, errorw (0x8704258a, _error)
	}
	_error = pprof.StartCPUProfile (_output)
	if _error != nil {
		return nil, errorw (0x8cceeb40, _error)
	}
	_close := func () (*Error) {
			defer _output.Close ()
			pprof.StopCPUProfile ()
			return nil
		}
	return _close, nil
}

func mainProfileMemory (_path string) (func () (*Error), *Error) {
	if _path == "" {
		return nil, errorw (0xfc9c43d8, nil)
	}
	logf ('i', 0x16f55857, "profiling memory, writing output to `%s`...", _path)
	_output, _error := os.OpenFile (_path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0o640)
	if _error != nil {
		return nil, errorw (0xd1f70512, _error)
	}
	_close := func () (*Error) {
			defer _output.Close ()
			runtime.GC ()
			_error := pprof.WriteHeapProfile (_output)
			if _error != nil {
				return errorw (0x3e5724d2, nil)
			}
			return nil
		}
	return _close, nil
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

func flagStringsOrDefault (_value *[]string, _default []string) ([]string) {
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

func flag2StringsOrDefault (_value_1 *[]string, _value_2 *[]string, _default []string) ([]string) {
	if _value_1 != nil {
		return *_value_1
	}
	if _value_2 != nil {
		return *_value_2
	}
	return _default
}

