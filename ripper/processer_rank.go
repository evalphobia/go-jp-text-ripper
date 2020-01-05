package ripper

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

// DoRank creates *RankProcessor from config and run it.
func DoRank(conf RankConfig) error {
	if err := conf.Init(); err != nil {
		return err
	}
	if err := conf.Validate(); err != nil {
		return err
	}

	conf.Logger.Infof("DoRank", "version:[%s] rev:[%s]", conf.Version, conf.Revision)
	r, err := NewRankProcessor(conf)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := r.WriteHeader(); err != nil {
		return err
	}

	return r.DoWithProgress()
}

// RankProcessor is struct for putting spaces between words
type RankProcessor struct {
	*CommonProcessor
	Config RankConfig
}

// NewRankProcessor returns initialized RankProcessor.
func NewRankProcessor(c RankConfig) (*RankProcessor, error) {
	common, err := NewCommonProcessor(c.CommonConfig)
	if err != nil {
		return nil, err
	}

	r := &RankProcessor{
		CommonProcessor: common,
		Config:          c,
	}
	return r, nil
}

// ReadHeader reads header columns and sets target column by index.
func (r *RankProcessor) ReadHeader() error {
	c := r.Config
	switch {
	case c.ColumnNumber > 0:
		return r.CommonProcessor.ReadHeaderWithIndex(c.ColumnNumber - 1)
	default:
		return r.readHeaderByName(c.Column)
	}
}

// readHeaderByName reads header columns and check target column is existed or not.
func (r *RankProcessor) readHeaderByName(col string) error {
	err := r.CommonProcessor.ReadHeader()
	if err != nil {
		return err
	}

	hasColumn := false
	header := r.inputHeader
	for idx, val := range header {
		if val == col {
			r.columnIndex = idx
			hasColumn = true
		}
	}
	if !hasColumn {
		return fmt.Errorf("cannnot find column name in header: col:[%s] headers:[%+v]", col, header)
	}
	return nil
}

// WriteHeader writes header columns
func (r *RankProcessor) WriteHeader() error {
	// read header if not read yet
	if len(r.inputHeader) == 0 {
		err := r.ReadHeader()
		if err != nil {
			return err
		}
	}

	r.outputHeader = []string{
		"type",
		"rank",
		"word",
		"countN",
		"countP",
	}

	// write to file
	return r.w.Write(r.outputHeader)
}

// DoWithProgress processes with showing progress.
func (r *RankProcessor) DoWithProgress() error {
	r.ShowProgress()

	conf := r.Config
	logger := conf.Logger
	logger.Infof("DoWithProgress", "read lines...")

	err := r.Do()
	if err != nil {
		logger.Errorf("DoWithProgress", "error on r.Do() err:[%s]", err.Error())
		return err
	}

	logger.Infof("DoWithProgress", "finish process")
	return nil
}

// Do processes word frequency ranking.
func (r *RankProcessor) Do() error {
	defer r.Close()
	c := r.Config
	logger := c.Logger

	rank, err := r.GetRank()
	if err != nil {
		return err
	}

	logger.Infof("Do", "Total Words:%d", rank.GetTotalWordSize())

	if err := r.output("top", rank.TopList); err != nil {
		return err
	}
	return r.output("last", rank.LastList)
}

// GetRank gets result of word freqency ranking.
func (r *RankProcessor) GetRank() (RankResult, error) {
	c := r.Config

	rank, err := r.getRank()
	if err != nil {
		return rank, err
	}

	rank.TopList, err = r.FilterByRank(rank, c.TopNumber, c.TopPercent, func(i int) int {
		return i
	})
	if err != nil {
		return rank, err
	}

	totalwords := rank.GetTotalWordSize()
	rank.LastList, err = r.FilterByRank(rank, c.LastNumber, c.LastPercent, func(i int) int {
		return totalwords - i - 1
	})
	return rank, err
}

