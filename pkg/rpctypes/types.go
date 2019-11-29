package rpctypes

type OutPoint struct {
	TxHash Hash   `json:"tx_hash"`
	Index  Uint32 `json:"index"`
}

type CellInput struct {
	Since          Uint64   `json:"since"`
	PreviousOutput OutPoint `json:"previous_output"`
}

type Script struct {
	CodeHash Hash           `json:"code_hash"`
	HashType ScriptHashType `json:"hash_type"`
	Args     Bytes          `json:"args"`
}

type CellOutput struct {
	Capacity Uint64  `json:"capacity"`
	Lock     Script  `json:"lock"`
	Type     *Script `json:"type,omitempty"`
}

type CellDep struct {
	OutPoint OutPoint `json:"out_point"`
	DepType  DepType  `json:"dep_type"`
}

type Transaction struct {
	Version     Uint32       `json:"version"`
	CellDeps    []CellDep    `json:"cell_deps"`
	HeaderDeps  []Hash       `json:"header_deps"`
	Inputs      []CellInput  `json:"inputs"`
	Outputs     []CellOutput `json:"outputs"`
	Witnesses   []Bytes      `json:"witnesses"`
	OutputsData []Bytes      `json:"outputs_data"`
}

type Header struct {
	Version          Uint32 `json:"version"`
	CompactTarget    Uint32 `json:"compact_target"`
	ParentHash       Hash   `json:"parent_hash"`
	Timestamp        Uint64 `json:"timestamp"`
	Number           Uint64 `json:"number"`
	Epoch            Uint64 `json:"epoch"`
	TransactionsRoot Hash   `json:"transactions_root"`
	ProposalsHash    Hash   `json:"proposals_hash"`
	UnclesHash       Hash   `json:"uncles_hash"`
	Dao              Bytes  `json:"dao"`
	// TODO: deal with Uint128 later
	Nonce Bytes `json:"nonce"`
}

type UncleBlock struct {
	Header    Header            `json:"header"`
	Proposals []ProposalShortId `json:"proposals"`
}

type Block struct {
	Header       Header            `json:"header"`
	Uncles       []UncleBlock      `json:"uncles"`
	Transactions []Transaction     `json:"transactions"`
	Proposals    []ProposalShortId `json:"proposals"`
}
