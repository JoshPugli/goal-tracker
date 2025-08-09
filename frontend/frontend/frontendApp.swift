//
//  frontendApp.swift
//  frontend
//
//  Created by callum on 2025-08-09.
//

import SwiftUI

@main
struct frontendApp: App {
    let persistenceController = PersistenceController.shared

    var body: some Scene {
        WindowGroup {
            HomeView()
                .environment(\.managedObjectContext, persistenceController.container.viewContext)
        }
    }
}
