package indexer

import "strings"

const notAssigned = 999999

// TorrentCriteria ...
type TorrentCriteria struct {
	MinSeeders      int
	RequiredKeyword []string
}

// BrowsingBase ...
type BrowsingBase struct {
	Category string
	Path     string
	Criteria TorrentCriteria
}

// TorrentWebSite ...
type TorrentWebSite struct {
	URL       string
	Browsings map[string]BrowsingBase
}

// ExtractNameAndTag ...
func ExtractNameAndTag(str string, tags []string) (s string, ret []string) {
	s = strings.ToLower(str)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, ".", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")

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
