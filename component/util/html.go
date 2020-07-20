package util // import "github.com/teamgrit-lab/cojam/component/util"

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/teamgrit-lab/cojam/mvc/domain"
)

// MakeTagCSS ...
func MakeTagCSS(tags []string) template.HTML {
	t := time.Now()
	ts := t.Format("20060102150405")
	var result string
	for _, value := range tags {
		result += `<link rel="stylesheet" type="text/css" href="/static/` + value + `?t=` + ts + `" />`
	}
	return template.HTML(result)
}

// MakeTagExternalCSS ...
func MakeTagExternalCSS(tags []string) template.HTML {
	var result string
	for _, value := range tags {
		result += `<link rel="stylesheet" type="text/css" href="` + value + `" />`
	}
	return template.HTML(result)
}

// MakeTagJavascript function gets the internal JavaScript library.
func MakeTagJavascript(tags []string) template.HTML {
	t := time.Now()
	ts := t.Format("20060102150405")
	var result string
	for _, value := range tags {
		result += `<script type="text/javascript" src="/static/` + value + `?t=` + ts + `"></script>`
	}
	return template.HTML(result)
}

// MakeTagExternalJavascript function imports an external JavaScript library.
func MakeTagExternalJavascript(tags []string) template.HTML {
	var result string
	for _, value := range tags {
		result += `<script type="text/javascript" src="` + value + `"></script>`
	}
	return template.HTML(result)
}

// MakeTagHashTag ...
// Buyer Live HashTag
func MakeTagHashTag(tag string) template.HTML {
	var result string
	tags := strings.Split(tag, "#")
	for _, v := range tags {
		if v == "" {
			continue
		}
		result += `<span># ` + strings.TrimSpace(v) + `</span>&nbsp`
	}
	return template.HTML(result)
}

// MakeSelectOption ...
func MakeSelectOption(codes []*domain.CommonCode, status string) template.HTML {
	result := "<option value=''>선택</option>"
	var selected string
	for _, v := range codes {
		if v.Code == status {
			selected = "selected"
		}
		result += fmt.Sprintf("<option value='%s' %s>%s</option>", v.Code, selected, v.CodeName)
		selected = ""
	}
	return template.HTML(result)
}
