package github

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Release struct {
    TagName string `json:"tag_name"`
    Name    string `json:"name"`
    HTMLURL string `json:"html_url"`
}

func LatestRelease(ctx context.Context, ownerRepo string) (*Release, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", ownerRepo)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil { return nil, err }
    req.Header.Set("Accept", "application/vnd.github+json")
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode == 404 { return nil, nil }
    if resp.StatusCode >= 300 { return nil, fmt.Errorf("github api error: %s", resp.Status) }
    var r Release
    if err := json.NewDecoder(resp.Body).Decode(&r); err != nil { return nil, err }
    return &r, nil
}


