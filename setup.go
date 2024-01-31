package main

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed packages/*
var content embed.FS

func main() {

	fmt.Println("¡Bienvenido!")

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener la ruta del directorio actual:", err)
		return
	}

	fmt.Println("El proyecto se creará en esta carpeta: ", currentDir)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("¿Desea continuar? (y/n): ")
	permission, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}

	for strings.TrimSpace(permission) != "y" {
		if strings.TrimSpace(permission) == "n" {
			os.Exit(1)
		}

		fmt.Println("Caracter no válido...")
		fmt.Print("(y/n): ")
		permission, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			return
		}
	}

	fmt.Print("Ingresa el nombre del proyecto: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}

	projectDir := filepath.Join(currentDir, strings.TrimSpace(name))
	if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
		fmt.Println("Error al crear la carpeta del proyecto:", err)
		os.Exit(1)
	}

	mainContent, err := content.ReadFile("packages/main.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	mainPath := filepath.Join(projectDir, "main.go")
	if err := writeToFile(mainPath, string(mainContent)); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		os.Exit(1)
	}

	envContent, err := content.ReadFile("packages/env.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	envPath := filepath.Join(projectDir, ".env")
	if err := writeToFile(envPath, string(envContent)); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		os.Exit(1)
	}

	gitignorePath := filepath.Join(projectDir, ".gitignore")
	file, err := os.Create(gitignorePath)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		os.Exit(1)
	}
	defer file.Close()

	serverDir := filepath.Join(projectDir, "server")
	if err := os.MkdirAll(serverDir, os.ModePerm); err != nil {
		fmt.Println("Error al crear la subcarpeta:", err)
		os.Exit(1)
	}

	serverContent, err := content.ReadFile("packages/server.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	serverPath := filepath.Join(serverDir, "server.go")
	if err := writeToFile(serverPath, string(serverContent)); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		os.Exit(1)
	}

	repositoryDir := filepath.Join(projectDir, "repository")
	if err := os.MkdirAll(repositoryDir, os.ModePerm); err != nil {
		fmt.Println("Error al crear la subcarpeta:", err)
		os.Exit(1)
	}

	repositoryContent, err := content.ReadFile("packages/repository.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	repositoryPath := filepath.Join(repositoryDir, "repository.go")
	if err := writeToFile(repositoryPath, string(repositoryContent)); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		os.Exit(1)
	}

	databaseDir := filepath.Join(projectDir, "database")
	if err := os.MkdirAll(databaseDir, os.ModePerm); err != nil {
		fmt.Println("Error al crear la subcarpeta:", err)
		os.Exit(1)
	}

	databaseContent, err := content.ReadFile("packages/database.txt")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	databasePath := filepath.Join(databaseDir, "database.go")
	if err := writeToFile(databasePath, string(databaseContent)); err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "mod", "init", fmt.Sprintf("github.com/%v", strings.TrimSpace(name)))
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error al ejecutar go mod init:", err)
		os.Exit(1)
	}

	cmd = exec.Command("go", "get", "github.com/gorilla/mux", "github.com/joho/godotenv", "github.com/lib/pq", "github.com/rs/cors")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error al instalar dependencias:", err)
		os.Exit(1)
	}

	fmt.Println("--------------------------")
	fmt.Println("Configuración completada.")
	fmt.Println("--------------------------")
}

func writeToFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
