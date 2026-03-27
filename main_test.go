package main

import (
	"context"
	"fmt"
	"testing"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name  string
		input input
		valid bool
	}{
		{"valid input", input{Owner: "octocat", Repository: "Hello-World"}, true},
		{"missing owner", input{Owner: "", Repository: "Hello-World"}, false},
		{"missing repository", input{Owner: "octocat", Repository: ""}, false},
		{"invalid owner", input{Owner: "octocat!", Repository: "Hello-World"}, false},
		{"invalid repository", input{Owner: "octocat", Repository: "Hello/ World"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateInput(tt.input)
			if validateInput(tt.input) != tt.valid {
				t.Errorf("expected valid: %v, got: %v", tt.valid, valid)
			}
		})
	}
}

func TestGetGithubRepositoryMetadata(t *testing.T) {
	tests := []struct {
		name    string
		fetcher githubRepoMetadataFetcher
		want    output
		wantErr bool
	}{
		{
			name:    "successful fetch",
			fetcher: NewMockGithubRepoMetadataFetcher(cannedResponse, nil),
			want:    output{Description: "Source for https://www.alphaminnesota.org"},
			wantErr: false,
		},
		{
			name:    "failed fetch",
			fetcher: NewMockGithubRepoMetadataFetcher("", fmt.Errorf("failed to fetch")),
			want:    output{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := gitHubClient{fetcher: tt.fetcher}
			_, got, err := client.repositoryMetadata(context.Background(), nil, input{Owner: "mboldt", Repository: "alphaminnesota.org"})
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("expected: %v, got: %v", tt.want, got)
			}
		})
	}
}

type mockMetadataFetcher struct {
	body []byte
	err  error
}

func (m mockMetadataFetcher) fetch(ctx context.Context, owner, repository string) ([]byte, error) {
	return m.body, m.err
}

func NewMockGithubRepoMetadataFetcher(body string, err error) githubRepoMetadataFetcher {
	return mockMetadataFetcher{body: []byte(body), err: err}
}

// From `curl https://api.github.com/repos/mboldt/alphaminnesota.org“
const cannedResponse = `
{
  "id": 1157494063,
  "node_id": "R_kgDORP31Lw",
  "name": "alphaminnesota.org",
  "full_name": "mboldt/alphaminnesota.org",
  "private": false,
  "owner": {
    "login": "mboldt",
    "id": 2256266,
    "node_id": "MDQ6VXNlcjIyNTYyNjY=",
    "avatar_url": "https://avatars.githubusercontent.com/u/2256266?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/mboldt",
    "html_url": "https://github.com/mboldt",
    "followers_url": "https://api.github.com/users/mboldt/followers",
    "following_url": "https://api.github.com/users/mboldt/following{/other_user}",
    "gists_url": "https://api.github.com/users/mboldt/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/mboldt/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/mboldt/subscriptions",
    "organizations_url": "https://api.github.com/users/mboldt/orgs",
    "repos_url": "https://api.github.com/users/mboldt/repos",
    "events_url": "https://api.github.com/users/mboldt/events{/privacy}",
    "received_events_url": "https://api.github.com/users/mboldt/received_events",
    "type": "User",
    "user_view_type": "public",
    "site_admin": false
  },
  "html_url": "https://github.com/mboldt/alphaminnesota.org",
  "description": "Source for https://www.alphaminnesota.org",
  "fork": false,
  "url": "https://api.github.com/repos/mboldt/alphaminnesota.org",
  "forks_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/forks",
  "keys_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/keys{/key_id}",
  "collaborators_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/collaborators{/collaborator}",
  "teams_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/teams",
  "hooks_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/hooks",
  "issue_events_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/issues/events{/number}",
  "events_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/events",
  "assignees_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/assignees{/user}",
  "branches_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/branches{/branch}",
  "tags_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/tags",
  "blobs_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/git/blobs{/sha}",
  "git_tags_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/git/tags{/sha}",
  "git_refs_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/git/refs{/sha}",
  "trees_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/git/trees{/sha}",
  "statuses_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/statuses/{sha}",
  "languages_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/languages",
  "stargazers_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/stargazers",
  "contributors_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/contributors",
  "subscribers_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/subscribers",
  "subscription_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/subscription",
  "commits_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/commits{/sha}",
  "git_commits_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/git/commits{/sha}",
  "comments_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/comments{/number}",
  "issue_comment_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/issues/comments{/number}",
  "contents_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/contents/{+path}",
  "compare_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/compare/{base}...{head}",
  "merges_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/merges",
  "archive_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/{archive_format}{/ref}",
  "downloads_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/downloads",
  "issues_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/issues{/number}",
  "pulls_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/pulls{/number}",
  "milestones_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/milestones{/number}",
  "notifications_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/notifications{?since,all,participating}",
  "labels_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/labels{/name}",
  "releases_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/releases{/id}",
  "deployments_url": "https://api.github.com/repos/mboldt/alphaminnesota.org/deployments",
  "created_at": "2026-02-13T22:09:51Z",
  "updated_at": "2026-03-04T15:23:26Z",
  "pushed_at": "2026-03-04T15:13:34Z",
  "git_url": "git://github.com/mboldt/alphaminnesota.org.git",
  "ssh_url": "git@github.com:mboldt/alphaminnesota.org.git",
  "clone_url": "https://github.com/mboldt/alphaminnesota.org.git",
  "svn_url": "https://github.com/mboldt/alphaminnesota.org",
  "homepage": null,
  "size": 52,
  "stargazers_count": 0,
  "watchers_count": 0,
  "language": "HTML",
  "has_issues": true,
  "has_projects": true,
  "has_downloads": true,
  "has_wiki": true,
  "has_pages": true,
  "has_discussions": false,
  "forks_count": 0,
  "mirror_url": null,
  "archived": false,
  "disabled": false,
  "open_issues_count": 0,
  "license": {
    "key": "mit",
    "name": "MIT License",
    "spdx_id": "MIT",
    "url": "https://api.github.com/licenses/mit",
    "node_id": "MDc6TGljZW5zZTEz"
  },
  "allow_forking": true,
  "is_template": false,
  "web_commit_signoff_required": false,
  "has_pull_requests": true,
  "pull_request_creation_policy": "all",
  "topics": [

  ],
  "visibility": "public",
  "forks": 0,
  "open_issues": 0,
  "watchers": 0,
  "default_branch": "main",
  "temp_clone_token": null,
  "network_count": 0,
  "subscribers_count": 0
}
`
