package template

import (
	"bytes"
	"text/template"
	"time"

	"github.com/go-sprout/sprout"
	"github.com/go-sprout/sprout/registry/std"
	"github.com/go-sprout/sprout/registry/strings"
	sproutTime "github.com/go-sprout/sprout/registry/time"
	"github.com/lestrrat-go/strftime"
	extFmt "github.com/ngyewch/sprout-ext/fmt"
	extTime "github.com/ngyewch/sprout-ext/time"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Helper struct {
	timeLocation    *time.Location
	dateTimeFormat  *strftime.Strftime
	templateFuncMap template.FuncMap
	templateCache   map[string]*template.Template
}

func NewHelper(timeLocation *time.Location, dateTimeFormat string) (*Helper, error) {
	sproutHandler := sprout.New()
	err := sproutHandler.AddRegistries(
		extTime.NewRegistry(),
		extFmt.NewRegistry(),
		sproutTime.NewRegistry(),
		strings.NewRegistry(),
		std.NewRegistry(),
	)
	if err != nil {
		return nil, err
	}

	dateTimeFormatStrftime, err := strftime.New(dateTimeFormat)
	if err != nil {
		return nil, err
	}

	templateFuncMap := sproutHandler.Build()

	templateHelper := &Helper{
		timeLocation:    timeLocation,
		dateTimeFormat:  dateTimeFormatStrftime,
		templateFuncMap: templateFuncMap,
		templateCache:   make(map[string]*template.Template),
	}

	templateFuncMap["standardDate"] = templateHelper.StandardDate

	return templateHelper, nil
}

func (helper *Helper) TimeLocation() *time.Location {
	return helper.timeLocation
}

func (helper *Helper) StandardDate(date any) string {
	return helper.dateTimeFormat.FormatString(helper.toTime(date))
}

func (helper *Helper) CreateTemplate(templateString string) (*template.Template, error) {
	return template.New("").Funcs(helper.templateFuncMap).Parse(templateString)
}

func (helper *Helper) RegisterFunction(name string, function any) {
	helper.templateFuncMap[name] = function
}

func (helper *Helper) ExecuteTemplate(templateString string, data any) (string, error) {
	outputBuf := bytes.NewBuffer(nil)
	tmpl, ok := helper.templateCache[templateString]
	if !ok {
		newTemplate, err := template.New("").Funcs(helper.templateFuncMap).Parse(templateString)
		if err != nil {
			return "", err
		}
		helper.templateCache[templateString] = newTemplate
		tmpl = newTemplate
	}
	err := tmpl.Execute(outputBuf, data)
	if err != nil {
		return "", err
	}
	return outputBuf.String(), nil
}

func (helper *Helper) toTime(date any) time.Time {
	switch date := date.(type) {
	case time.Time:
		return date.In(helper.timeLocation)
	case *time.Time:
		return (*date).In(helper.timeLocation)
	case int64:
		return time.Unix(date, 0).In(helper.timeLocation)
	case int:
		return time.Unix(int64(date), 0).In(helper.timeLocation)
	case int32:
		return time.Unix(int64(date), 0).In(helper.timeLocation)
	case *timestamppb.Timestamp:
		return date.AsTime().In(helper.timeLocation)
	}
	return time.Now().Local()
}
