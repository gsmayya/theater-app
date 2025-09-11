import SwiftUI
import Combine

struct AdminLoginView: View {
    @StateObject private var authService = AuthService.shared
    @State private var email = ""
    @State private var password = ""
    @State private var errorMessage: String?
    @State private var showingDashboard = false
    
    var body: some View {
        NavigationView {
            VStack(spacing: 24) {
                // Header
                VStack(spacing: 16) {
                    Image(systemName: "person.crop.circle.badge.checkmark")
                        .font(.system(size: 80))
                        .foregroundColor(.blue)
                    
                    Text("Admin Login")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                    
                    Text("Sign in to manage theater shows")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                        .multilineTextAlignment(.center)
                }
                .padding(.top, 40)
                
                // Login Form
                VStack(spacing: 16) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Email")
                            .font(.headline)
                            .foregroundColor(.primary)
                        
                        TextField("admin@theater.com", text: $email)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .keyboardType(.emailAddress)
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                    }
                    
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Password")
                            .font(.headline)
                            .foregroundColor(.primary)
                        
                        SecureField("Enter password", text: $password)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                    }
                    
                    // Error Message
                    if let errorMessage = errorMessage {
                        Text(errorMessage)
                            .foregroundColor(.red)
                            .font(.caption)
                            .padding()
                            .background(Color.red.opacity(0.1))
                            .cornerRadius(8)
                    }
                    
                    // Login Button
                    Button(action: login) {
                        if authService.isLoading {
                            ProgressView()
                                .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        } else {
                            Text("Sign In")
                                .fontWeight(.semibold)
                        }
                    }
                    .disabled(authService.isLoading || email.isEmpty || password.isEmpty)
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(email.isEmpty || password.isEmpty ? Color.gray : Color.blue)
                    .foregroundColor(.white)
                    .cornerRadius(12)
                }
                .padding(.horizontal, 32)
                
                // Demo Credentials
                VStack(spacing: 8) {
                    Text("Demo Credentials")
                        .font(.headline)
                        .foregroundColor(.secondary)
                    
                    VStack(spacing: 4) {
                        Text("Email: admin@theater.com")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        Text("Password: admin123")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding()
                    .background(Color(.systemGroupedBackground))
                    .cornerRadius(8)
                }
                .padding(.horizontal, 32)
                
                Spacer()
            }
            .navigationTitle("Admin")
            .navigationBarTitleDisplayMode(.inline)
            .onChange(of: authService.isAuthenticated) { isAuthenticated in
                if isAuthenticated && authService.isAdmin {
                    showingDashboard = true
                }
            }
            .sheet(isPresented: $showingDashboard) {
                AdminDashboardView()
            }
        }
    }
    
    private func login() {
        errorMessage = nil
        
        authService.adminLogin(email: email, password: password)
            .sink(
                receiveCompletion: { completion in
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { _ in
                    // Login successful - handled in onChange
                }
            )
            .store(in: &cancellables)
    }
    
    @State private var cancellables = Set<AnyCancellable>()
}

#Preview {
    AdminLoginView()
}
