//
//  swift_scaling_enigmaApp.swift
//  swift-scaling-enigma
//
//  Created by Peter Bishop on 3/16/25.
//

import SwiftUI
import SwiftData

@main
struct swift_scaling_enigmaApp: App {
    var sharedModelContainer: ModelContainer = {
        let schema = Schema([
            Item.self,
        ])
        let modelConfiguration = ModelConfiguration(schema: schema, isStoredInMemoryOnly: false)

        do {
            return try ModelContainer(for: schema, configurations: [modelConfiguration])
        } catch {
            fatalError("Could not create ModelContainer: \(error)")
        }
    }()

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
        .modelContainer(sharedModelContainer)
    }
}
