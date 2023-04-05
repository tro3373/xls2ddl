package cmd

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

const defaultTableOption = "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci"

type Table struct {
	TableOption      string
	Table            string
	TableName        string
	AdditionalOption string
	Cols             []*Col
}

type Col struct {
	Column      string
	ColumnName  string
	Type        string
	Digit       string
	Decimal     string
	NN          string
	PK          string
	Description string
}

func NewTable(config Config, f *excelize.File, sheet string) (*Table, error) {
	cells := []string{
		config.TableCell,
		config.TableNameCell,
		config.AdditionalOptionCell,
	}
	vals, err := readCellValues(f, sheet, cells)
	if err != nil {
		return nil, err
	}

	i := config.ColumnStartRow
	var cols []*Col
	for {
		col, err := newCol(config, f, sheet, i)
		if err != nil {
			return nil, err
		}
		if col == nil {
			break
		}
		cols = append(cols, col)
		i++
	}

	tableOption := defaultTableOption
	if len(config.TableOption) != 0 {
		tableOption = config.TableOption
	}

	table := Table{
		TableOption:      tableOption,
		Table:            vals[0],
		TableName:        vals[1],
		AdditionalOption: vals[2],
		Cols:             cols,
	}
	return &table, nil
}

func newCol(config Config, f *excelize.File, sheet string, row int) (*Col, error) {
	cells := []string{
		colCell(config.ColumnCol, row),
		colCell(config.ColumnNameCol, row),
		colCell(config.TypeCol, row),
		colCell(config.DigitCol, row),
		colCell(config.DecimalCol, row),
		colCell(config.NNCol, row),
		colCell(config.PKCol, row),
		colCell(config.DescriptionCol, row),
	}
	vals, err := readCellValues(f, sheet, cells)
	if err != nil {
		return nil, err
	}
	if vals[0] == "" {
		return nil, nil
	}

	col := Col{
		Column:      vals[0],
		ColumnName:  vals[1],
		Type:        vals[2],
		Digit:       vals[3],
		Decimal:     vals[4],
		NN:          vals[5],
		PK:          vals[6],
		Description: vals[7],
	}
	return &col, nil
}

func colCell(col string, row int) string {
	return fmt.Sprintf("%s%d", col, row)
}

func readCellValues(f *excelize.File, sheet string, cells []string) ([]string, error) {
	var res []string
	for _, cell := range cells {
		val, err := f.GetCellValue(sheet, cell)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func (col *Col) ToDDL() string {
	digitDecimal := col.Digit
	colType := parseColType(col.Type)
	if len(digitDecimal) != 0 {
		if isDecimalableType(colType) && len(col.Decimal) != 0 {
			digitDecimal += "," + col.Decimal
		}
	}
	if len(digitDecimal) != 0 {
		digitDecimal = fmt.Sprintf("(%s)", digitDecimal)
	}
	if strings.EqualFold(col.Column, "id") {
		colType = "INT UNSIGNED AUTO_INCREMENT"
		digitDecimal = ""
	}
	if isNoDigitDecimalType(colType) {
		digitDecimal = ""
	}

	nn := ""
	if len(col.NN) != 0 {
		nn = "NOT NULL"
	}
	return fmt.Sprintf("`%s` %s%s %s COMMENT '%s'", col.Column, colType, digitDecimal, nn, col.ColumnName)
}

func parseColType(colType string) string {
	if colType == "DEC" {
		return "DECIMAL"
	}
	return colType
}
func isDecimalableType(colType string) bool {
	switch colType {
	case "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC":
		return true
	}
	return false
}
func isNoDigitDecimalType(colType string) bool {
	switch colType {
	case "DATE", "TIME", "DATETIME", "TIMESTAMP", "YEAR":
		return true
	}
	return false
}

func (table *Table) ToDDL() string {
	var lines []string
	var cols []string
	var pks []string
	lines = append(lines, fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", table.Table))
	lines = append(lines, fmt.Sprintf("CREATE TABLE `%s` (", table.Table))
	for _, col := range table.Cols {
		cols = append(cols, col.ToDDL())
		if len(col.PK) != 0 {
			pks = append(pks, col.Column)
		}
	}
	lines = append(lines, " "+strings.Join(cols, "\n,"))
	if len(pks) != 0 {
		lines = append(lines, fmt.Sprintf(", PRIMARY KEY (`%s`)", strings.Join(pks, "`,`")))
	}
	if len(table.AdditionalOption) != 0 {
		lines = append(lines, ", "+table.AdditionalOption)
	}
	lines = append(lines, fmt.Sprintf(") %s", table.TableOption))
	lines = append(lines, fmt.Sprintf("comment='%s';", table.TableName))
	lines = append(lines, "") // Add row after ddl end
	return strings.Join(lines, "\n")
}
