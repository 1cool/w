package cmd

import (
	"bytes"
	template2 "github.com/1cool/w/template"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ActionType int

const (
	ActionCreate ActionType = iota + 1
	ActionUpdate
)

type Generator struct {
	Module         string // go.mod module
	TargetFile     string // 写入目标文件
	SnakeName      string
	TmplName       string
	UpdatePosition string       // 更新位置
	Action         ActionType   // 动作；创建，还是更新
	Content        bytes.Buffer // 内容
	WorkDir        string       // exec命令执行工作目录
	InternalDir    string
	ConfigDir      string
	DatabaseDir    string
	HandlerDir     string
	RequestDir     string
	ResponseDir    string
	ModelDir       string
	RepositoryDir  string
	ServiceDir     string
	tmpl           *template.Template
}

var gen *Generator

func (gen *Generator) NewEntSchema() (string, error) {
	return runCmdCommand("go", gen.InternalDir, []string{
		"run",
		"-mod=mod",
		"entgo.io/ent/cmd/ent",
		"new",
		strcase.ToCamel(gen.SnakeName),
	}...)
}
func (gen *Generator) setDir() {
	gen.InternalDir = filepath.Join(gen.WorkDir, "internal")
	gen.ServiceDir = filepath.Join(gen.InternalDir, "service")
	gen.RepositoryDir = filepath.Join(gen.InternalDir, "repository")
	gen.HandlerDir = filepath.Join(gen.InternalDir, "httptransport")
	gen.DatabaseDir = filepath.Join(gen.InternalDir, "database")
	gen.ModelDir = filepath.Join(gen.InternalDir, "model")
	gen.ConfigDir = filepath.Join(gen.InternalDir, "config")
	gen.RequestDir = filepath.Join(gen.HandlerDir, "request")
	gen.ResponseDir = filepath.Join(gen.HandlerDir, "response")

	os.MkdirAll(gen.InternalDir, os.ModePerm)
	os.MkdirAll(gen.RepositoryDir, os.ModePerm)
	os.MkdirAll(gen.ServiceDir, os.ModePerm)
	os.MkdirAll(gen.HandlerDir, os.ModePerm)
	os.MkdirAll(gen.RequestDir, os.ModePerm)
	os.MkdirAll(gen.ResponseDir, os.ModePerm)
	os.MkdirAll(gen.DatabaseDir, os.ModePerm)
	os.MkdirAll(gen.ModelDir, os.ModePerm)
	os.MkdirAll(gen.ConfigDir, os.ModePerm)
}

func (gen *Generator) EntGenerate() (string, error) {
	return runCmdCommand("go", gen.InternalDir, []string{
		"generate",
		"./ent",
	}...)
}

func (gen *Generator) InitHandler() error {
	filename := filepath.Join(gen.HandlerDir, "gintransport.go")
	isExist := fileIsExisted(filename)
	if isExist {
		return nil
	}

	err := os.MkdirAll(gen.HandlerDir, os.ModePerm)
	if err != nil {
		return err
	}

	// 渲染模板并将结果写入到指定文件中
	buf := bytes.Buffer{}
	err = gen.tmpl.ExecuteTemplate(&buf, "handler.tmpl", gen)
	if err != nil {
		log.Fatal(err)
	}

	return WriteToFile(filename, buf.Bytes())
}

// InitInterface init interface file
func (gen *Generator) InitInterface(t string) error {
	var (
		filename string
		dirname  string
	)
	if t == "service" {
		filename = filepath.Join(gen.ServiceDir, t+".go")
		dirname = gen.ServiceDir
	} else {
		filename = filepath.Join(gen.RepositoryDir, t+".go")
		dirname = gen.RepositoryDir
	}

	isExist := fileIsExisted(filename)
	if isExist {
		return nil
	}

	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}

	// 渲染模板并将结果写入到指定文件中
	buf := bytes.Buffer{}
	err = gen.tmpl.ExecuteTemplate(&buf, t+"_interface", gen)
	if err != nil {
		log.Fatal(err)
	}

	err = WriteToFile(filename, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// InitGenerate generate func
// filename mysql.go
// fileDirname /internal/database/
// templateName mysql.tmpl
func (gen *Generator) InitGenerate(filename, fileDirname, templateName string) error {
	newFileTarget := filepath.Join(fileDirname, filename)
	if fileIsExisted(newFileTarget) {
		return nil
	}

	err := os.MkdirAll(fileDirname, os.ModePerm)
	if err != nil {
		return err
	}

	// 渲染模板并将结果写入到指定文件中
	buf := bytes.Buffer{}
	err = gen.tmpl.ExecuteTemplate(&buf, templateName, gen)
	if err != nil {
		log.Fatal(err)
	}

	return WriteToFile(newFileTarget, buf.Bytes())
}

// NewGenerate func
// the t is new action's type.eg: service、repository、httptransport
//func (gen *Generator) NewGenerate() error {
//	if gen.Type != TypeHandler {
//		interfaceData, err := os.ReadFile(gen.TargetInterfaceFile)
//		if err != nil {
//			return err
//		}
//
//		newInterfaceData := strings.Replace(string(interfaceData), "// {{ .InjectInterface }}", gen.TemplateInterface.String(), 1)
//		err = WriteToFile(gen.TargetInterfaceFile, []byte(newInterfaceData))
//
//		if err != nil {
//			return err
//		}
//	}
//
//	if !fileIsExisted(gen.TargetFuncFile) {
//		create := bytes.Buffer{}
//		create.WriteString("package " + gen.Type)
//		create.WriteString("\r\n")
//
//		if gen.Type != TypeHandler {
//			create.WriteString("import \"context\"")
//		} else {
//			create.WriteString(`
//import (
//   "github.com/gin-gonic/gin"
//   "net/http"
//)
//`)
//		}
//
//		create.WriteString("\r\n")
//		create.Write(gen.TemplateFunc.Bytes())
//		return WriteToFile(gen.TargetFuncFile, create.Bytes())
//	}
//
//	implData, err := os.ReadFile(gen.TargetFuncFile)
//	if err != nil {
//		return err
//	}
//
//	replaceOld := "// {{ .InjectInterfaceImpl }}"
//	if gen.Type == TypeHandler {
//		replaceOld = "// {{ .InjectHandlerFunc }}"
//	}
//
//	newImplData := strings.Replace(string(implData), replaceOld, gen.TemplateFunc.String(), 1)
//	return WriteToFile(gen.TargetFuncFile, []byte(newImplData))
//}

// write
// all command finally use this function to write content to file.
func (gen *Generator) write() error {
	gen.Content = bytes.Buffer{}
	err := gen.tmpl.ExecuteTemplate(&gen.Content, gen.TmplName, gen)
	if err != nil {
		return err
	}
	if gen.Action == ActionCreate {
		return WriteToFile(gen.TargetFile, gen.Content.Bytes())
	}
	if gen.Action == ActionUpdate {
		oldFileContent, err := os.ReadFile(gen.TargetFile)
		if err != nil {
			return err
		}
		newFileContent := strings.Replace(string(oldFileContent), gen.UpdatePosition, gen.Content.String(), 1)
		return WriteToFile(gen.TargetFile, []byte(newFileContent))
	}
	return nil
}

func (gen *Generator) Generate(targetFile, tmpl string, actionType ActionType) error {
	gen.TargetFile = targetFile
	gen.TmplName = tmpl
	gen.Action = actionType
	return gen.write()
}

func NewGenerate(name string) *Generator {
	return &Generator{
		SnakeName: strcase.ToSnake(name),
		tmpl: template.Must(
			template.New("").
				Funcs(template.FuncMap{
					"ToSlash":     filepath.ToSlash,
					"FirstLetter": firstLetter,
					"CaseName":    strcase.ToCamel,
					"LowerCamel":  strcase.ToLowerCamel,
					"Snake":       strcase.ToSnake,
					"Plural":      inflection.Plural,
				}).
				ParseFS(template2.TemplateDir, "tmpl/*.tmpl"),
		),
	}
}
