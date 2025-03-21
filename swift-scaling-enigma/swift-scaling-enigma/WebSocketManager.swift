//
//  WebSocketManager.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/17/25.
//

import Foundation

struct Event: Codable {
    let type: String
    let message: String
}

class WebSocketManager: ObservableObject {
    private var webSocketTask: URLSessionWebSocketTask?
    private let urlSession = URLSession(configuration: .default)
    
    @Published var messages: [Event] = [] // Stores received messages

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

    func sendMessage(_ event: Event) {
        guard let jsonData = try? JSONEncoder().encode(event),
              let jsonString = String(data: jsonData, encoding: .utf8) else {
            print("Failed to encode message")
            return
        }
        let messageToSend = URLSessionWebSocketTask.Message.string(jsonString)
        webSocketTask?.send(messageToSend) { error in
            if let error = error {
                print("Error sending event: \(error.localizedDescription)")
            } else {
                print("📤 Sent: \(jsonString)")
            }
        }
    }

    func receiveMessage() {
        webSocketTask?.receive { [weak self] result in
            switch result {
            case .success(let message):
                switch message {
                case .string(let text):
                    if let data = text.data(using: .utf8),
                       let event = try? JSONDecoder().decode(Event.self, from: data) {
                        DispatchQueue.main.async {
                            self?.messages.append(event)
                        }
                        print("Received: \(event)")
                    } else {
                        print("Failed to decode response")
                    }
                default:
                    print("Received unsupported message type")
                }
                self?.receiveMessage() // Keep listening
            case .failure(let error):
                print("Error receiving websocket message: \(error.localizedDescription)")
            }
        }
    }

    func disconnect() {
        webSocketTask?.cancel(with: .goingAway, reason: nil)
        print("WebSocket disconnected")
    }
}

