// This program creates a Linux container using system calls
// instead of a separate tool (like Docker). The idea is to
// see what functionality Linux provides by itself.
//
// This program MUST BE RUN on a Linux machine.
//
// Sources:
//  "Linux Containers and Virtualization: A Kernel Perspective" by Shashank Mohan Jain (p.93-106)
//  https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909
//	https://medium.com/@jain.sm/writing-your-own-linux-container-259054465bd1
//
// To compile:
// 		$ GOOS=linux go build -o tc tinyContainer.go
//
// You must compile for Linux (notice GOOS=linux), it will not compile
// on Mac or Windows without it.
//
// Usage:
//		$ pwd
//		/root/tiny-container
//		$ ./tc -root=/root/tiny-container -shell=./ts
//		Creating container...
//		Tiny shell started
//		$ ls
//		...
//
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var root, shell, xaction string

func init() {
	flag.StringVar(&root, "root", "", "Full path of directory to mount as root in the container. Required.")
	flag.StringVar(&shell, "shell", "", "Path to shell program to run (relative to root once mounted). Required.")
	flag.StringVar(&xaction, "x-action", "create", "Used internally. Please ignore.")
	flag.Parse()
}

func main() {

	// Root and shell are required args
	if root == "" || shell == "" {
		printHelp()
		os.Exit(1)
	}

	if xaction == "create" {
		// Create the wrapper for the container
		createContainer(root, shell)
		os.Exit(0)
	}

	if xaction == "launch-shell" {
		// This is used internally to lanch the shell inside
		// the container. Therefore this command must be exectuted
		// AFTER the container has been created.
		runShell(root, shell)
		os.Exit(0)
	}

	printHelp()
	os.Exit(1)
}

func createContainer(root string, shell string) {
	fmt.Printf("Creating container...\n")

	args := []string{"-root=" + root, "-shell=" + shell, "-x-action=launch-shell"}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// These flags are what instruct Linux to create a new container
	// (notice NEWNS, NEWUTS, ...) as it runs the command.
	var flags uintptr
	flags = syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: flags,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
	}
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /proc/self/exe container - %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Exited container\n")
}

func runShell(root string, shell string) {
	fmt.Printf("Launching shell session...\n")
	fmt.Printf("\troot.: %s\n", root)
	fmt.Printf("\tshell: %s\n", shell)

	cmd := exec.Command(shell)

	cmd.Env = []string{"tiny_demo=something tiny"}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set the hostname
	err := syscall.Sethostname([]byte("tinyhost"))
	if err != nil {
		fmt.Printf("Error setting hostname - %s\n", err)
	}

	// Pivot to our new root folder
	err = pivotRoot(root)
	if err != nil {
		fmt.Printf("Error running pivot_root - %s\n", err)
		os.Exit(1)
	}

	// Launch the new shell session
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running the shell %s - %s\n", shell, err)
		os.Exit(1)
	}

	fmt.Printf("Exited shell session\n")
}

func pivotRoot(newRoot string) error {
	putold := filepath.Join(newRoot, "/.pivot_root")

	// Bind mount `newroot` to itself.
	// This is a slight hack needed to satisfy the `pivot_root`
	// requirement that `newroot` and `putold` must not be on
	// the same filesystem as the current root
	err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return err
	}

	// create putold directory
	err = os.MkdirAll(putold, 0700)
	if err != nil {
		return err
	}

	// call pivot_root
	err = syscall.PivotRoot(newRoot, putold)
	if err != nil {
		return err
	}

	// ensure current working directory is set to new root
	err = os.Chdir("/")
	if err != nil {
		return err
	}

	//umount putold, which now lives at /.pivot_root
	putold = "/.pivot_root"
	err = syscall.Unmount(putold, syscall.MNT_DETACH)
	if err != nil {
		return err
	}

	// remove putold
	err = os.RemoveAll(putold)
	if err != nil {
		return err
	}
	return nil
}

func printHelp() {
	fmt.Println("tinyContainer (tc) parameters:")
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("")
	fmt.Println("  $ pwd")
	fmt.Println("  /root/tiny-container")
	fmt.Println("  $ ./tc -root=/root/tiny-container -shell=/ts")
	fmt.Println("")
}
