import SwiftUI

struct RegisterForm: View {
    @Binding var email: String
    @Binding var password: String
    @Binding var username: String
    @Binding var first: String
    @Binding var last: String

    var body: some View {
        Section(header: Text("Account")) {
            TextField("Email", text: $email)
                .keyboardType(.emailAddress)
                .textContentType(.username)
                .autocapitalization(.none)
            SecureField("Password", text: $password)
                .textContentType(.newPassword)
            TextField("Username", text: $username)
                .textContentType(.nickname)
        }
        Section(header: Text("Profile")) {
            TextField("First name", text: $first)
                .textContentType(.givenName)
            TextField("Last name", text: $last)
                .textContentType(.familyName)
        }
    }
}


