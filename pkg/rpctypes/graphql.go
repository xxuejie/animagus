package rpctypes

type GraphqlBytes struct {
	Hash    *Hash   `json:"hash,omitempty"`
	Length  *Uint32 `json:"length,omitempty"`
	Content *Raw    `json:"content,omitempty"`
}

type GraphqlHeaderDep struct {
	Hash   *Hash   `json:"hash,omitempty"`
	Header *Header `json:"header,omitempty"`
}
