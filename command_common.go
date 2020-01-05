package main

// CommonOption of sub commands.
type CommonOption struct {
	Column           string `cli:"c,column" usage:"target column name in input file"`
	ColumnNumber     int    `cli:"columnn" usage:"target column index in input file (1st col=1)"`
	Input            string `cli:"*i,input" usage:"input file path --input='/path/to/input.csv'"`
	Output           string `cli:"o,output" usage:"output file path --output='./my_result.csv'"`
	Dictionary       string `cli:"dic" usage:"custom dictionary path (mecab ipa dictionaly)"`
	StopWord         string `cli:"stopword" usage:"stop word list file path"`
	ShowResult       bool   `cli:"show" usage:"print separated words to console"`
	UseOriginalForm  bool   `cli:"original" usage:"output original form of word"`
	UseNoun          bool   `cli:"noun" usage:"output 'noun' type of word"`
	UseVerb          bool   `cli:"verb" usage:"output 'verb' type of word"`
	UseAdjective     bool   `cli:"adjective" usage:"output 'adjective' type of word"`
	UseNeologd       bool   `cli:"neologd" usage:"use prefilter for neologd"`
	ProgressInterval int    `cli:"progress" usage:"print current progress (sec)" dft:"30"`
	MinLetterSize    int    `cli:"min" usage:"minimum letter size for output" dft:"1"`
}
