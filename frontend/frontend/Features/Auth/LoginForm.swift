import SwiftUI

struct LoginForm: View {
    @Binding var email: String
    @Binding var password: String

    var body: some View {
        Section(header: Text("Sign In")) {
            TextField("Email", text: $email)
                .keyboardType(.emailAddress)
                .textContentType(.username)
                .autocapitalization(.none)
            SecureField("Password", text: $password)
                .textContentType(.password)
        }
    }
}


