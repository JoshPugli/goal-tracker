import SwiftUI

struct RootView: View {
    @EnvironmentObject var auth: AuthManager
    var body: some View {
        Group {
            if auth.token == nil { AuthFlowView() } else { HomeView() }
        }
    }
}


