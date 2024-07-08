package cmd

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator struct {
	WorkDir                 string // exec命令执行工作目录
	Module                  string // go.mod module
	InternalDir             string
	ServiceDir              string
	RepositoryDir           string
	HandlerDir              string
	ModelDir                string
	ConfigDir               string
	DatabaseDir             string
	Type                    string
	CamelName               string
	LowerCamelName          string
	SnakeName               string
	FuncName                string
	InjectHereImpl          string //
	InjectInterface         string //
	InjectInterfaceFuncHere string //
	InjectInterfaceImpl     string //
	InjectHandlerFunc       string
	TargetFuncFile          string
	TargetInterfaceFile     string
	InjectInterfaceEntity   bytes.Buffer
	TemplateFunc            bytes.Buffer
	TemplateInterface       bytes.Buffer
	tmpl                    *template.Template
}

var gen *Generator

func (gen *Generator) NewEntSchema() (string, error) {
	return runCmdCommand("go", gen.InternalDir, []string{
		"run",
		"-mod=mod",
		"entgo.io/ent/cmd/ent",
		"new",
		gen.CamelName,
	}...)
}
func (gen *Generator) setVariable() {
	// dirname
	gen.InternalDir = filepath.Join(gen.WorkDir, "internal")
	gen.ServiceDir = filepath.Join(gen.InternalDir, "service")
	gen.RepositoryDir = filepath.Join(gen.InternalDir, "repository")
	gen.HandlerDir = filepath.Join(gen.InternalDir, "httptransport")
	gen.DatabaseDir = filepath.Join(gen.InternalDir, "database")
	gen.ModelDir = filepath.Join(gen.InternalDir, "model")
	gen.ConfigDir = filepath.Join(gen.InternalDir, "config")
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
func (gen *Generator) NewGenerate() error {
	if gen.Type != TypeHandler {
		interfaceData, err := os.ReadFile(gen.TargetInterfaceFile)
		if err != nil {
			return err
		}

		newInterfaceData := strings.Replace(string(interfaceData), "// {{ .InjectInterface }}", gen.TemplateInterface.String(), 1)
		err = WriteToFile(gen.TargetInterfaceFile, []byte(newInterfaceData))

		if err != nil {
			return err
		}
	}

	if !fileIsExisted(gen.TargetFuncFile) {
		create := bytes.Buffer{}
		create.WriteString("package " + gen.Type)
		create.WriteString("\r\n")

		if gen.Type != TypeHandler {
			create.WriteString("import \"context\"")
		} else {
			create.WriteString(`
import (
   "github.com/gin-gonic/gin"
   "net/http"
)
`)
		}

		create.WriteString("\r\n")
		create.Write(gen.TemplateFunc.Bytes())
		return WriteToFile(gen.TargetFuncFile, create.Bytes())
	}

	implData, err := os.ReadFile(gen.TargetFuncFile)
	if err != nil {
		return err
	}

	replaceOld := "// {{ .InjectInterfaceImpl }}"
	if gen.Type == TypeHandler {
		replaceOld = "// {{ .InjectHandlerFunc }}"
	}

	newImplData := strings.Replace(string(implData), replaceOld, gen.TemplateFunc.String(), 1)
	return WriteToFile(gen.TargetFuncFile, []byte(newImplData))
}
