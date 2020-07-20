package domain

type EsBody struct {
	Took int        `json:"took"`
	Hits EsMainHits `json:"hits"`
}

type EsTotal struct {
	Value int `json:"value"`
}

type EsSource struct {
	IndexedAt string `json:"indexed_at"`
	Date      string `json:"date"`
	From      string `json:"from"`
	Text      string `json:"text"`
	Room      uint64 `json:"room"`
}

type EsHitItem struct {
	Index  string   `json:"_index"`
	Source EsSource `json:"_source"`
}

type EsMainHits struct {
	Total EsTotal     `json:"total"`
	Hits  []EsHitItem `json:"hits"`
}
