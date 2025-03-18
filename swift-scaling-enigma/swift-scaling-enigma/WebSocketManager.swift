//
//  WebSocketManager.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/17/25.
//

import Foundation

class WebSocketManager: ObservableObject {
    private var webSocketTask: URLSessionWebSocketTask?
    private let urlSession = URLSession(configuration: .default)
    
    @Published var messages: [String] = [] // Stores received messages

    init() {
        connect()
    }

    func connect() {
        guard let url = URL(string: "ws://localhost:8080/ws") else { return }
        webSocketTask = urlSession.webSocketTask(with: url)
        webSocketTask?.resume()
        print("Connected to WebSocket server")
        
        receiveMessage() // Start listening for messages
    }

    func sendMessage(_ message: String) {
        let messageToSend = URLSessionWebSocketTask.Message.string(message)
        webSocketTask?.send(messageToSend) { error in
            if let error = error {
                print("Error sending event: \(error.localizedDescription)")
            } else {
                print("📤 Sent: \(message)")
            }
        }
    }

    func receiveMessage() {
        webSocketTask?.receive { [weak self] result in
            switch result {
            case .success(let message):
                switch message {
                case .string(let text):
                    DispatchQueue.main.async {
                        self?.messages.append(text)
                    }
                    print("Received: \(text)")
                default:
                    print("Received unsupported message type")
                }
                self?.receiveMessage() // Keep listening
            case .failure(let error):
                print("Error recieving websocket message: \(error.localizedDescription)")
            }
        }
    }

    func disconnect() {
        webSocketTask?.cancel(with: .goingAway, reason: nil)
        print("WebSocket disconnected")
    }
}
