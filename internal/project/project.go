package project

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"locky/internal/resources"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"text/template"
)

type Config struct {
	Name string
	Open bool
}

type gradleSettings struct {
	Name string
}

type gradleBuild struct {
	KotlinVersion, KtorVersion, MainClass string
}

func GenServer(config Config, logger io.Writer) error {
	projectPath := filepath.Join(".", config.Name)
	fmt.Fprintf(logger, "Creating project directory %s\n", projectPath)
	err := CreateDir(projectPath)

	if err != nil {
		return err
	}

	fmt.Fprintln(logger, "Generating project files")
	srcDir := filepath.Join(projectPath, "src", "main", "kotlin")
	err = os.MkdirAll(srcDir, 0744)

	if err != nil {
		return err
	}

	fns := []func() error {
		func() error {
			return Write("server.kt", filepath.Join(srcDir, "main.kt"))
		},

		func() error {
			return WriteByTemplate("build.gradle", gradleBuild {
				KotlinVersion: "1.4.21",
				KtorVersion: "1.5.0",
				MainClass: "MainKt",
			}, filepath.Join(projectPath, "build.gradle"))
		},

		func() error {
			return WriteByTemplate("settings.gradle", gradleSettings{config.Name}, filepath.Join(projectPath, "settings.gradle"))
		},

		func() error {
			return Write("gradle.properties", filepath.Join(projectPath, "gradle.properties"))
		},
	}

	for _, fn := range fns {
		err := fn()

		if err != nil {
			return err
		}
	}

	if config.Open {
		fmt.Fprintln(os.Stderr, "Opening project in IDEA")
		_ = openInIdea(projectPath)
	}

	fmt.Fprintln(logger, "Downloading Gradle wrapper")
	cmd := exec.Command("/usr/bin/env", "gradle", "wrapper")
	cmd.Dir = projectPath
	err = cmd.Run()

	return err
}

func openInIdea(path string) error {
	usr, err := user.Current()

	if err != nil {
		return err
	}

	ideaPath, ok := findIdea(filepath.Join(usr.HomeDir, "Library"))

	if !ok {
		return errors.New("cannot find path for IDEA executable")
	}

	cmd := exec.Command("nohup", ideaPath, ".")
	cmd.Dir = path
	return cmd.Start()
}

func findIdea(path string) (string, bool) {
	var result string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "idea" && info.Mode() & 0111 != 0 {
			result = path
			return io.EOF
		}

		return nil
	})

	return result, err == io.EOF
}

func Write(resPath string, outputFile string) error {
	b, err := resources.GetBin(resPath)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputFile, *b, 0644)
}

func WriteByTemplate(resPath string, data interface{}, outputFile string) error {
	s, err := resources.Get(resPath)

	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("").Parse(s))

	f, err := os.Create(outputFile)

	if err != nil {
		return err
	}

	defer f.Close()
	err = tmpl.Execute(f, data)

	return err
}

func CreateDir(name string) error {
	err := os.Mkdir(name, 0744)

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create project directory %s", name))
	}

	return nil
}
