package main

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"strings"
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
	ZipPath                  string
}

func init() {
	mkDir()
	createTmplFiles()
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
	zipConf := viperConfigVariable("zipPath")
	archiveFormat := ".nupkg"

	// Formating string for .nuspec file in order to use the provided id inside the config.yaml file
	tpp := idConf
	rtpp := fmt.Sprintf("content/%s.nuspec", tpp)
	rtppf := fmt.Sprintf("%s.nuspec", tpp)

	// Formating string for using your zipped binary name from the config.yaml file
	zbc := zipConf
	zbcc := fmt.Sprintf("tools/%s", zbc)
	zbccf := fmt.Sprintf("tools/%s", zbc)
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		rtpp:                          rtppf,
		"content/[Content_Types].xml": "[Content_Types].xml",
		"content/_rels/.rels":         "_rels/.rels",
		"content/package/services/metadata/core-properties/81fb83d7949f4e33baf8f5b203521668.psmdcp": "package/services/metadata/core-properties/81fb83d7949f4e33baf8f5b203521668.psmdcp",
		zbcc:                            zbccf,
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

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func mkDir() {
	err := os.Mkdir("content", 0755)
	check(err)
	err = os.Mkdir("content/_rels", 0755)
	check(err)
	err = os.Mkdir("content/package", 0755)
	check(err)
	err = os.Mkdir("content/package/services", 0755)
	check(err)
	err = os.Mkdir("content/package/services/metadata", 0755)
	check(err)
	err = os.Mkdir("content/package/services/metadata/core-properties", 0755)
	check(err)
	err = os.Mkdir("templates", 0755)
	check(err)
}

func createTmplFiles() {
	relsStr := fmt.Sprint(`<?xml version="1.0" encoding="utf-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Type="http://schemas.microsoft.com/packaging/2010/07/manifest" Target="/{{.Id}}.nuspec" Id="R4b6b1994e8284062" />
    <Relationship Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="/package/services/metadata/core-properties/81fb83d7949f4e33baf8f5b203521668.psmdcp" Id="Rc4eb3718cc22453f" />
</Relationships>`)
	f, err := os.Create("templates/rels.tmpl")
	check(err)
	defer f.Close()
	io.Copy(f, strings.NewReader(relsStr))

	nuspecStr := fmt.Sprint(`<?xml version="1.0"?>
<package xmlns="http://schemas.microsoft.com/packaging/2010/07/nuspec.xsd">
    <metadata>
        <id>{{.Id}}</id>
        <version>{{.Version}}</version>
        <title>{{.Title}}</title>
        <authors>{{.Authors}}</authors>
        <owners>{{.Owners}}</owners>
        <requireLicenseAcceptance>{{.RequireLicenseAcceptance}}</requireLicenseAcceptance>
        <description>{{.Description}}</description>
        <summary>{{.Summary}}</summary>
        <tags>{{.Tags}}</tags>
    </metadata>
</package>`)
	g, err := os.Create("templates/nuspec.tmpl")
	check(err)
	defer g.Close()
	io.Copy(g, strings.NewReader(nuspecStr))

	coPropStr := fmt.Sprint(`<?xml version="1.0" encoding="utf-8"?>
<coreProperties xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://schemas.openxmlformats.org/package/2006/metadata/core-properties">
    <dc:creator>{{.Authors}}</dc:creator>
    <dc:description>{{.Description}}</dc:description>
    <dc:identifier>{{.Id}}</dc:identifier>
    <version>{{.Version}}</version>
    <keywords>{{.Tags}}</keywords>
    <dc:title>{{.Title}}</dc:title>
    <lastModifiedBy>choco, Version=0.11.3.0, Culture=neutral, PublicKeyToken=79d02ea9cad655eb;Microsoft Windows NT 10.0.19044.0;.NET Framework 4</lastModifiedBy>
</coreProperties>`)
	h, err := os.Create("templates/coreproperties.tmpl")
	check(err)
	defer h.Close()
	io.Copy(h, strings.NewReader(coPropStr))

	contentTypeStr := fmt.Sprint(`<?xml version="1.0" encoding="utf-8"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml" />
<Default Extension="nuspec" ContentType="application/octet" />
<Default Extension="zip" ContentType="application/octet" />
<Default Extension="ps1" ContentType="application/octet" />
<Default Extension="psmdcp" ContentType="application/vnd.openxmlformats-package.core-properties+xml" /></Types>`)
	j, err := os.Create("content/[Content_Types].xml")
	check(err)
	defer j.Close()
	io.Copy(j, strings.NewReader(contentTypeStr))
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
