import SwiftUI

struct RegisterPage: View {
    @EnvironmentObject var auth: AuthManager
    @State private var email = ""
    @State private var password = ""
    @State private var first = ""
    @State private var loading = false
    @State private var errorMsg: String?
    let onToggleToLogin: () -> Void

    var body: some View {
        NavigationStack {
            Form {
                Section(header: Text("Account")) {
                    TextField("First name", text: $first)
                        .textContentType(.givenName)
                    TextField("Email", text: $email)
                        .keyboardType(.emailAddress)
                        .textContentType(.username)
                        .autocapitalization(.none)
                    SecureField("Password", text: $password)
                        .textContentType(.newPassword)
                }

                if let msg = errorMsg { Text(msg).foregroundStyle(.red) }

                Button(action: submit) {
                    if loading { ProgressView() } else { Text("Create Account") }
                }
                .disabled(loading || email.isEmpty || password.isEmpty || first.isEmpty)

                Button("Have an account? Sign in", action: onToggleToLogin)
                    .buttonStyle(.borderless)
            }
            .navigationTitle("Register")
        }
    }

    private func submit() {
        loading = true
        errorMsg = nil
        Task {
            do {
                let uname: String = {
                    let base = first.trimmingCharacters(in: .whitespacesAndNewlines)
                    if !base.isEmpty { return base.replacingOccurrences(of: "\\s+", with: "", options: .regularExpression).lowercased() }
                    return email.split(separator: "@").first.map(String.init) ?? "user"
                }()
                try await auth.register(.init(email: email, username: uname, first_name: first, last_name: "", password: password))
            } catch {
                errorMsg = "Registration failed. Please check your details."
            }
            loading = false
        }
    }
}

struct RegisterPage_Previews: PreviewProvider {
    static var previews: some View {
        RegisterPage(onToggleToLogin: { })
            .environmentObject(AuthManager.shared)
    }
}


