package config

import "testing"

func TestExtractDirectClones(t *testing.T) {
	rules := []MirrorRule{
		{Rule: `path == "rust-kotlin/ashell"`, Action: RuleActionInclude},
		{Rule: `path == "https://github.com/cidverse/reposync"`, Action: RuleActionInclude},
		{Rule: `path == "LeastBit/Claude_skills_zh-CN"`, Action: RuleActionInclude},
		{Rule: `path == "https://github.com/mukul975/Anthropic-Cybersecurity-Skills"`, Action: RuleActionInclude},
		{Rule: `path == "travisvn/awesome-claude-skills"`, Action: RuleActionInclude},
	}

	clones := ExtractDirectClones(rules)

	if len(clones) != 2 {
		t.Errorf("expected 2 direct clones, got %d", len(clones))
	}

	if clones[0].URL != "https://github.com/cidverse/reposync" {
		t.Errorf("expected URL https://github.com/cidverse/reposync, got %s", clones[0].URL)
	}
	if clones[0].Namespace != "cidverse" {
		t.Errorf("expected namespace cidverse, got %s", clones[0].Namespace)
	}
	if clones[0].Name != "reposync" {
		t.Errorf("expected name reposync, got %s", clones[0].Name)
	}

	if clones[1].URL != "https://github.com/mukul975/Anthropic-Cybersecurity-Skills" {
		t.Errorf("expected URL https://github.com/mukul975/Anthropic-Cybersecurity-Skills, got %s", clones[1].URL)
	}
	if clones[1].Namespace != "mukul975" {
		t.Errorf("expected namespace mukul975, got %s", clones[1].Namespace)
	}
	if clones[1].Name != "Anthropic-Cybersecurity-Skills" {
		t.Errorf("expected name Anthropic-Cybersecurity-Skills, got %s", clones[1].Name)
	}
}

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		url       string
		namespace string
		name      string
	}{
		{"https://github.com/owner/repo", "owner", "repo"},
		{"https://github.com/owner/repo.git", "owner", "repo"},
		{"https://github.com/org/project-name", "org", "project-name"},
	}

	for _, tt := range tests {
		namespace, name := parseGitHubURL(tt.url)
		if namespace != tt.namespace {
			t.Errorf("parseGitHubURL(%s) namespace = %s, want %s", tt.url, namespace, tt.namespace)
		}
		if name != tt.name {
			t.Errorf("parseGitHubURL(%s) name = %s, want %s", tt.url, name, tt.name)
		}
	}
}

func TestIsURLRule(t *testing.T) {
	tests := []struct {
		rule   string
		isURL  bool
	}{
		{`path == "owner/repo"`, false},
		{`path == "https://github.com/owner/repo"`, true},
		{`path == "http://github.com/owner/repo"`, true},
	}

	for _, tt := range tests {
		if got := isURLRule(tt.rule); got != tt.isURL {
			t.Errorf("isURLRule(%s) = %v, want %v", tt.rule, got, tt.isURL)
		}
	}
}
