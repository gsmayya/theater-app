import SwiftUI
import Combine

struct ProfileView: View {
    @StateObject private var authService = AuthService.shared
    @State private var showingLogin = false
    @State private var showingAdminLogin = false
    @State private var showingAdminDashboard = false
    @State private var showingDebugSettings = false
    
    var body: some View {
        NavigationView {
            VStack {
                if !authService.isAuthenticated {
                    // Not logged in state
                    notLoggedInView
                } else {
                    // Logged in state
                    loggedInView
                }
            }
            .navigationTitle("Profile")
            .navigationBarTitleDisplayMode(.large)
            .sheet(isPresented: $showingLogin) {
                LoginView()
            }
            .sheet(isPresented: $showingAdminLogin) {
                AdminLoginView()
            }
            .sheet(isPresented: $showingAdminDashboard) {
                AdminDashboardView()
            }
            .sheet(isPresented: $showingDebugSettings) {
                DebugSettingsView()
            }
        }
    }
    
    // MARK: - Not Logged In View
    private var notLoggedInView: some View {
        VStack(spacing: 24) {
            // Profile Icon
            Image(systemName: "person.circle")
                .font(.system(size: 100))
                .foregroundColor(.secondary)
            
            // Welcome Text
            VStack(spacing: 8) {
                Text("Welcome to Theater App")
                    .font(.title2)
                    .fontWeight(.bold)
                
                Text("Sign in to access your bookings and manage your account")
                    .font(.body)
                    .foregroundColor(.secondary)
                    .multilineTextAlignment(.center)
            }
            
            // Action Buttons
            VStack(spacing: 12) {
                Button("Sign In") {
                    showingLogin = true
                }
                .buttonStyle(.borderedProminent)
                .controlSize(.large)
                
                Button("Admin Login") {
                    showingAdminLogin = true
                }
                .buttonStyle(.bordered)
                .controlSize(.large)
            }
            
            Spacer()
        }
        .padding()
    }
    
    // MARK: - Logged In View
    private var loggedInView: some View {
        ScrollView {
            VStack(spacing: 24) {
                // User Info Section
                userInfoSection
                
                // Quick Actions
                quickActionsSection
                
                // Settings Section
                settingsSection
                
                // Admin Section (if admin)
                if authService.isAdmin {
                    adminSection
                }
                
                // Sign Out Button
                signOutSection
            }
            .padding()
        }
    }
    
    // MARK: - User Info Section
    private var userInfoSection: some View {
        VStack(spacing: 16) {
            // Profile Picture
            Image(systemName: "person.circle.fill")
                .font(.system(size: 80))
                .foregroundColor(.blue)
            
            // User Details
            VStack(spacing: 4) {
                Text(authService.currentUser?.name ?? "User")
                    .font(.title2)
                    .fontWeight(.bold)
                
                Text(authService.currentUser?.email ?? "")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                
                if let phone = authService.currentUser?.phone {
                    Text(phone)
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                
                // Provider Badge
                if let provider = authService.currentUser?.provider {
                    Text(provider.capitalized)
                        .font(.caption)
                        .fontWeight(.semibold)
                        .padding(.horizontal, 8)
                        .padding(.vertical, 4)
                        .background(providerColor(provider))
                        .foregroundColor(.white)
                        .cornerRadius(8)
                }
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Quick Actions Section
    private var quickActionsSection: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Quick Actions")
                .font(.headline)
                .fontWeight(.semibold)
            
            VStack(spacing: 12) {
                ActionRow(
                    icon: "ticket.fill",
                    title: "My Bookings",
                    subtitle: "View your ticket bookings",
                    action: {
                        // Navigate to bookings
                    }
                )
                
                ActionRow(
                    icon: "calendar",
                    title: "Calendar",
                    subtitle: "View shows by date",
                    action: {
                        // Navigate to calendar
                    }
                )
                
                ActionRow(
                    icon: "house.fill",
                    title: "Browse Shows",
                    subtitle: "Discover new performances",
                    action: {
                        // Navigate to home
                    }
                )
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Settings Section
    private var settingsSection: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Settings")
                .font(.headline)
                .fontWeight(.semibold)
            
            VStack(spacing: 12) {
                ActionRow(
                    icon: "bell.fill",
                    title: "Notifications",
                    subtitle: "Manage your notification preferences",
                    action: {
                        // Open notification settings
                    }
                )
                
                ActionRow(
                    icon: "gear",
                    title: "Preferences",
                    subtitle: "Customize your app experience",
                    action: {
                        // Open preferences
                    }
                )
                
                ActionRow(
                    icon: "questionmark.circle.fill",
                    title: "Help & Support",
                    subtitle: "Get help or contact support",
                    action: {
                        // Open help
                    }
                )
                
                ActionRow(
                    icon: "gear",
                    title: "Debug Settings",
                    subtitle: "Configure app behavior",
                    action: {
                        showingDebugSettings = true
                    }
                )
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Admin Section
    private var adminSection: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Admin")
                .font(.headline)
                .fontWeight(.semibold)
            
            VStack(spacing: 12) {
                ActionRow(
                    icon: "person.crop.circle.badge.checkmark",
                    title: "Admin Dashboard",
                    subtitle: "Manage shows and bookings",
                    action: {
                        showingAdminDashboard = true
                    }
                )
                
                ActionRow(
                    icon: "plus.circle.fill",
                    title: "Create Show",
                    subtitle: "Add new theater performances",
                    action: {
                        showingAdminDashboard = true
                    }
                )
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Sign Out Section
    private var signOutSection: some View {
        Button(action: {
            authService.logout()
        }) {
            HStack {
                Image(systemName: "arrow.right.square.fill")
                Text("Sign Out")
                    .fontWeight(.semibold)
            }
            .foregroundColor(.red)
            .frame(maxWidth: .infinity)
            .padding()
            .background(Color.red.opacity(0.1))
            .cornerRadius(12)
        }
    }
    
    // MARK: - Helper Methods
    private func providerColor(_ provider: String) -> Color {
        switch provider.lowercased() {
        case "google":
            return .blue
        case "admin":
            return .purple
        case "local":
            return .green
        default:
            return .gray
        }
    }
}

// MARK: - Action Row
struct ActionRow: View {
    let icon: String
    let title: String
    let subtitle: String
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 12) {
                Image(systemName: icon)
                    .font(.title2)
                    .foregroundColor(.blue)
                    .frame(width: 30)
                
                VStack(alignment: .leading, spacing: 2) {
                    Text(title)
                        .font(.subheadline)
                        .fontWeight(.medium)
                        .foregroundColor(.primary)
                    
                    Text(subtitle)
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                
                Spacer()
                
                Image(systemName: "chevron.right")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            .padding(.vertical, 8)
        }
        .buttonStyle(PlainButtonStyle())
    }
}

#Preview {
    ProfileView()
}
