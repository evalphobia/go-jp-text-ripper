package plugin

import (
	"strconv"

	"github.com/evalphobia/go-jp-text-ripper/ripper"
)

// NounNameCountPlugin calculates personal name word count from tokenized words
var NounNameCountPlugin = &ripper.Plugin{
	Title: "noun_name_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("人名"))
	},
}

// NounNumberCountPlugin calculates number word count from tokenized words
var NounNumberCountPlugin = &ripper.Plugin{
	Title: "noun_number_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("数"))
	},
}

// NounLocationCountPlugin calculates location word count from tokenized words
var NounLocationCountPlugin = &ripper.Plugin{
	Title: "noun_location_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("地域"))
	},
}

// NounOrganizationCountPlugin calculates organization word count from tokenized words
var NounOrganizationCountPlugin = &ripper.Plugin{
	Title: "noun_organization_count",
	Fn: func(text *ripper.TextData) string {
		return strconv.Itoa(text.GetWords().CountFeatures("組織"))
	},
}

// NounHasFullNamePlugin calculates personal full name from tokenized words
var NounHasFullNamePlugin = &ripper.Plugin{
	Title: "noun_has_fullname",
	Fn: func(text *ripper.TextData) string {
		w := text.GetWords()
		switch {
		case !w.HasFeatures("姓"):
			return "false"
		case !w.HasFeatures("名"):
			return "false"
		}
		return "true"
	},
}
