package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"k-microserv-kuber.com/orders/gateway"
	"k-microserv-kuber.com/users"
)

type userGateway struct {
	addr string
}

func NewUserGateway(addr string) *userGateway {
	return &userGateway{addr: addr}
}

func (g *userGateway) GetUser(ctx context.Context, userId int) (*users.User, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/user", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("userId", fmt.Sprintf("%d", userId))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var u *users.User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}

	return u, nil
}

func (g *userGateway) PutUser(ctx context.Context, userId int, email string) error {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/user", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("userId", fmt.Sprintf("%d", userId))
	values.Add("email", email)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", resp)
	}

	return nil
}
