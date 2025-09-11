import SwiftUI
import Combine

struct BookingsView: View {
    @StateObject private var authService = AuthService.shared
    @StateObject private var apiService = APIService.shared
    @State private var bookings: [Booking] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showingLogin = false
    
    var body: some View {
        NavigationView {
            VStack {
                if !authService.isAuthenticated {
                    // Login Prompt
                    loginPromptView
                } else if isLoading {
                    // Loading State
                    loadingView
                } else if bookings.isEmpty {
                    // Empty State
                    emptyStateView
                } else {
                    // Bookings List
                    bookingsListView
                }
            }
            .navigationTitle("My Bookings")
            .navigationBarTitleDisplayMode(.large)
            .onAppear {
                if authService.isAuthenticated {
                    loadBookings()
                }
            }
            .sheet(isPresented: $showingLogin) {
                LoginView()
            }
        }
    }
    
    // MARK: - Login Prompt View
    private var loginPromptView: some View {
        VStack(spacing: 20) {
            Image(systemName: "person.circle")
                .font(.system(size: 80))
                .foregroundColor(.secondary)
            
            Text("Sign In to View Bookings")
                .font(.title2)
                .fontWeight(.semibold)
            
            Text("Please sign in to view your ticket bookings and manage your reservations.")
                .font(.body)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
            
            Button("Sign In") {
                showingLogin = true
            }
            .buttonStyle(.borderedProminent)
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
    
    // MARK: - Loading View
    private var loadingView: some View {
        VStack(spacing: 16) {
            ProgressView()
                .scaleEffect(1.2)
            Text("Loading bookings...")
                .foregroundColor(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
    
    // MARK: - Empty State View
    private var emptyStateView: some View {
        VStack(spacing: 20) {
            Image(systemName: "ticket")
                .font(.system(size: 60))
                .foregroundColor(.secondary)
            
            Text("No Bookings Yet")
                .font(.title2)
                .fontWeight(.semibold)
            
            Text("You haven't made any bookings yet. Browse our shows and book your tickets!")
                .font(.body)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
            
            Button("Browse Shows") {
                // This would navigate to the home tab
            }
            .buttonStyle(.borderedProminent)
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
    
    // MARK: - Bookings List View
    private var bookingsListView: some View {
        ScrollView {
            LazyVStack(spacing: 16) {
                ForEach(bookings) { booking in
                    BookingCard(booking: booking)
                }
            }
            .padding()
        }
    }
    
    // MARK: - Methods
    private func loadBookings() {
        guard let user = authService.currentUser else { return }
        
        isLoading = true
        errorMessage = nil
        
        // Try to load bookings by email first, then by phone
        let emailPublisher = apiService.getBookingsByContact(contactValue: user.email)
        let phonePublisher = user.phone != nil ? 
            apiService.getBookingsByContact(contactValue: user.phone!) : 
            Just([]).setFailureType(to: Error.self).eraseToAnyPublisher()
        
        Publishers.Merge(emailPublisher, phonePublisher)
            .sink(
                receiveCompletion: { completion in
                    isLoading = false
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { bookings in
                    self.bookings = bookings
                }
            )
            .store(in: &cancellables)
    }
    
    @State private var cancellables = Set<AnyCancellable>()
}

// MARK: - Booking Card
struct BookingCard: View {
    let booking: Booking
    @State private var show: Show?
    @State private var isLoading = true
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            // Header
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text(show?.title ?? "Loading...")
                        .font(.headline)
                        .fontWeight(.semibold)
                    
                    Text("Booking #\(booking.bookingId.prefix(8))")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                
                Spacer()
                
                StatusBadge(status: booking.status)
            }
            
            // Booking Details
            VStack(spacing: 8) {
                DetailRow(icon: "calendar", text: formatDate(booking.bookingDate))
                DetailRow(icon: "ticket", text: "\(booking.numberOfTickets) ticket\(booking.numberOfTickets == 1 ? "" : "s")")
                DetailRow(icon: "person", text: booking.customerName ?? "Guest")
                DetailRow(icon: "dollarsign.circle", text: String(format: "$%.2f", booking.totalAmountInDollars))
            }
            
            // Actions
            HStack(spacing: 12) {
                if booking.status == "pending" {
                    Button("Confirm") {
                        // Handle confirmation
                    }
                    .buttonStyle(.bordered)
                    .foregroundColor(.green)
                    
                    Button("Cancel") {
                        // Handle cancellation
                    }
                    .buttonStyle(.bordered)
                    .foregroundColor(.red)
                }
                
                Spacer()
                
                Button("View Details") {
                    // Handle view details
                }
                .buttonStyle(.bordered)
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
        .onAppear {
            loadShowDetails()
        }
    }
    
    private func loadShowDetails() {
        // In a real app, you'd load the show details
        // For now, we'll just simulate it
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
            isLoading = false
            // show would be set here
        }
    }
    
    private func formatDate(_ dateString: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        if let date = formatter.date(from: dateString) {
            formatter.dateStyle = .medium
            return formatter.string(from: date)
        }
        return dateString
    }
}

// MARK: - Status Badge
struct StatusBadge: View {
    let status: String
    
    var body: some View {
        Text(status.capitalized)
            .font(.caption)
            .fontWeight(.semibold)
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
            .background(backgroundColor)
            .foregroundColor(textColor)
            .cornerRadius(8)
    }
    
    private var backgroundColor: Color {
        switch status.lowercased() {
        case "confirmed":
            return .green.opacity(0.2)
        case "pending":
            return .orange.opacity(0.2)
        case "cancelled":
            return .red.opacity(0.2)
        default:
            return .gray.opacity(0.2)
        }
    }
    
    private var textColor: Color {
        switch status.lowercased() {
        case "confirmed":
            return .green
        case "pending":
            return .orange
        case "cancelled":
            return .red
        default:
            return .gray
        }
    }
}

// MARK: - Login View
struct LoginView: View {
    @Environment(\.dismiss) private var dismiss
    @StateObject private var authService = AuthService.shared
    @State private var email = ""
    @State private var password = ""
    @State private var errorMessage: String?
    
    var body: some View {
        NavigationView {
            VStack(spacing: 24) {
                // Header
                VStack(spacing: 16) {
                    Image(systemName: "person.circle.fill")
                        .font(.system(size: 80))
                        .foregroundColor(.blue)
                    
                    Text("Sign In")
                        .font(.largeTitle)
                        .fontWeight(.bold)
                    
                    Text("Sign in to view your bookings")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                .padding(.top, 40)
                
                // Login Form
                VStack(spacing: 16) {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Email")
                            .font(.headline)
                            .foregroundColor(.primary)
                        
                        TextField("Enter your email", text: $email)
                            .textFieldStyle(RoundedBorderTextFieldStyle())
                            .keyboardType(.emailAddress)
                            .autocapitalization(.none)
                    }
                    
                    VStack(alignment: .leading, spacing: 8) {
                        Text("Password")
                            .font(.headline)
                            .foregroundColor(.primary)
                        
                        SecureField("Enter your password", text: $password)
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
                    
                    // Google Sign In
                    Button(action: loginWithGoogle) {
                        HStack {
                            Image(systemName: "globe")
                            Text("Continue with Google")
                                .fontWeight(.semibold)
                        }
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.white)
                        .foregroundColor(.black)
                        .overlay(
                            RoundedRectangle(cornerRadius: 12)
                                .stroke(Color.gray, lineWidth: 1)
                        )
                        .cornerRadius(12)
                    }
                    .disabled(authService.isLoading)
                }
                .padding(.horizontal, 32)
                
                Spacer()
            }
            .navigationTitle("Sign In")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
            .onChange(of: authService.isAuthenticated) { isAuthenticated in
                if isAuthenticated {
                    dismiss()
                }
            }
        }
    }
    
    private func login() {
        errorMessage = nil
        
        authService.login(email: email, password: password)
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
    
    private func loginWithGoogle() {
        errorMessage = nil
        
        authService.loginWithGoogle()
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
    BookingsView()
}
