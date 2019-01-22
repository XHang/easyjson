package gen

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type tagHandler func(tag, preTag, nextTag string) error

func defaultTagHandler(tag, preTag, nextTag string) error {
	return nil
}

var (
	// defaultTagHandlers enumerates all the default tag handler
	defaultTagHandlers = map[string]tagHandler{
		"<-":       defaultTagHandler,
		"->":       defaultTagHandler,
		"PK":       defaultTagHandler,
		"NULL":     defaultTagHandler,
		"NOT":      defaultTagHandler,
		"AUTOINCR": defaultTagHandler,
		"DEFAULT":  defaultTagHandler,
		"CREATED":  defaultTagHandler,
		"UPDATED":  defaultTagHandler,
		"DELETED":  defaultTagHandler,
		"VERSION":  defaultTagHandler,
		"UTC":      defaultTagHandler,
		"LOCAL":    defaultTagHandler,
		"NOTNULL":  defaultTagHandler,
		"INDEX":    defaultTagHandler,
		"UNIQUE":   defaultTagHandler,
		"CACHE":    defaultTagHandler,
		"NOCACHE":  defaultTagHandler,
		"COMMENT":  defaultTagHandler,
	}
)

func parseXormFieldName(f reflect.StructField) (ret string, err error) {

	ormTagStr := f.Tag.Get("xorm")

	if ormTagStr == "" {
		return
	}

	tags := splitTag(ormTagStr)

	if len(tags) == 0 {
		return
	}

	if tags[0] == "-" || strings.ToUpper(tags[0]) == "EXTENDS" {
		ret = "-"
		return
	}

	var tagName, preTag, nextTag string
	// var params []string

	for j, key := range tags {

		k := strings.ToUpper(key)
		tagName = k
		// params = []string{}

		pStart := strings.Index(k, "(")
		if pStart == 0 {
			err = errors.New("( could not be the first charactor")
			return
		}
		if pStart > -1 {
			if !strings.HasSuffix(k, ")") {
				err = fmt.Errorf("field %s tag %s cannot match ) charactor", f.Name, key)
				return
			}

			tagName = k[:pStart]
			// params = strings.Split(key[pStart+1:len(k)-1], ",")
		}

		if j > 0 {
			preTag = strings.ToUpper(tags[j-1])
		}
		if j < len(tags)-1 {
			nextTag = tags[j+1]
		} else {
			nextTag = ""
		}

		if h, ok := defaultTagHandlers[tagName]; ok {
			if _err := h(tagName, preTag, nextTag); _err != nil {
				err = _err
				return
			}
		} else {
			if strings.HasPrefix(key, "'") && strings.HasSuffix(key, "'") {
				tagName = key[1 : len(key)-1]
			} else {
				tagName = key
			}
		}
	}

	ret = tagName
	return
}

func splitTag(tag string) (tags []string) {
	tag = strings.TrimSpace(tag)
	var hasQuote = false
	var lastIdx = 0
	for i, t := range tag {
		if t == '\'' {
			hasQuote = !hasQuote
		} else if t == ' ' {
			if lastIdx < i && !hasQuote {
				tags = append(tags, strings.TrimSpace(tag[lastIdx:i]))
				lastIdx = i + 1
			}
		}
	}
	if lastIdx < len(tag) {
		tags = append(tags, strings.TrimSpace(tag[lastIdx:]))
	}
	return
}
