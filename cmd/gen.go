package cmd

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func Gen(config Config, files []string) error {
	for _, file := range files {
		err := generateDdl(config, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateDdl(config Config, file string) (resErr error) {
	log.Infof("==> File:%s", file)

	f, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			resErr = err
		}
	}()

	var ignoreRegExps []*regexp.Regexp
	for _, pattern := range config.IgnoreSheet {
		ignoreRegExps = append(ignoreRegExps, regexp.MustCompile(pattern))
	}

	sheets := f.GetSheetList()
SHEET_LOOP:
	for _, sheet := range sheets {
		log.Infof("==> Sheet: %s", sheet)
		for _, reqExp := range ignoreRegExps {
			if reqExp.MatchString(sheet) {
				log.Info("====> Skip")
				continue SHEET_LOOP
			}
		}
		table, err := NewTable(config, f, sheet)
		if err != nil {
			return err
		}
		fmt.Println(table.ToDDL())
	}
	return nil
}
