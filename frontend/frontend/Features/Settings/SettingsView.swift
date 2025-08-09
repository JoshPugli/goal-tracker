import SwiftUI

struct SettingsView: View {
    @EnvironmentObject var auth: AuthManager
    @Environment(\.dismiss) private var dismiss
    @State private var remindersEnabled = true

    var body: some View {
        NavigationStack {
            Form {
                Section(header: Text("General")) {
                    Toggle("Enable reminders", isOn: $remindersEnabled)
                    NavigationLink("Manage habits") {
                        Text("Coming soon")
                            .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .center)
                    }
                }

                Section(header: Text("Account")) {
                    Button(role: .destructive) {
                        auth.logout()
                        dismiss()
                    } label: {
                        Text("Log Out")
                    }
                }
            }
            .navigationTitle("Settings")
        }
    }
}


