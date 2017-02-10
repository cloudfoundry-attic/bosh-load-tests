package action

type Table struct {
	Rows [][]string `json:"Rows"`
}

type Output struct {
	Tables []Table `json:"Tables"`
}