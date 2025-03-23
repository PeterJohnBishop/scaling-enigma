//
//  ChatViewModel.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/23/25.
//

import SwiftUI
import Observation

@Observable class ChatMessageViewModel {
    
    func startNewChat(users: [UserModel]) async -> Bool {
        // create a new chat object and add user ids to the users array
        // check that user.length >= 2
        // save chat object to postgres db
        // update each user with the chat id returned from postgres
        return true
    }
    
    func getChats() async -> Bool {
        return true
    }
    
    func getChatById() async -> Bool {
        return true
    }
    
    func upDateChat(chat: ChatModel) async -> Bool {
        return true
    }
    
    func sendMessage(chat: ChatModel) async -> Bool {
        // create a new message object and add content
        // check that either text or media content != empty before sending
        // update the current chat with the message
        // emit websocket message to prompt all users to update message feed
        return true
    }
    
    func getMessages() async -> Bool {
        return true
    }
    
    func deleteMessage(chat: ChatModel, messageId: String) async -> Bool {
        // update chat to remove the messageId
        // emit websocket message to prompt all users to update message feed
        // delete message by messageId
        return true
    }
    
    func deleteChat(chatId: String) async -> Bool {
        return true
    }
    
    
    
    
}
