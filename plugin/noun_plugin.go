package plugin

import (
	"strconv"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// NameCountPlugin calculates personal name word count from tokenized words
var NameCountPlugin = &ripper.Plugin{
	Title: "name_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("人名"))
	},
}

// NumberCountPlugin calculates number word count from tokenized words
var NumberCountPlugin = &ripper.Plugin{
	Title: "number_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("数"))
	},
}

// LocationCountPlugin calculates location word count from tokenized words
var LocationCountPlugin = &ripper.Plugin{
	Title: "location_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("地域"))
	},
}

// OrganizationCountPlugin calculates organization word count from tokenized words
var OrganizationCountPlugin = &ripper.Plugin{
	Title: "organization_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("組織"))
	},
}
