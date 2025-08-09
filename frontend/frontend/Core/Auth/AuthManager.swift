import Foundation
import Security
import Combine

final class AuthManager: ObservableObject {
    static let shared = AuthManager()

    @Published private(set) var token: String?

    let baseURL = URL(string: "https://48458bf86bcf.ngrok-free.app")!

    private let keychainService = "grindhouse.auth"
    private let keychainAccount = "jwt"

    private init() { self.token = loadToken() }

    func setToken(_ newToken: String?) {
        token = newToken
        if let t = newToken { _ = saveToken(t) } else { _ = deleteToken() }
    }

    struct LoginRequest: Codable { let email: String; let password: String }
    struct RegisterRequest: Codable { let email: String; let username: String; let first_name: String; let last_name: String; let password: String }
    struct AuthResponse: Codable { let token: String; let user: [String: AnyCodable]? }

    func login(email: String, password: String) async throws {
        let url = baseURL.appendingPathComponent("/api/auth/login")
        var req = URLRequest(url: url)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(LoginRequest(email: email, password: password))
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard let http = resp as? HTTPURLResponse, (200..<300).contains(http.statusCode) else { throw URLError(.userAuthenticationRequired) }
        let decoded = try JSONDecoder().decode(AuthResponse.self, from: data)
        setToken(decoded.token)
    }

    func register(_ r: RegisterRequest) async throws {
        let url = baseURL.appendingPathComponent("/api/auth/register")
        var req = URLRequest(url: url)
        req.httpMethod = "POST"
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        req.httpBody = try JSONEncoder().encode(r)
        let (data, resp) = try await URLSession.shared.data(for: req)
        guard let http = resp as? HTTPURLResponse, (200..<300).contains(http.statusCode) else { throw URLError(.userAuthenticationRequired) }
        let decoded = try JSONDecoder().decode(AuthResponse.self, from: data)
        setToken(decoded.token)
    }

    func logout() { setToken(nil) }
    func attachAuth(to request: inout URLRequest) { if let t = token { request.setValue("Bearer \(t)", forHTTPHeaderField: "Authorization") } }

    private func saveToken(_ token: String) -> Bool {
        let data = Data(token.utf8)
        let query: [String: Any] = [kSecClass as String: kSecClassGenericPassword, kSecAttrService as String: keychainService, kSecAttrAccount as String: keychainAccount, kSecValueData as String: data]
        SecItemDelete(query as CFDictionary)
        return SecItemAdd(query as CFDictionary, nil) == errSecSuccess
    }
    private func loadToken() -> String? {
        let query: [String: Any] = [kSecClass as String: kSecClassGenericPassword, kSecAttrService as String: keychainService, kSecAttrAccount as String: keychainAccount, kSecReturnData as String: true, kSecMatchLimit as String: kSecMatchLimitOne]
        var item: CFTypeRef?
        guard SecItemCopyMatching(query as CFDictionary, &item) == errSecSuccess, let data = item as? Data else { return nil }
        return String(data: data, encoding: .utf8)
    }
    private func deleteToken() -> Bool { let q: [String: Any] = [kSecClass as String: kSecClassGenericPassword, kSecAttrService as String: keychainService, kSecAttrAccount as String: keychainAccount]; return SecItemDelete(q as CFDictionary) == errSecSuccess }
}

struct AnyCodable: Codable {}


