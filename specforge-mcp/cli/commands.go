package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func handleConnect(args []string) {
	fs := flag.NewFlagSet("connect", flag.ExitOnError)
	server := fs.String("server", "http://localhost:8081", "MCP Server URL")
	token := fs.String("token", "", "API Token")
	project := fs.String("project", "", "Project ID")
	fs.Parse(args)

	if *token == "" || *project == "" {
		fmt.Println("Error: --token and --project are required")
		fs.Usage()
		os.Exit(1)
	}

	cfg := &Config{
		MCPServerURL: *server,
		APIToken:     *token,
		ProjectID:    *project,
	}

	if err := saveConfig(cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration saved successfully.")
	handleHandshake(nil)
}

func handleHandshake(args []string) {
	result, err := callMCP("initialize", map[string]interface{}{})
	if err != nil {
		fmt.Printf("Handshake failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Handshake successful: %s\n", string(result))
}

func handleTools(args []string) {
	result, err := callMCP("list_tools", map[string]interface{}{})
	if err != nil {
		fmt.Printf("Failed to list tools: %v\n", err)
		os.Exit(1)
	}

	var res struct {
		Tools []interface{} `json:"tools"`
	}
	if err := json.Unmarshal(result, &res); err != nil {
		fmt.Printf("Failed to parse tools: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available Tools:")
	for _, t := range res.Tools {
		b, _ := json.MarshalIndent(t, "  ", "  ")
		fmt.Println(string(b))
	}
}

func handleCreateSnapshot(args []string) {
	cfg, _ := loadConfig()
	params := map[string]interface{}{
		"name": "create_snapshot",
		"arguments": map[string]interface{}{
			"project_id": cfg.ProjectID,
		},
	}
	result, err := callMCP("call_tool", params)
	if err != nil {
		fmt.Printf("Failed to create snapshot: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Snapshot created: %s\n", string(result))
}

func handlePostSnapshot(args []string) {
	cfg, _ := loadConfig()
	params := map[string]interface{}{
		"name": "post_snapshot",
		"arguments": map[string]interface{}{
			"project_id": cfg.ProjectID,
		},
	}
	result, err := callMCP("call_tool", params)
	if err != nil {
		fmt.Printf("Failed to post snapshot: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Snapshot posted: %s\n", string(result))
}

func handleImportProject(args []string) {
	cfg, _ := loadConfig()
	params := map[string]interface{}{
		"name": "import_project",
		"arguments": map[string]interface{}{
			"project_id": cfg.ProjectID,
		},
	}
	result, err := callMCP("call_tool", params)
	if err != nil {
		fmt.Printf("Failed to import project: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Import triggered: %s\n", string(result))
}

func handleVerify(args []string) {
	handleHandshake(nil)
}
