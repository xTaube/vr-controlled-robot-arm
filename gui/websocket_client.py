import websocket
import threading


class WsClient:
    def __init__(self, address: str):
        self.ws = websocket.WebSocketApp(
            address,
            on_message=self.on_message,
            on_error=self.on_error,
            on_close=self.on_close,
        )
        self.ws.on_open = self.on_open
        self.ws_thread = threading.Thread(target=self.ws.run_forever)
        self.ws_thread.daemon = True
        self.ws_thread.start()

    def send_message(self, message: str) -> None:
        self.ws.send(message)

    def on_message(self, ws, message: str) -> None:
        print(f"Received message: {message}")

    def on_error(self, ws, error: str) -> None:
        print(f"Error: {error}")

    def on_close(self, ws) -> None:
        print("WebSocket closed")

    def on_open(self, ws) -> None:
        print("WebSocket connection opened")

    def close_connection(self) -> None:
        self.ws.close()
