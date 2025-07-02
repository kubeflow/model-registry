package converter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
)

// visitor
type visitor struct {
	t        *testing.T
	entities map[string]*oapiEntity
}

func newVisitor(t *testing.T, _ *ast.File) visitor {
	return visitor{
		t: t,
		entities: map[string]*oapiEntity{
			"RegisteredModel": {
				obj: openapi.RegisteredModel{},
			},
			"ModelVersion": {
				obj: openapi.ModelVersion{},
			},
			"DocArtifact": {
				obj: openapi.DocArtifact{},
			},
			"ModelArtifact": {
				obj: openapi.ModelArtifact{},
			},
			"DataSet": {
				obj: openapi.DataSet{},
			},
			"Metric": {
				obj: openapi.Metric{},
			},
			"Parameter": {
				obj: openapi.Parameter{},
			},
			"ServingEnvironment": {
				obj: openapi.ServingEnvironment{},
			},
			"InferenceService": {
				obj: openapi.InferenceService{},
			},
			"ServeModel": {
				obj: openapi.ServeModel{},
			},
			"Artifact": {
				obj: openapi.Artifact{},
			},
			"Experiment": {
				obj: openapi.Experiment{},
			},
			"ExperimentRun": {
				obj: openapi.ExperimentRun{},
			},
		},
	}
}

func (v *visitor) extractGroup(regex *regexp.Regexp, s string) string {
	extracted := regex.FindStringSubmatch(s)
	if len(extracted) != 2 {
		v.t.Errorf("unable to extract groups from %s for %s", regex.String(), s)
	}
	// the first one is the wole matched string, the second one is the group
	return extracted[1]
}

func (v *visitor) getEntity(name string) *oapiEntity {
	val, ok := v.entities[name]
	if !ok {
		v.t.Errorf("openapi entity not found in the entities map: %s", name)
	}
	return val
}

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.InterfaceType:
		for _, m := range d.Methods.List {
			methodName := m.Names[0].Name

			if converterMethodPattern.MatchString(methodName) {
				entityName := v.extractGroup(converterMethodPattern, methodName)
				entity := v.getEntity(entityName)
				// there should be just one doc comment matching ignoreDirectivePattern
				for _, c := range m.Doc.List {
					if ignoreDirectivePattern.MatchString(c.Text) {
						entity.notEditableFields = v.extractGroup(ignoreDirectivePattern, c.Text)
					}
				}
			} else if overrideNotEditableMethodPattern.MatchString(methodName) {
				entityName := v.extractGroup(overrideNotEditableMethodPattern, methodName)
				entity := v.getEntity(entityName)
				// there should be just one doc comment matching ignoreDirectivePattern
				for _, c := range m.Doc.List {
					if ignoreDirectivePattern.MatchString(c.Text) {
						entity.ignoredFields = v.extractGroup(ignoreDirectivePattern, c.Text)
					}
				}
			}
		}
		v.checkEntities()
	}
	return v
}

// checkEntities check if all editable fields are listed in the goverter ignore directive of OverrideNotEditableFor
func (v *visitor) checkEntities() {
	errorMsgs := map[string][]string{}
	for k, v := range v.entities {
		msgs := checkEntity(v)
		if len(msgs) > 0 {
			errorMsgs[k] = msgs
		}
	}

	if len(errorMsgs) > 0 {
		missingFieldsMsg := ""
		for k, fields := range errorMsgs {
			missingFieldsMsg += fmt.Sprintf("%s: %v\n", k, fields)
		}
		v.t.Errorf("missing fields to be ignored for OverrideNotEditableFor* goverter methods:\n%v", missingFieldsMsg)
	}
}

// checkEntity check if there are missing fields to be ignored in the override method
func checkEntity(entity *oapiEntity) []string {
	res := []string{}
	objType := reflect.TypeOf(entity.obj)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if !strings.Contains(entity.notEditableFields, field.Name) && !strings.Contains(entity.ignoredFields, field.Name) {
			// check if the not editable field (first check) is not present in the ignored fields (second check)
			// if this condition is true, we missed that field in the Override method ignore list
			res = append(res, field.Name)
		}
	}
	return res
}

// test

var (
	converterMethodPattern           *regexp.Regexp = regexp.MustCompile(`Convert(?P<entity>\w+)Update`)
	overrideNotEditableMethodPattern *regexp.Regexp = regexp.MustCompile(`OverrideNotEditableFor(?P<entity>\w+)`)
	ignoreDirectivePattern           *regexp.Regexp = regexp.MustCompile(`// goverter:ignore (?P<fields>.+)`)
)

func setup(t *testing.T) *assert.Assertions {
	return assert.New(t)
}

func TestOverrideNotEditableFields(t *testing.T) {
	_ = setup(t)

	fset := token.NewFileSet() // positions are relative to fset
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("error getting current working directory")
	}
	filePath := fmt.Sprintf("%s/openapi_converter.go", wd)
	f, _ := parser.ParseFile(fset, filePath, nil, parser.ParseComments)

	v := newVisitor(t, f)
	ast.Walk(v, f)
}

type oapiEntity struct {
	obj               any
	notEditableFields string
	ignoredFields     string
}
