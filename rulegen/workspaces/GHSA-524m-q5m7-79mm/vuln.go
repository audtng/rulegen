package main

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin:       func(r *http.Request) bool { return true }, // allow multi-domain
	EnableCompression: true,                                       // experimental compression
}

// Client is a middleman between the websocket connection and the hub.
