import SwiftUI

struct LoginPage: View {
    @EnvironmentObject var auth: AuthManager
    @State private var email = ""
    @State private var password = ""
    @State private var loading = false
    @State private var errorMsg: String?
    let onToggleToRegister: () -> Void

    var body: some View {
        let isDisabled = (loading || email.isEmpty || password.isEmpty)

        NavigationStack {
            ZStack {
                Color.clear
                VStack(spacing: 16) {
                    Spacer(minLength: 120)
                    Form {
                        Section {
                            TextField("email", text: $email)
                                .keyboardType(.emailAddress)
                                .textContentType(.username)
                                .autocapitalization(.none)
                            SecureField("password", text: $password)
                                .textContentType(.password)
                        }
                        if let msg = errorMsg { Text(msg).foregroundStyle(.red) }
                    }
                    .scrollContentBackground(.hidden)
                    .frame(maxWidth: 360)
                    .background()
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .center)
            .safeAreaInset(edge: .bottom) {
                VStack(spacing: 8) {
                    Button("need an account? register", action: onToggleToRegister)
                        .buttonStyle(.plain)
                        .foregroundStyle(.blue)

                    let isDisabled = (loading || email.isEmpty || password.isEmpty)
                    Button(action: submit) {
                        Text("sign in").frame(maxWidth: .infinity)
                    }
                    .disabled(isDisabled)
                    .padding(.vertical, 24)
                    .foregroundStyle(isDisabled ? Color.white.opacity(0.7) : Color.white)
                    .background(isDisabled ? Color.blue.opacity(0.4) : Color.blue)
                    .clipShape(RoundedRectangle(cornerRadius: 24))
                }
                .padding(.horizontal, 32)
            }
            .navigationTitle("Sign In")
        }
    }

    private func submit() {
        loading = true
        errorMsg = nil
        Task {
            do {
                try await auth.login(email: email, password: password)
            } catch {
                errorMsg = "Authentication failed. Please check your details."
            }
            loading = false
        }
    }
}

struct LoginPage_Previews: PreviewProvider {
    static var previews: some View {
        LoginPage(onToggleToRegister: { })
            .environmentObject(AuthManager.shared)
    }
}


