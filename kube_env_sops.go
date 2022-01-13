package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
    "path/filepath"
	"log"
)

const yml =
`apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
secretGenerator:
- name: environment-secrets
  envs:
  - .env.local
generatorOptions:
  disableNameSuffixHash: true
`

// Returns the current working directory
func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd="./"
	}
	return cwd
}

// Returns whether or not the path exists
func path_exists(file_path string) bool {
    _, err := os.Stat(file_path)
	if err != nil {
		return false
	}
    return true
}

// Creates a new file with data
func create_file_with_data(file_path string, data string) {
    f, err := os.Create(file_path)

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    _, err2 := f.WriteString(data)

    if err2 != nil {
        log.Fatal(err2)
    }
}

// Removes a file by file path
func remove_file(file_path string) {
	os.Remove(file_path)
	// TODO: Print on debug
	// fmt.Println("Removed file: " + file_path)
    // err := os.Remove(file_path)
	// if err != nil {
    //     fmt.Println(err)
    // }
}

// Checks if a command exists
func command_exists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Execute a command
func exec_command(cwd_path string, name string, args ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Dir = cwd_path
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()
	err := cmd.Run()

	if err != nil {
		return (fmt.Sprint(err) + ": " + stderr.String()), err
	}

	return stdout.String(), err
}

func main() {

	// Ensure that the required commands are installed
	required_commands := []string{"kubectl", "sops"}
	
	for _, command := range required_commands {
		if(!command_exists(command)) {
			log.Fatal(command + " is required. Please install.")
		}
	}

	// Set the path that files should be generated in
	cwd_path := flag.String("cwd_path", getCwd(), "The path to generate the encrypted secrets.")
	dot_env_enc_file_name := flag.String("dot_env_enc_file_name", ".env-enc.yml", "An optional name for the encrypted yml file.")	
	save_secret := flag.Bool("save", true, "If set to false, the encrypted secret will be output to stdout.")

    // Parse command line into the defined flags
    flag.Parse()
  	
	kustomization_file_path := filepath.Join(*cwd_path, "kustomization.yaml")
	dot_env_local_file_path := filepath.Join(*cwd_path, ".env.local")	
	dot_env_dec_file_path := filepath.Join(*cwd_path, ".env-dec.yml")
	dot_env_enc_file_path := filepath.Join(*cwd_path, *dot_env_enc_file_name)

	// Delay the removal of the file after main is complete
	defer remove_file(kustomization_file_path)
	defer remove_file(dot_env_dec_file_path)

	// TODO: Print if debug mode
	// Print the input parameters
	// fmt.Println("Debug Info: ")
	// fmt.Println("cwd_path: ", *cwd_path)

	// If the .env.local does not exist
	if(!path_exists(dot_env_local_file_path)) {
		log.Fatal("A .env.local does not exist. Please create a .env.local file.\n\n")
	}

	// Create the kustomization file
	create_file_with_data(kustomization_file_path, yml)

	// Execute the kustomization
	std_out, err := exec_command(*cwd_path, "kubectl", "kustomize", *cwd_path)
	
	if err != nil {
		log.Fatal(err)
	}

	// Create the .env-dec.yml file
	create_file_with_data(dot_env_dec_file_path, std_out)

	// Execute the sops encryption
	std_out2, err2 := exec_command(*cwd_path, "sops", "--encrypt", dot_env_dec_file_path)
	
	if err2 != nil {
		log.Fatal(err)
	}

	if *save_secret {
		// Create the .env-enc.yml file
		create_file_with_data(dot_env_enc_file_path, std_out2)

		// Print a success to the console
		fmt.Println("Successfully created the encrypted secret: " + *dot_env_enc_file_name)
	} else {
		fmt.Print(std_out2)
	}
}
