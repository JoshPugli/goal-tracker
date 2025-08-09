import SwiftUI

@main
struct frontendApp: App {
    let persistenceController = PersistenceController.shared
    @StateObject private var auth = AuthManager.shared

    var body: some Scene {
        WindowGroup {
            RootView()
                .environment(\.managedObjectContext, persistenceController.container.viewContext)
                .environmentObject(auth)
        }
    }
}


