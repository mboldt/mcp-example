package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type input struct {
	Owner      string `json:"owner" jsonschema:"the owner of the GitHub repository"`
	Repository string `json:"repository" jsonschema:"the name of the GitHub repository"`
}

type output struct {
	Description string `json:"description" jsonschema:"the description of the GitHub repository"`
}

type githubRepoMetadataFetcher interface {
	fetch(ctx context.Context, owner, repository string) ([]byte, error)
}

type liveGithubRepoMetadataFetcher struct{}

func (f liveGithubRepoMetadataFetcher) fetch(ctx context.Context, owner, repository string) ([]byte, error) {
	res, err := http.Get("https://api.github.com/repos/" + owner + "/" + repository)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repository metadata: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository metadata: %s", res.Status)
	}

	return body, nil
}

type gitHubClient struct {
	fetcher githubRepoMetadataFetcher
}

func (c gitHubClient) repositoryMetadata(ctx context.Context, req *mcp.CallToolRequest, input input) (
	*mcp.CallToolResult,
	output,
	error,
) {
	if !validateInput(input) {
		return nil, output{}, fmt.Errorf("invalid input")
	}

	body, err := c.fetcher.fetch(ctx, input.Owner, input.Repository)
	if err != nil {
		return nil, output{}, err
	}

	var output output
	json.Unmarshal(body, &output)

	return nil, output, nil
}

func validateInput(input input) bool {
	matcher := regexp.MustCompile(`^[a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*$`)
	return matcher.MatchString(input.Owner) && matcher.MatchString(input.Repository)
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "github-metadata", Version: "v1.0.0"}, nil)
	client := gitHubClient{fetcher: liveGithubRepoMetadataFetcher{}}
	mcp.AddTool(server, &mcp.Tool{Name: "repository-metadata", Description: "Get metadata about a GitHub repository"}, client.repositoryMetadata)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
