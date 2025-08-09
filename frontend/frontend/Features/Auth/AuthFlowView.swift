import SwiftUI

struct AuthFlowView: View {
    @EnvironmentObject var auth: AuthManager
    @State private var email = ""
    @State private var password = ""
    @State private var first = ""
    @State private var last = ""
    @State private var username = ""
    @State private var isRegister = false
    @State private var errorMsg: String?
    @State private var loading = false

    var body: some View {
        NavigationStack {
            Form {
                if isRegister {
                    RegisterForm(email: $email, password: $password, username: $username, first: $first, last: $last)
                } else {
                    LoginForm(email: $email, password: $password)
                }

                if let msg = errorMsg { Text(msg).foregroundStyle(.red) }

                Button(action: submit) {
                    if loading { ProgressView() } else { Text(isRegister ? "Create Account" : "Sign In") }
                }
                .disabled(loading || email.isEmpty || password.isEmpty || (isRegister && username.isEmpty))

                Button(isRegister ? "Have an account? Sign in" : "Need an account? Register") {
                    isRegister.toggle()
                }
                .buttonStyle(.borderless)
            }
            .navigationTitle(isRegister ? "Register" : "Sign In")
        }
    }

    private func submit() {
        loading = true
        errorMsg = nil
        Task {
            do {
                if isRegister {
                    try await auth.register(.init(email: email, username: username, first_name: first, last_name: last, password: password))
                } else {
                    try await auth.login(email: email, password: password)
                }
            } catch {
                errorMsg = "Authentication failed. Please check your details."
            }
            loading = false
        }
    }
}


