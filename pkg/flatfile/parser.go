package flatfile

import (
	"fmt"
	"strconv"
	"strings"
)

type flatfileTag struct {
	start int
	end   string
}

var parseFuncMap = map[string]func(string, *flatfileTag) error{
	"start": parseStart,
	"end":   parseEnd,
}

//validOptions is a list of the keys loaded from parseFuncMap. This is used purely to display options to user
var validOptions []string

func init() {
	for k := range parseFuncMap {
		validOptions = append(validOptions, k)
	}
}

func parseFlatfileTag(tag string, parserTag *flatfileTag) error {
	var err error
	//split tag into key and value
	params := strings.Split(tag, ",")
	if len(params) < 1 {
		return fmt.Errorf("flatfile.parseFlatfileTag: Not enough flatfile tag parameters provided.\nStart parameter is required.\nMust be in the form of `flatfile:start[,end]`")
	}
	for i, param := range params {
		//check for named options
		if strings.Contains(param, "=") {
			options := strings.Split(param, "=")
			if len(options) != 2 {
				return fmt.Errorf("flatfile.parseFlatfileTag: Invalid formatting of named option, '%v'.\nMust be in the form of `key=value`", param)
			}
			if funcVal, exists := parseFuncMap[options[0]]; exists {
				err = funcVal(options[1], parserTag)
			} else {
				return fmt.Errorf("flatfile.parseFlatfileTag: Invalid tag parameter '%v'.\nMust be one of %v", options[0], validOptions)
			}
		} else {
			//check for positional options
			switch i {
			case 0:
				err = parseStart(param, parserTag)
			case 1:
				err = parseEnd(param, parserTag)
			}
		}
		if err != nil {
			return fmt.Errorf("flatfile.parseFlatfileTag: Failed to parse tag parameter '%v'.\n%w", param, err)
		}
	}
	if parserTag.start == 0 {
		return fmt.Errorf("flatfile.parseFlatfileTag: Start parameter is required.\nMust be in the form of `flatfile:start[,end]`")
	}
	if parserTag.end == "" {
		parserTag.end = strconv.Itoa(parserTag.start)
	}
	return nil
}

func parseStart(param string, parserTag *flatfileTag) error {
	start, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("flatfile.parseStart: Failed to parse start parameter '%v'.\n%w", param, err)
	}
	if start < 1 {
		return fmt.Errorf("flatfile.parseStart: Start parameter must be greater than 0. Fields are indexed starting at 1.\n%v", param)
	}

	parserTag.start = start
	return nil
}

func parseEnd(param string, parserTag *flatfileTag) error {
	_, err := strconv.Atoi(param)
	if err == nil {
		parserTag.end = param
		return nil
	}
	if param != "" {
		return fmt.Errorf("flatfile.parseEnd: Failed to parse end parameter '%v'. Must be a field index or '*'.\n%w", param, err)
	}
	parserTag.end = "*"
	return nil
}
