import SwiftUI

struct AuthFlowView: View {
    @State private var isRegister = false

    var body: some View {
        Group {
            if isRegister {
                RegisterPage(onToggleToLogin: { isRegister = false })
            } else {
                LoginPage(onToggleToRegister: { isRegister = true })
            }
        }
    }
}