func (r *RankProcessor) getRank() (RankResult, error) {
	defer r.r.Close()
	c := r.Config
	logger := c.Logger
	idx := r.columnIndex
	tok := r.tok

	lastLineNo := 1
	lastLineText := ""
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		logger.Errorf("GetRank", "unknown error occurred on Line:[%d] Text:[%s]\n", lastLineNo, lastLineText)
	}()

	result := RankResult{}
	resultMap := make(map[string]int, 1024)
	isUnique := c.UseUnique
	totalCount := 0
	for {
		lastLineNo++
		line, err := r.r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Errorf("GetRank", "r.r.Read() err:[%s]\n", err.Error())
			return result, err
		}

		text := &TextData{}

		// tokenize text
		lastLineText = line[idx]
		text.raw = line[idx]
		text.normalized = r.applyPreFilters(text.raw)
		text.words, text.nonWords = tok.Tokenize(text.normalized)
		if err != nil {
			return result, err
		}

		// create result line
		words := text.words.GetWords()
		wordMap := make(map[string]int, len(words))
		for _, w := range words {
			wordMap[w]++
		}
		for word, count := range wordMap {
			if isUnique {
				resultMap[word]++
				totalCount++
				continue
			}
			resultMap[word] += count
			totalCount += count
		}
	}
	result.List = createWordCountList(resultMap)
	result.TotalCount = totalCount
	return result, nil
}

// FilterByRank filters word freqency ranking to use only top rank (or last).
func (r *RankProcessor) FilterByRank(rank RankResult, maxN int, maxP float64, getIndex func(i int) int) (wordCountList, error) {
	totalwords := rank.GetTotalWordSize()
	totalcount := rank.TotalCount
	sumP := 0.0

	results := make(wordCountList, 0, totalwords)
	// filter by rank of top/last
	for i := 0; ; i++ {
		if i >= totalwords {
			break
		}
		idx := getIndex(i)
		v := rank.List[idx]
		rankN := idx + 1
		countP := float64(v.count) / float64(totalcount)

		if rankN > maxN && sumP > maxP {
			break
		}
		sumP += countP

		v.percent = countP
		rank.List[idx] = v
		results = append(results, v)
	}
	return results, nil
}

func (r *RankProcessor) output(typ string, list wordCountList) error {
	c := r.Config
	logger := c.Logger

	for i, v := range list {
		rankN := i + 1

		if c.ShowResult {
			logger.Infof("output", "[%s] #%d %s:%d (%.05f)", typ, rankN, v.word, v.count, v.percent)
		}

		results := []string{
			typ,
			strconv.Itoa(rankN),
			v.word,
			strconv.Itoa(v.count),
			strconv.FormatFloat(v.percent, 'f', 5, 64),
		}
		err := r.w.Write(results)
		if err != nil {
			logger.Errorf("output", "r.w.Write() err:[%s]\n", err.Error())
			return err
		}
	}
	return nil
}

// RankResult has a word frequency ranking result.
type RankResult struct {
	TotalCount int
	List       wordCountList
	TopList    wordCountList
	LastList   wordCountList
}

// GetTotalWordSize returns word types count.
func (r RankResult) GetTotalWordSize() int {
	return len(r.List)
}

// GetTopWords returns high frequency words.
func (r RankResult) GetTopWords() []string {
	list := make([]string, len(r.TopList))
	for i, v := range r.TopList {
		list[i] = v.word
	}
	return list
}

// GetLastWords returns low frequency words.
func (r RankResult) GetLastWords() []string {
	list := make([]string, len(r.LastList))
	for i, v := range r.LastList {
		list[i] = v.word
	}
	return list
}

// for sorting word rank
type wordCount struct {
	word    string
	count   int
	percent float64
}

type wordCountList []wordCount

func createWordCountList(data map[string]int) wordCountList {
	list := make(wordCountList, len(data))
	i := 0
	for k, v := range data {
		list[i] = wordCount{
			word:  k,
			count: v,
		}
		i++
	}
	sort.Sort(sort.Reverse(list))
	return list
}

func (l wordCountList) Len() int           { return len(l) }
func (l wordCountList) Less(i, j int) bool { return l[i].count < l[j].count }
func (l wordCountList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
