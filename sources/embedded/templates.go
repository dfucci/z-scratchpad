

package embedded


import _ "embed"




//go:embed templates/libraries-index.html
var LibrariesIndexHtml string

//go:embed templates/libraries-index.txt
var LibrariesIndexText string


//go:embed templates/documents-index.html
var DocumentsIndexHtml string

//go:embed templates/documents-index.txt
var DocumentsIndexText string


//go:embed templates/library-view.html
var LibraryViewHtml string

//go:embed templates/library-view.txt
var LibraryViewText string


//go:embed templates/document-view.html
var DocumentViewHtml string

//go:embed templates/document-view.txt
var DocumentViewText string


//go:embed templates/document-export.html
var DocumentExportHtml string

//go:embed templates/document-export.txt
var DocumentExportText string

//go:embed templates/document-source.txt
var DocumentExportSource string


//go:embed templates/global-navigation.html
var GlobalNavigationHtml string

//go:embed templates/global-partials.html
var GlobalPartialsHtml string


//go:embed templates/version.html
var VersionHtml string

