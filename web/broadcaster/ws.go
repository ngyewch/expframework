package broadcaster

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ngyewch/expframework/codec"
)

type WebSocketBroadcaster[T any, PT *T] struct {
	codec          codec.Codec[T, PT]
	pingInterval   time.Duration
	subProtocols   []string
	originPatterns []string
	connections    mapset.Set[*websocket.Conn]
	mutex          sync.Mutex
}

func NewWebSocketBroadcaster[T any, PT *T](codec codec.Codec[T, PT], pingInterval time.Duration, subProtocols []string, originPatterns []string) *WebSocketBroadcaster[T, PT] {
	return &WebSocketBroadcaster[T, PT]{
		codec:          codec,
		pingInterval:   pingInterval,
		subProtocols:   subProtocols,
		originPatterns: originPatterns,
		connections:    mapset.NewSet[*websocket.Conn](),
	}
}

func (ws *WebSocketBroadcaster[T, PT]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   ws.subProtocols,
		OriginPatterns: ws.originPatterns,
	})
	if err != nil {
		slog.Error("could not accept websocket",
			slog.Any("err", err),
		)
		return
	}

	ws.addConnection(c)

	ctx := c.CloseRead(r.Context())
	for {
		err = c.Ping(ctx)
		if err != nil {
			if !errors.Is(ctx.Err(), context.Canceled) {
				slog.Error("could not ping websocket, removing",
					slog.Any("err", err),
				)
			}
			ws.removeConnection(c)
			break
		}
		select {
		case <-r.Context().Done():
			return
		case <-time.After(ws.pingInterval):
		}
	}
}

func (ws *WebSocketBroadcaster[T, PT]) addConnection(c *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.connections.Add(c)
}

func (ws *WebSocketBroadcaster[T, PT]) removeConnection(c *websocket.Conn) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.connections.Remove(c)
}

func (ws *WebSocketBroadcaster[T, PT]) Publish(ctx context.Context, m PT) error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	b, err := ws.codec.MarshalBinary(m)
	if err != nil {
		return err
	}
	for _, c := range ws.connections.ToSlice() {
		err := c.Write(ctx, websocket.MessageBinary, b)
		if err != nil {
			slog.Error("could not send message",
				slog.Any("err", err),
			)
			go func(c *websocket.Conn) {
				ws.removeConnection(c)
			}(c)
		}
	}
	return nil
}
