package parsers

import (
	"regexp"

	"github.com/codeshelldev/goredacted/internals/parser"
)

var nameParser = parser.Parser{
	Name: "name",
	Parse: func(this parser.Parser, content string) string {
		result := content

		for _, pattern := range CustomParserPatterns(this.Name) {
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}

			result = re.ReplaceAllStringFunc(result, func(match string) string {
				return SelectReplacement(match, FIRST_NAMES)
			})
		}

		return result
	},
}

func init() {
	parser.Register(nameParser)
}

var FIRST_NAMES = []string{
	"alice","brian","charlotte","daniel","emily","georgia","henry","isabel",
	"jackson","katherine","leonard","madeline","nathaniel","olivia","patrick",
	"quinn","rebecca","sebastian","theresa","ulysses","wyatt","xenia","zeke",
	"amelia","bradley","cassandra","derek","elizabeth","frederick","gabriella",
	"harold","iris","jonathan","kelly","lawrence","melissa","nicholas","ophelia",
	"philip","quincy","rosalie","stanley","tabitha","ursula","wilhelmina","xander",
	"yvana","zachary","adrian","bianca","colin","danica","edward","fiona","gregory",
	"helen","ian","jocelyn","kael","lila","marcus","nina","orlando","paige","quinton",
	"rachel","samuel","tessa","ulyana","vincent","wendy","xavier","yasmine","zane",
	"amelie","bruno","celeste","dominic","eva","felix","gretchen","harvey","imogen",
	"jasper","karina","lucas","monica","nathan","octavia","phillip","quincy","renata",
	"silas","theron","ursula","valeria","wyatt","xandra","yasmin","zander",
}

var LAST_NAMES = []string{
	"morgan","foster","hill","carter","russell","bell","adams","james","wood",
	"clark","morris","king","lee","turner","harris","mitchell","murphy","price",
	"evans","collins","bryant","simmons","perry","cook","ward","ross","morgan",
	"powell","taylor","morris","wood","ellis","bennett","hart","robinson","jenkins",
	"brooks","murray","wright","hamilton","harrison","richardson","griffin","stevens",
	"sanders","fleming","may","cole","knight","miles","parker","fleming","holt","smith",
	"dunn","fairchild","morgan","price","hunter","dawson","foster","lawrence","bennett",
	"hamilton","james","coleman","garcia","roberts","carter","hill","brown","hansen",
	"hart","morris","ford","clark","jenkins","evans","mitchell","robinson","lee",
	"sanders","price","taylor","hughes","turner","stewart","wright","morgan","brooks",
	"ross","simmons","dunn","cole","fleming","hamilton","jenkins","morris","holt",
	"foster",
}
