package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"k-microserv-kuber.com/orders/gateway"
	"k-microserv-kuber.com/products"
)

type productGateway struct {
	addr string
}

func NewProductGateway(addr string) *productGateway {
	return &productGateway{addr: addr}
}

func (g *productGateway) GetProduct(ctx context.Context, productId int) (*products.Product, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/products", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("productId", fmt.Sprintf("%d", productId))
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

	var p *products.Product
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}

	return p, nil
}

func (g *productGateway) InsertProduct(ctx context.Context, productId int, productName string, price float64) error {
	req, err := http.NewRequest(http.MethodPost, g.addr+"/products", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("productId", fmt.Sprintf("%d", productId))
	values.Add("name", productName)
	values.Add("price", fmt.Sprintf("%v", price))
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
