//
//  UserModel.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/20/25.
//

import Foundation

struct UserModel: Codable, Equatable {
    var id: String
    var name: String
    var email: String
    var password: String?
    var created_at: Date?
    var updated_at: Date?

    private enum CodingKeys: String, CodingKey {
        case id, name, email, password, created_at, updated_at
    }

    static let iso8601Formatter: ISO8601DateFormatter = {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        return formatter
    }()
    
    init(id: String = "", name: String = "", email: String = "", password: String = "", created_at: Date? = nil, updated_at: Date? = nil) {
                self.id = id
                self.name = name
                self.email = email
                self.password = password
            self.created_at = created_at
            self.updated_at = updated_at
            }

    init(from decoder: Decoder) throws {
            let container = try decoder.container(keyedBy: CodingKeys.self)
            
            // Required fields
            id = try container.decode(String.self, forKey: .id)
            name = try container.decode(String.self, forKey: .name)
            email = try container.decode(String.self, forKey: .email)

            // Optional fields
            password = try container.decodeIfPresent(String.self, forKey: .password)

            // Decode dates safely
            if let createdAtString = try? container.decode(String.self, forKey: .created_at),
               let createdAtDate = UserModel.iso8601Formatter.date(from: createdAtString) {
                created_at = createdAtDate
            } else {
                created_at = nil
            }

            if let updatedAtString = try? container.decode(String.self, forKey: .updated_at),
               let updatedAtDate = UserModel.iso8601Formatter.date(from: updatedAtString) {
                updated_at = updatedAtDate
            } else {
                updated_at = nil
            }
        }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(id, forKey: .id)
        try container.encode(name, forKey: .name)
        try container.encode(email, forKey: .email)
        try container.encode(password, forKey: .password) 
        if let created_at = created_at {
            try container.encode(UserModel.iso8601Formatter.string(from: created_at), forKey: .created_at)
        }
        if let updated_at = updated_at {
            try container.encode(UserModel.iso8601Formatter.string(from: updated_at), forKey: .updated_at)
        }
    }
}
