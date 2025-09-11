import Foundation
import Combine

// MARK: - Authentication Service
class AuthService: ObservableObject {
    static let shared = AuthService()
    
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var isLoading = false
    
    private let userDefaults = UserDefaults.standard
    private let userKey = "current_user"
    private let authKey = "is_authenticated"
    
    private init() {
        loadUserFromStorage()
    }
    
    // MARK: - Authentication Methods
    
    func login(email: String, password: String) -> AnyPublisher<Bool, Error> {
        isLoading = true
        
        // Mock authentication - in a real app, this would call your backend
        return Future<Bool, Error> { [weak self] promise in
            DispatchQueue.main.asyncAfter(deadline: .now() + 1.0) {
                // Simulate successful login for any email/password
                let user = User(
                    id: UUID().uuidString,
                    email: email,
                    name: email.components(separatedBy: "@").first?.capitalized ?? "User",
                    phone: nil,
                    avatar: nil,
                    provider: "local"
                )
                
                self?.currentUser = user
                self?.isAuthenticated = true
                self?.saveUserToStorage(user: user)
                self?.isLoading = false
                
                promise(.success(true))
            }
        }
        .eraseToAnyPublisher()
    }
    
    func loginWithGoogle() -> AnyPublisher<Bool, Error> {
        isLoading = true
        
        // Mock Google authentication
        return Future<Bool, Error> { [weak self] promise in
            DispatchQueue.main.asyncAfter(deadline: .now() + 1.5) {
                let user = User(
                    id: UUID().uuidString,
                    email: "user@gmail.com",
                    name: "Google User",
                    phone: nil,
                    avatar: nil,
                    provider: "google"
                )
                
                self?.currentUser = user
                self?.isAuthenticated = true
                self?.saveUserToStorage(user: user)
                self?.isLoading = false
                
                promise(.success(true))
            }
        }
        .eraseToAnyPublisher()
    }
    
    func logout() {
        currentUser = nil
        isAuthenticated = false
        clearUserFromStorage()
    }
    
    // MARK: - Admin Authentication
    
    func adminLogin(email: String, password: String) -> AnyPublisher<Bool, Error> {
        isLoading = true
        
        // Mock admin authentication
        return Future<Bool, Error> { [weak self] promise in
            DispatchQueue.main.asyncAfter(deadline: .now() + 1.0) {
                // Simple admin check - in real app, this would be more secure
                if email.lowercased().contains("admin") || password == "admin123" {
                    let user = User(
                        id: UUID().uuidString,
                        email: email,
                        name: "Admin User",
                        phone: nil,
                        avatar: nil,
                        provider: "admin"
                    )
                    
                    self?.currentUser = user
                    self?.isAuthenticated = true
                    self?.saveUserToStorage(user: user)
                    self?.isLoading = false
                    
                    promise(.success(true))
                } else {
                    self?.isLoading = false
                    promise(.failure(AuthError.invalidCredentials))
                }
            }
        }
        .eraseToAnyPublisher()
    }
    
    var isAdmin: Bool {
        return currentUser?.provider == "admin"
    }
    
    // MARK: - Storage Methods
    
    private func saveUserToStorage(user: User) {
        do {
            let data = try JSONEncoder().encode(user)
            userDefaults.set(data, forKey: userKey)
            userDefaults.set(true, forKey: authKey)
        } catch {
            print("Failed to save user: \(error)")
        }
    }
    
    private func loadUserFromStorage() {
        guard let data = userDefaults.data(forKey: userKey),
              let user = try? JSONDecoder().decode(User.self, from: data) else {
            return
        }
        
        currentUser = user
        isAuthenticated = userDefaults.bool(forKey: authKey)
    }
    
    private func clearUserFromStorage() {
        userDefaults.removeObject(forKey: userKey)
        userDefaults.set(false, forKey: authKey)
    }
}

// MARK: - Authentication Errors
enum AuthError: Error, LocalizedError {
    case invalidCredentials
    case networkError
    case unknownError
    
    var errorDescription: String? {
        switch self {
        case .invalidCredentials:
            return "Invalid email or password"
        case .networkError:
            return "Network connection error"
        case .unknownError:
            return "An unknown error occurred"
        }
    }
}
