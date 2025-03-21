//
//  UserViewModel.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/20/25.
//

import Foundation
import SwiftUI
import Foundation
import Observation
import CryptoKit

@Observable class UserViewModel {
    var user: UserModel = UserModel()
    var users: [UserModel] = []
    var error: String? = nil
    var baseURL: String = "http://localhost:8080"
    
    func createNewUser() async -> Bool {
        print("Creating a new user.")
        guard let url = URL(string: "\(baseURL)/users/new") else { return false }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let body: [String: Any] = [
            "name": user.name,
            "email": user.email,
            "password": user.password,
        ]
        
        guard let jsonData = try? JSONSerialization.data(withJSONObject: body, options: []) else { return false}
        
        request.httpBody = jsonData
        
        do {
            let (_, response) = try await URLSession.shared.data(for: request)
            
            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 201 {
                print("New user successfully added to Postgres DB.")
                return true
            } else {
                self.error = "Error creating user: \(response)"
                print(self.error ?? "Server Error")
                return false
            }
        } catch {
            self.error = "Error submitting data for new user: \(error.localizedDescription)"
            print(self.error ?? "Server Error")
            return false
        }
    }
    
    struct LoginResponse: Codable {
        let message: String
        let refreshToken: String
        let token: String
        let user: UserModel

        private enum CodingKeys: String, CodingKey {
            case message
            case refreshToken = "refresh_token"
            case token
            case user
        }
    }
    
    func Login() async -> Bool {
        print("Attempting user login...")
        guard let url = URL(string: "\(baseURL)/users/login") else { return false }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let body: [String: Any] = [
            "email": user.email,
            "password": user.password ?? ""
        ]
        
        guard let jsonData = try? JSONSerialization.data(withJSONObject: body, options: []) else { return false }
        request.httpBody = jsonData
        
        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 {
                // Print raw response data
                if let rawJson = String(data: data, encoding: .utf8) {
                    print("📢 Raw JSON Response: \(rawJson)")
                }

                let decoder = JSONDecoder()
                decoder.dateDecodingStrategy = .iso8601

                do {
                    let loginResponse = try decoder.decode(LoginResponse.self, from: data)
                    print("✅ Decoded Login Response: \(loginResponse)")

                    UserDefaults.standard.setValue(loginResponse.token, forKey: "authToken")
                    UserDefaults.standard.setValue(loginResponse.refreshToken, forKey: "refresh_token")

                    if let encodedUser = try? JSONEncoder().encode(loginResponse.user) {
                        UserDefaults.standard.setValue(encodedUser, forKey: "currentUser")
                    }

                    print("✅ User login successful. Token and user info saved.")
                    return true
                } catch {
                    print("❌ Decoding error: \(error)")
                    print("📢 Response Data: \(String(data: data, encoding: .utf8) ?? "Invalid Data")")
                    
                    // Step 2: Try decoding into a Dictionary to inspect missing fields
                    if let jsonObject = try? JSONSerialization.jsonObject(with: data, options: []) as? [String: Any] {
                        print("🧐 Parsed JSON as Dictionary: \(jsonObject)")
                    } else {
                        print("❌ Failed to parse JSON into Dictionary.")
                    }
                    
                    return false
                }
            } else {
                print("❌ Login Error: \(response)")
                return false
            }
        } catch {
            print("❌ Network Error: \(error.localizedDescription)")
            return false
        }
    }

}
