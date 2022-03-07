package main

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/viper"
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

type nuspecConfs struct {
	Id                       string
	Version                  string
	Title                    string
	Authors                  string
	Owners                   string
	RequireLicenseAcceptance string
	Description              string
	Summary                  string
	Tags                     string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.tmpl"))
}

func viperConfigVariable(key string) string {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading the configuration file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}

func archive() {
	idConf := viperConfigVariable("id")
	archiveFormat := ".zip"
	// ToDo:
	// fmt.SprintF for map strings, so we can insert id name from var into the string
	tpp := idConf
	rtpp := fmt.Sprintf("content/%s.nuspec", tpp)
	rtppf := fmt.Sprintf("%s.nuspec", tpp)
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		//"content/bolt_exec_puppet.nuspec": "bolt_exec_puppet.nuspec",
		rtpp:                          rtppf,
		"content/[Content_Types].xml": "[Content_Types].xml",
		"content/_rels/.rels":         "_rels/.rels",
		"content/package/services/metadata/core-properties/81fb83d7949f4e33baf8f5b203521668.psmdcp": "package/services/metadata/core-properties/81fb83d7949f4e33baf8f5b203521668.psmdcp",
		"tools/bolt_exec_puppet.zip":    "tools/bolt_exec_puppet.zip",
		"tools/chocolateyinstall.ps1":   "tools/chocolateyinstall.ps1",
		"tools/chocolateyuninstall.ps1": "tools/chocolateyuninstall.ps1",
	})
	if err != nil {
		log.Fatal(err)
	}

	// create the output file we'll write to
	out, err := os.Create(idConf + archiveFormat)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Archival: archiver.Zip{},
	}

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Loading Configuration keys for .nuspec file
	idConfig := viperConfigVariable("id")
	versionConfig := viperConfigVariable("version")
	titleConfig := viperConfigVariable("title")
	authorsConfig := viperConfigVariable("authors")
	ownersConfig := viperConfigVariable("owners")
	requireLicenseConfig := viperConfigVariable("requireLicenseAcceptance")
	descriptionConfig := viperConfigVariable("description")
	summaryConfig := viperConfigVariable("summary")
	tagsConfig := viperConfigVariable("tags")

	// Creation of .nuspec file
	data := idConfig
	response := fmt.Sprintf("content/%s.nuspec", data)
	nuspecFile, err := os.Create(response)
	if err != nil {
		log.Fatal("error creating nuspec file", err)
	}
	defer nuspecFile.Close()

	// Creation of .rels file
	relsFile, err := os.Create("content/_rels/.rels")
	if err != nil {
		log.Fatal("error creating .rels file", err)
	}
	defer relsFile.Close()

	// Creation of coreproperties file
	cPropFile, err := os.Create("content/package/services/metadata/core-properties/81fb83d7949f4e33baf8f5b203521668.psmdcp")
	if err != nil {
		log.Fatal("error creating core-properties file", err)
	}
	defer cPropFile.Close()

	// Merging data from Configuration file to Template Engine
	nuspecData := nuspecConfs{
		Id:                       idConfig,
		Version:                  versionConfig,
		Title:                    titleConfig,
		Authors:                  authorsConfig,
		Owners:                   ownersConfig,
		RequireLicenseAcceptance: requireLicenseConfig,
		Description:              descriptionConfig,
		Summary:                  summaryConfig,
		Tags:                     tagsConfig,
	}

	relsData := nuspecConfs{
		Id: idConfig,
	}

	cPropData := nuspecConfs{
		Authors:     authorsConfig,
		Description: descriptionConfig,
		Id:          idConfig,
		Version:     versionConfig,
		Tags:        tagsConfig,
		Title:       titleConfig,
	}

	// Executing nuspec template file
	err = tpl.ExecuteTemplate(nuspecFile, "nuspec.tmpl", nuspecData)
	if err != nil {
		log.Fatal("error executing nuspec.tmpl", err)
	} else {
		log.Println("Successfully executed nuspec.tmpl")
	}

	// Executing rels template file
	err = tpl.ExecuteTemplate(relsFile, "rels.tmpl", relsData)
	if err != nil {
		log.Fatal("error executing rels.tmpl", err)
	} else {
		log.Println("Successfully executed rels.tmpl")
	}

	// Executing core-properties template file
	err = tpl.ExecuteTemplate(cPropFile, "coreproperties.tmpl", cPropData)
	if err != nil {
		log.Fatal("error executing coreproperties.tmpl", err)
	} else {
		log.Println("Successfully executed coreproperties.tmpl")
	}

	archive()
}
