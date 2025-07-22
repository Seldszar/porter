package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/sync/errgroup"
)

var (
	clients = make([]*Client, 0)
)

type Client struct {
	ctx context.Context
	url string

	conn   *websocket.Conn
	logger *slog.Logger
}

func (c *Client) run() error {
	conn, _, err := websocket.Dial(c.ctx, c.url, nil)

	if err != nil {
		return err
	}

	c.conn = conn

	for {
		mt, p, err := conn.Read(c.ctx)

		if err != nil {
			return err
		}

		c.logger.Info("Message received", slog.Any("message_type", mt), slog.Any("payload", p))

		for _, client := range clients {
			if client == c {
				continue
			}

			client.Write(mt, p)
		}
	}
}

func (c *Client) Write(mt websocket.MessageType, p []byte) error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Write(c.ctx, mt, p)
}

func (c *Client) Connect() error {
	c.logger.Info("Starting client...")

	for {
		err := c.run()

		if err != nil {
			c.logger.Error("Client disconnected", slog.Any("error", err))
		}

		time.Sleep(1000)
	}
}

func main() {
	g := new(errgroup.Group)

	for _, url := range os.Args[1:] {
		client := &Client{
			ctx:    context.Background(),
			logger: slog.With(slog.String("url", url)),
			url:    url,
		}

		g.Go(func() error {
			return client.Connect()
		})

		clients = append(clients, client)
	}

	g.Wait()
}
