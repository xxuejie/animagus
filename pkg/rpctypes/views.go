package rpctypes

type TransactionView struct {
	Transaction
	Hash Hash `json:"hash"`
}

type TransactionWithStatusView struct {
	Transaction TransactionView `json:"transaction"`
	TxStatus    TxStatus        `json:"tx_status"`
}

type HeaderView struct {
	Header
	Hash Hash `json:"hash"`
}

type UncleBlockView struct {
	Header    HeaderView        `json:"header"`
	Proposals []ProposalShortId `json:"proposals"`
}

type BlockView struct {
	Header       HeaderView        `json:"header"`
	Uncles       []UncleBlockView  `json:"uncles"`
	Transactions []TransactionView `json:"transactions"`
	Proposals    []ProposalShortId `json:"proposals"`
}
