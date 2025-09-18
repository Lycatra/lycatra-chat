//go:build dev

package main

import (
    "context"
    "fmt"
    "time"
    cfg "github.com/Lycatra/lycatra-chat/internal/config"
    gh "github.com/Lycatra/lycatra-chat/internal/github"
    mx "github.com/Lycatra/lycatra-chat/internal/matrix"
)

// Minimal demo: on start, fetch latest release of this repo and post to Matrix if configured
func init() {
    go func() {
        time.Sleep(500 * time.Millisecond)
        conf := cfg.FromEnv()
        if conf.Matrix.AccessToken == "" || conf.Matrix.Homeserver == "" || conf.Matrix.RoomID == "" {
            return
        }
        rel, err := gh.LatestRelease(context.Background(), "Lycatra/lycatra-chat")
        if err != nil || rel == nil { return }
        client := mx.NewClient(conf.Matrix.Homeserver, conf.Matrix.AccessToken)
        _ = client.SendMessage(context.Background(), conf.Matrix.RoomID, fmt.Sprintf("Latest release: %s (%s) %s", rel.Name, rel.TagName, rel.HTMLURL))
    }()
}


