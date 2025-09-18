package matrix

import (
    "context"
    "fmt"
    "net/http"
    "strings"
    "time"
)

type Client struct {
    Homeserver  string
    AccessToken string
    httpClient  *http.Client
}

func NewClient(hs, token string) *Client {
    return &Client{
        Homeserver: hs,
        AccessToken: token,
        httpClient: &http.Client{Timeout: 15 * time.Second},
    }
}

func (c *Client) SendMessage(ctx context.Context, roomID, body string) error {
    if c == nil || c.AccessToken == "" || c.Homeserver == "" { return fmt.Errorf("matrix not configured") }
    url := strings.TrimRight(c.Homeserver, "/") + "/_matrix/client/v3/rooms/" + roomID + "/send/m.room.message"
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(fmt.Sprintf(`{"msgtype":"m.text","body":"%s"}`, body)))
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+c.AccessToken)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 { return fmt.Errorf("matrix send failed: %s", resp.Status) }
    return nil
}


