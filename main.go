package main

import (
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

//func mkDir() {
//	err := os.Mkdir("content", 0755)
//	if err != nil {
//		log.Fatal("Error creating content Folder", err)
//	}
//
//	err := os.Mkdir("_rels", 0755)
//	if err != nil {
//		log.Fatal("Error creating _rels folder", err)
//	}
//
//	defer os.RemoveAll("content")
//}

func main() {
	//mkDir()
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
	nuspecFile, err := os.Create("content/bolt_exec_puppet.nuspec")
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
}
