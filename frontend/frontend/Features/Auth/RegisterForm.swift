import SwiftUI

struct RegisterForm: View {
    @Binding var email: String
    @Binding var password: String
    @Binding var first: String

    var body: some View {
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
    }
}

struct RegisterForm_Previews: PreviewProvider {
    @State static var email = "test@example.com"
    @State static var password = "password"
    @State static var first = "Adam"
    
    static var previews: some View {
        Form {
            RegisterForm(
                email: $email, password: $password, first: $first
            )
        }
    }
}


