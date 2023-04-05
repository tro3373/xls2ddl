package cmd

type Config struct {
	TableOption          string
	IgnoreSheet          []string
	TableCell            string
	TableNameCell        string
	AdditionalOptionCell string
	ColumnStartRow       int
	ColumnCol            string
	ColumnNameCol        string
	TypeCol              string
	DigitCol             string
	DecimalCol           string
	NNCol                string
	PKCol                string
	DescriptionCol       string
}
