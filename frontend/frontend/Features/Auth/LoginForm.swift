import SwiftUI

struct LoginForm: View {
    @Binding var email: String
    @Binding var password: String

    var body: some View {
        Section() {
            TextField("Email", text: $email)
                .keyboardType(.emailAddress)
                .textContentType(.username)
                .autocapitalization(.none)
            SecureField("Password", text: $password)
                .textContentType(.password)
        }
    }
}

struct LoginForm_Previews: PreviewProvider {
    @State static var email = "test@example.com"
    @State static var password = "password"
    @State static var loading = false

    static var previews: some View {
        NavigationStack {
            Form {
                LoginForm(email: $email, password: $password)

                Button(action: { /* simulate sign in */ }) {
                    if loading { ProgressView() } else { Text("Sign In") }
                }
                .disabled(email.isEmpty || password.isEmpty)

                Button("Need an account? Register") { /* simulate toggle */ }
                    .buttonStyle(.borderless)
            }
            .navigationTitle("Sign In")
        }
    }
}
