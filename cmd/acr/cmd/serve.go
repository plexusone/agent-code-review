package cmd

import (
	"context"
	"fmt"

	"github.com/plexusone/agent-code-review/internal/mcp"
	"github.com/plexusone/agent-code-review/pkg/config"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as an MCP server",
	Long: `Run agent-code-review as an MCP (Model Context Protocol) server.

This allows Claude Code or other MCP-compatible clients to use code review tools.

The server communicates over stdin/stdout using JSON-RPC.

Example:
  acr serve

Add to Claude Code MCP config:
  {
    "mcpServers": {
      "code-review": {
        "command": "acr",
        "args": ["serve"]
      }
    }
  }`,
	RunE: runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := config.CreateClient(ctx)
	if err != nil {
		return fmt.Errorf("creating client: %w", err)
	}

	server := mcp.NewServer(client)
	return server.Run(ctx)
}
