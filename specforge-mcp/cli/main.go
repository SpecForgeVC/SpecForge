package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "connect":
		handleConnect(args)
	case "handshake":
		handleHandshake(args)
	case "tools":
		handleTools(args)
	case "create-snapshot":
		handleCreateSnapshot(args)
	case "post-snapshot":
		handlePostSnapshot(args)
	case "import-project":
		handleImportProject(args)
	case "verify":
		handleVerify(args)
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("SpecForge MCP CLI")
	fmt.Println("Usage: specforge-mcp <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  connect         Authenticate and save config")
	fmt.Println("  handshake       Check connection to MCP server")
	fmt.Println("  tools           List available MCP tools")
	fmt.Println("  create-snapshot Create a project snapshot")
	fmt.Println("  post-snapshot   Post a changelog snapshot")
	fmt.Println("  import-project  Trigger project import")
	fmt.Println("  verify          Verify connection and auth")
	fmt.Println("  help            Show this help message")
}
