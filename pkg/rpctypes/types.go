package rpctypes

type OutPoint struct {
	TxHash Hash   `json:"tx_hash"`
	Index  Uint32 `json:"index"`

	GraphqlCell     *CellOutput   `json:"cell,omitempty"`
	GraphqlCellData *GraphqlBytes `json:"cell_data,omitempty"`
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

	GraphqlResolvedOutPoints []OutPoint `json:"resolved_out_points,omitempty"`
}

type RawTransaction struct {
	Version     Uint32       `json:"version"`
	CellDeps    []CellDep    `json:"cell_deps"`
	HeaderDeps  []Hash       `json:"header_deps"`
	Inputs      []CellInput  `json:"inputs"`
	Outputs     []CellOutput `json:"outputs"`
	OutputsData []Bytes      `json:"outputs_data"`

	GraphqlResolvedHeaderDeps []GraphqlHeaderDep `json:"resolved_header_deps,omitempty"`
}

type Transaction struct {
	RawTransaction
	Witnesses []Bytes `json:"witnesses"`
}

type RawHeader struct {
	Version          Uint32 `json:"version"`
	CompactTarget    Uint32 `json:"compact_target"`
	Timestamp        Uint64 `json:"timestamp"`
	Number           Uint64 `json:"number"`
	Epoch            Uint64 `json:"epoch"`
	ParentHash       Hash   `json:"parent_hash"`
	TransactionsRoot Hash   `json:"transactions_root"`
	ProposalsHash    Hash   `json:"proposals_hash"`
	UnclesHash       Hash   `json:"uncles_hash"`
	Dao              Raw    `json:"dao"`
}

type Header struct {
	RawHeader
	Nonce Uint128 `json:"nonce"`
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

type CellbaseWitness struct {
	Lock    Script `json:"lock"`
	Message Bytes  `json:"message"`
}

type WitnessArgs struct {
	Lock       *Bytes `json:"lock,omitempty"`
	InputType  *Bytes `json:"input_type,omitempty"`
	OutputType *Bytes `json:"output_type,omitempty"`
}
