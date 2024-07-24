import asyncio
import websockets


# Funkcja obsługująca połączenie z klientem
async def handle_connection(websocket, path):
    print(f"Client connected: {websocket.remote_address}")
    try:
        async for message in websocket:
            print(f"Received message from {websocket.remote_address}: {message}")
    except websockets.ConnectionClosed as e:
        print(f"Client disconnected: {websocket.remote_address} ({e.code}, {e.reason})")
    finally:
        print(f"Connection closed: {websocket.remote_address}")


# Funkcja główna uruchamiająca serwer
async def main():
    async with websockets.serve(handle_connection, "localhost", 8765):
        print("Server started at ws://localhost:8765")
        await asyncio.Future()  # Run forever


# Uruchomienie pętli zdarzeń
asyncio.run(main())
