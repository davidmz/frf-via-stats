package main

type srcInfo struct {
	viaJSON
	Count int `json:"count"`
}

type statInfo struct {
	UserName string     `json:"username"`
	Sources  []*srcInfo `json:"sources"`
}

type userJSON struct {
	UserName string `json:"id"`
	Type     string `json:"type"`
}

type viaJSON struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type entryJSON struct {
	From userJSON `json:"from"`
	Via  viaJSON  `json:"via"` // always not nil but may have zero value
}
