package wsstat

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// MeasureLatency is a wrapper around a one-hit usage of the WSStat instance. It establishes a
// WebSocket connection, sends a message, reads the response, and closes the connection.
// Note: sets all times in the Result object.
func MeasureLatency(url *url.URL, msg string, customHeaders http.Header) (Result, []byte, error) {
	ws := New()
	defer ws.Close()

	if err := ws.Dial(url, customHeaders); err != nil {
		logger.Debug().Err(err).Msg("Failed to establish WebSocket connection")
		return Result{}, nil, err
	}
	ws.WriteMessage(websocket.TextMessage, []byte(msg))
	_, p, err := ws.ReadMessage()
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to read message")
		return Result{}, nil, err
	}
	ws.Close()

	return *ws.Result, p, nil
}

// MeasureLatencyJSON is a wrapper around a one-hit usage of the WSStat instance. It establishes a
// WebSocket connection, sends a JSON message, reads the response, and closes the connection.
// Note: sets all times in the Result object.
func MeasureLatencyJSON(url *url.URL, v interface{}, customHeaders http.Header) (Result, interface{}, error) {
	ws := New()
	defer ws.Close()

	if err := ws.Dial(url, customHeaders); err != nil {
		logger.Debug().Err(err).Msg("Failed to establish WebSocket connection")
		return Result{}, nil, err
	}
	p, err := ws.OneHitMessageJSON(v)
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to send message")
		return Result{}, nil, err
	}
	ws.Close()

	return *ws.Result, p, nil
}

// MeasureLatencyPing is a wrapper around a one-hit usage of the WSStat instance. It establishes a
// WebSocket connection, sends a ping message, awaits the pong response, and closes the connection.
// Note: sets all times in the Result object.
func MeasureLatencyPing(url *url.URL, customHeaders http.Header) (Result, error) {
	ws := New()
	defer ws.Close()

	if err := ws.Dial(url, customHeaders); err != nil {
		logger.Debug().Err(err).Msg("Failed to establish WebSocket connection")
		return Result{}, err
	}
	err := ws.PingPong()
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to send ping")
		return Result{}, err
	}
	ws.Close()

	return *ws.Result, nil
}
