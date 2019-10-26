package indexer

import "strings"

const notAssigned = 999999

// TorrentCriteria ...
type TorrentCriteria struct {
	MinSeeders int      `json:"min_seeders"`
	Tags       []string `json:"tags"`
}

// BrowsingBase ...
type BrowsingBase struct {
	Category string          `json:"category"`
	Path     string          `json:"path"`
	Criteria TorrentCriteria `json:"criteria"`
}

// TorrentWebSite ...
type TorrentWebSite struct {
	URL       string                  `json:"url"`
	Browsings map[string]BrowsingBase `json:"browsings"`
}

// Settings ...
type Settings struct {
	Timeout    TimeoutCrawling           `json:"timeout"`
	DB         DBSettings                `json:"database"`
	Extraction ExtractionSettings        `json:"extraction_name_entity"`
	Crawler    map[string]TorrentWebSite `json:"crawlers"`
}

// Replace ...
type Replace struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ExtractionSettings ...
type ExtractionSettings struct {
	Replacements   []Replace `json:"replacements"`
	TryExtractYear bool      `json:"TryExtractYear"`
}

var exSettings ExtractionSettings

// SetSettings ...
func SetSettings(s Settings) {
	timeout = s.Timeout
	dbsettings = s.DB
	exSettings = s.Extraction

	for n, c := range crawlers {
		if s, ok := s.Crawler[n]; ok {
			c.SetSettings(s)
		} else {
			panic("No setting for crawler" + n)
		}
	}
}

// ExtractNameAndTag ...
func ExtractNameAndTag(str string, tags []string) (s string, ret []string) {
	s = strings.ToLower(str)
	for _, r := range exSettings.Replacements {
		s = strings.ReplaceAll(s, r.From, r.To)
	}

	firstIndexTag := notAssigned
	for _, t := range tags {
		t = strings.ToLower(t)
		index := strings.Index(s, t)
		if index != -1 {
			ret = append(ret, t)
			if index < firstIndexTag {
				firstIndexTag = index
			}
		}

	}

	if firstIndexTag > 0 && firstIndexTag != notAssigned {
		s = s[0:firstIndexTag]
	}
	return s, ret
}
