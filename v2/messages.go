package courier

import (
	"context"
	"encoding/json"
	"net/http"
)

// ProvidersChannelResponse represents the channel section of the ProvidersResponse
type ProvidersChannelResponse struct {
	Key      string
	Name     string
	Template string
}

// ProvidersResponse represents the providers section of the MessageResponse
type ProvidersResponse struct {
	Channel          *ProvidersChannelResponse
	Error            string
	Status           string
	Delivered        int64
	Sent             int64
	Clicked          int64
	Provider         string
	ProviderResponse interface{}
	Reference        interface{}
}

// MessageResponse represents the return of the /messages/* endpoints on the Courier API
type MessageResponse struct {
	ID           string
	Event        string
	Notification string
	Status       string
	Error        string
	Reason       string
	Recipient    string
	Enqueued     int64
	Delivered    int64
	Sent         int64
	Clicked      int64
	Providers    []*ProvidersResponse
	Tags         []string
}

type MessagesResponse struct {
	Paging  *PagingResponse
	Results []*MessageResponse
}

// GetMessage calls the /messages/:id endpoint of the Courier API
func (c *Client) GetMessage(ctx context.Context, messageID string) (*MessageResponse, error) {
	var response MessageResponse
	err := c.API.SendRequestWithJSON(ctx, "GET", "/messages/"+messageID, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetMessages calls the /messages/ endpoint of the Courier API
func (c *Client) GetMessages(ctx context.Context, cursor, tags string) (*MessagesResponse, error) {
	url := c.API.BaseURL + "/messages"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if cursor != "" {
		q.Add("cursor", cursor)
	}
	if tags != "" {
		q.Add("tags", tags)
	}
	req.URL.RawQuery = q.Encode()

	bytes, err := c.API.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	var data MessagesResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
