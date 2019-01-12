// visgostruct
// CLI application to extract structs in golang sources,
// and draw relations in PlantUML format.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

// FieldInformation is meta inforamtions of fields in struct.
type FieldInformation struct {
	Name    string
	Type    string
	HasA    string
	Tag     string
	Comment string
}

// StructInformation is meta information of struct.
type StructInformation struct {
	Name   string
	Fileds []*FieldInformation
}

// PrettyPrint is debugging function to print struct informations.
func (i StructInformation) PrettyPrint() {
	fmt.Println(i.Name)
	for _, field := range i.Fileds {
		if len(field.Comment) > 0 {
			fmt.Printf("\t%s\n", field.Comment)
		}
		fmt.Printf("\t%s\t%s\t%s\n", field.Name, field.Type, field.Tag)
	}
}

// SprintClass returns string of struct in PlantUML format.
func (i StructInformation) SprintClass(includeFields, enableComment, enableTag, byNote bool) string {
	if len(i.Fileds) == 0 {
		return ""
	}
	var uml string
	uml += fmt.Sprintf("class %s {\n", i.Name)
	if includeFields {
		for _, field := range i.Fileds {
			uml += fmt.Sprintf("{field} +%s <%s>", field.Name, field.Type)
			if !byNote {
				if enableTag && len(field.Tag) > 0 {
					uml += fmt.Sprintf(" `%s`", field.Tag)
				}
				if enableComment && len(field.Comment) > 0 {
					uml += fmt.Sprintf(" %s\n", field.Comment)
				}
			}
			uml += fmt.Sprintln()
		}
	}
	uml += fmt.Sprintf("}\n")
	if includeFields && byNote && (enableComment || enableTag) {
		uml += fmt.Sprintf("note right of %s\n", i.Name)
		for _, field := range i.Fileds {
			uml += fmt.Sprintf("%s:", field.Name)
			if enableComment {
				uml += fmt.Sprintf(" %s", field.Comment)
			}
			if enableTag {
				uml += fmt.Sprintf(" `%s`", field.Tag)
			}
			uml += fmt.Sprintln()
		}
		uml += fmt.Sprintf("end note\n")
	}
	return uml
}

// SprintRelations returns string of relations in PlantUML format.
func (i StructInformation) SprintRelations(classes map[string]*StructInformation) string {
	var uml string
	relations := map[string]int{}
	for _, field := range i.Fileds {
		_, ok := classes[field.HasA]
		if ok {
			relations[field.HasA]++
		}
	}
	for relation, count := range relations {
		if count > 1 {
			uml += fmt.Sprintf("%s --* \"%d\" %s\n", relation, count, i.Name)
		} else {
			uml += fmt.Sprintf("%s --* %s\n", relation, i.Name)
		}
	}
	return uml
}

func main() {
	app := cli.NewApp()
	app.Name = "visgostruct"
	app.Usage = "extract structs in golang sources and draw relations in PlantUML format."
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "include, i",
			Usage: "pattern to extract structs",
			Value: "",
		},
		cli.StringFlag{
			Name:  "exclude, e",
			Usage: "pattern to ignore structs",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "fields, f",
			Usage: "include field definitions",
		},
		cli.BoolFlag{
			Name:  "comment, c",
			Usage: "enable comment",
		},
		cli.BoolFlag{
			Name:  "tag, t",
			Usage: "enable tag",
		},
		cli.BoolFlag{
			Name:  "note, n",
			Usage: "comment and tag shown in note",
		},
		cli.StringFlag{
			Name:  "root, r",
			Usage: "extract recursively from specified struct(not implemented)",
		},
	}

	app.Action = func(context *cli.Context) error {
		// precompile regexp for include/exclude
		var include, exclude *regexp.Regexp
		if len(context.String("include")) > 0 {
			include = regexp.MustCompile(context.String("include"))
		}
		if len(context.String("exclude")) > 0 {
			exclude = regexp.MustCompile(context.String("exclude"))
		}

		// extract all structs found in sources
		classes := map[string]*StructInformation{}
		for _, arg := range context.Args() {
			structs := ParseFile(arg)
			for _, info := range structs {
				if exclude != nil && exclude.MatchString(info.Name) {
					continue
				}
				if include != nil && !include.MatchString(info.Name) {
					continue
				}
				classes[info.Name] = info
			}
		}

		// collect root and descendent if "root" flag is given
		if len(context.String("root")) > 0 {
			info, ok := classes[context.String("root")]
			if !ok {
				log.Fatal(fmt.Errorf("root struct %s not found", context.String("root")))
				os.Exit(1)
			}
			structs := []*StructInformation{info}
			selected := map[string]*StructInformation{}
			for len(structs) > 0 {
				children := []*StructInformation{}
				for _, info := range structs {
					selected[info.Name] = info
					for _, field := range info.Fileds {
						child, ok := classes[field.HasA]
						if ok {
							children = append(children, child)
						}
					}
				}
				structs = children
			}
			classes = selected
		}

		// print PlantUML
		fmt.Println("@startuml{}")
		fmt.Println("left to right direction")
		for _, info := range classes {
			fmt.Print(info.SprintClass(context.Bool("fields"), context.Bool("comment"), context.Bool("tag"), context.Bool("note")))
		}
		for _, info := range classes {
			fmt.Print(info.SprintRelations(classes))
		}
		fmt.Println("@enduml")
		return nil
	}
	app.Run(os.Args)
}

// ParseFile parse source file to extract structs
func ParseFile(path string) []*StructInformation {
	informations := []*StructInformation{}
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, path, nil, parser.ParseComments)

	ast.Inspect(f, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}
		if decl.Tok != token.TYPE {
			return true
		}
		for _, spec := range decl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			information := &StructInformation{}
			information.Name = ts.Name.Name
			information.Fileds = []*FieldInformation{}
			for _, field := range st.Fields.List {
				var tag string
				if field.Tag != nil {
					tag = field.Tag.Value
				}
				ft, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}
				info := &FieldInformation{}
				if len(field.Names) > 0 {
					info.Name = field.Names[0].Name
				} else {
					info.Name = ft.Name
				}
				info.Type = ft.Name
				info.HasA = strings.Trim(ft.Name, "*[]")
				info.Tag = strings.Trim(tag, "`")
				var commentString string
				if field.Comment != nil {
					for _, comment := range field.Comment.List {
						commentString += comment.Text
					}
				}
				info.Comment = commentString
				information.Fileds = append(information.Fileds, info)
			}
			informations = append(informations, information)
		}
		return true
	})
	return informations
}
