import SwiftUI
import Combine

struct BookingView: View {
    let show: Show
    @StateObject private var authService = AuthService.shared
    @StateObject private var apiService = APIService.shared
    @Environment(\.dismiss) private var dismiss
    
    @State private var selectedShowTime: ShowTime?
    @State private var numberOfTickets = 1
    @State private var customerName = ""
    @State private var customerPhone = ""
    @State private var customerEmail = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var bookingResult: Booking?
    @State private var showingSuccess = false
    
    private var totalPrice: Double {
        guard let selectedShowTime = selectedShowTime else { return 0 }
        return Double(show.price) / 100.0 * Double(numberOfTickets)
    }
    
    private var isFormValid: Bool {
        selectedShowTime != nil &&
        numberOfTickets > 0 &&
        !customerName.isEmpty &&
        !customerPhone.isEmpty
    }
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Show Header
                    showHeaderView
                    
                    if !authService.isAuthenticated {
                        // Login Section
                        loginSection
                    } else {
                        // Booking Form
                        bookingFormSection
                        
                        // Booking Summary
                        bookingSummarySection
                    }
                }
                .padding()
            }
            .navigationTitle("Book Tickets")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
            .alert("Booking Confirmed!", isPresented: $showingSuccess) {
                Button("OK") {
                    dismiss()
                }
            } message: {
                if let booking = bookingResult {
                    Text("Your booking ID is \(booking.bookingId). You will receive a confirmation email shortly.")
                }
            }
        }
    }
    
    // MARK: - Show Header View
    private var showHeaderView: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(show.title)
                .font(.title)
                .fontWeight(.bold)
            
            Text(show.details)
                .font(.body)
                .foregroundColor(.secondary)
            
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("Location")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    Text(show.location)
                        .font(.subheadline)
                        .fontWeight(.medium)
                }
                
                Spacer()
                
                VStack(alignment: .trailing, spacing: 4) {
                    Text("Price")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    Text(String(format: "$%.2f", show.priceInDollars))
                        .font(.title2)
                        .fontWeight(.bold)
                        .foregroundColor(.blue)
                }
            }
        }
        .padding()
        .background(Color(.systemGroupedBackground))
        .cornerRadius(12)
    }
    
    // MARK: - Login Section
    private var loginSection: some View {
        VStack(spacing: 16) {
            Text("Login to Book Tickets")
                .font(.title2)
                .fontWeight(.semibold)
            
            VStack(spacing: 12) {
                TextField("Email", text: $customerEmail)
                    .textFieldStyle(RoundedBorderTextFieldStyle())
                    .keyboardType(.emailAddress)
                    .autocapitalization(.none)
                
                SecureField("Password", text: .constant(""))
                    .textFieldStyle(RoundedBorderTextFieldStyle())
            }
            
            Button(action: login) {
                if authService.isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                } else {
                    Text("Login")
                        .fontWeight(.semibold)
                }
            }
            .disabled(authService.isLoading)
            .frame(maxWidth: .infinity)
            .padding()
            .background(Color.blue)
            .foregroundColor(.white)
            .cornerRadius(12)
            
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
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 4, x: 0, y: 2)
    }
    
    // MARK: - Booking Form Section
    private var bookingFormSection: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text("Booking Details")
                .font(.title2)
                .fontWeight(.semibold)
            
            // Show Time Selection
            VStack(alignment: .leading, spacing: 8) {
                Text("Select Show Time")
                    .font(.headline)
                
                if show.showTimes.isEmpty {
                    Text("No show times available")
                        .foregroundColor(.secondary)
                        .padding()
                        .frame(maxWidth: .infinity)
                        .background(Color(.systemGroupedBackground))
                        .cornerRadius(8)
                } else {
                    ForEach(show.showTimes) { showTime in
                        ShowTimeRow(
                            showTime: showTime,
                            isSelected: selectedShowTime?.id == showTime.id,
                            onTap: { selectedShowTime = showTime }
                        )
                    }
                }
            }
            
            // Number of Tickets
            VStack(alignment: .leading, spacing: 8) {
                Text("Number of Tickets")
                    .font(.headline)
                
                HStack {
                    Button(action: { if numberOfTickets > 1 { numberOfTickets -= 1 } }) {
                        Image(systemName: "minus.circle.fill")
                            .font(.title2)
                            .foregroundColor(.blue)
                    }
                    .disabled(numberOfTickets <= 1)
                    
                    Text("\(numberOfTickets)")
                        .font(.title2)
                        .fontWeight(.semibold)
                        .frame(minWidth: 50)
                    
                    Button(action: { 
                        if let selectedShowTime = selectedShowTime {
                            let maxTickets = selectedShowTime.availableSeats ?? show.availableTickets
                            if numberOfTickets < maxTickets { numberOfTickets += 1 }
                        } else {
                            numberOfTickets += 1
                        }
                    }) {
                        Image(systemName: "plus.circle.fill")
                            .font(.title2)
                            .foregroundColor(.blue)
                    }
                    .disabled(selectedShowTime != nil && numberOfTickets >= (selectedShowTime?.availableSeats ?? show.availableTickets))
                }
            }
            
            // Customer Information
            VStack(alignment: .leading, spacing: 12) {
                Text("Your Information")
                    .font(.headline)
                
                VStack(spacing: 12) {
                    TextField("Full Name", text: $customerName)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                    
                    TextField("Phone Number", text: $customerPhone)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                        .keyboardType(.phonePad)
                    
                    TextField("Email", text: $customerEmail)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                        .keyboardType(.emailAddress)
                        .autocapitalization(.none)
                }
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
            
            // Book Button
            Button(action: createBooking) {
                if isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                } else {
                    Text("Confirm Booking")
                        .fontWeight(.semibold)
                }
            }
            .disabled(!isFormValid || isLoading)
            .frame(maxWidth: .infinity)
            .padding()
            .background(isFormValid ? Color.blue : Color.gray)
            .foregroundColor(.white)
            .cornerRadius(12)
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 4, x: 0, y: 2)
    }
    
    // MARK: - Booking Summary Section
    private var bookingSummarySection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Booking Summary")
                .font(.title2)
                .fontWeight(.semibold)
            
            VStack(spacing: 8) {
                SummaryRow(label: "User", value: authService.currentUser?.email ?? "Guest")
                SummaryRow(label: "Show", value: show.title)
                
                if let selectedShowTime = selectedShowTime {
                    SummaryRow(label: "Date", value: formatDate(selectedShowTime.date))
                    SummaryRow(label: "Time", value: formatTime(selectedShowTime.time))
                }
                
                SummaryRow(label: "Tickets", value: "\(numberOfTickets)")
                SummaryRow(label: "Total Price", value: String(format: "$%.2f", totalPrice))
            }
        }
        .padding()
        .background(Color(.systemGroupedBackground))
        .cornerRadius(12)
    }
    
    // MARK: - Methods
    private func login() {
        authService.login(email: customerEmail, password: "password")
            .sink(
                receiveCompletion: { completion in
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { _ in
                    // Login successful
                }
            )
            .store(in: &cancellables)
    }
    
    private func loginWithGoogle() {
        authService.loginWithGoogle()
            .sink(
                receiveCompletion: { completion in
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { _ in
                    // Login successful
                }
            )
            .store(in: &cancellables)
    }
    
    private func createBooking() {
        guard let selectedShowTime = selectedShowTime else { return }
        
        isLoading = true
        errorMessage = nil
        
        let bookingRequest = BookingRequest(
            showId: show.id,
            contactType: "mobile",
            contactValue: customerPhone,
            numberOfTickets: numberOfTickets,
            customerName: customerName,
            bookingDate: selectedShowTime.date
        )
        
        apiService.createBooking(bookingRequest)
            .sink(
                receiveCompletion: { completion in
                    isLoading = false
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { booking in
                    bookingResult = booking
                    showingSuccess = true
                }
            )
            .store(in: &cancellables)
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
    
    private func formatTime(_ timeString: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        if let time = formatter.date(from: timeString) {
            formatter.dateFormat = "h:mm a"
            return formatter.string(from: time)
        }
        return timeString
    }
    
    @State private var cancellables = Set<AnyCancellable>()
}

// MARK: - Show Time Row
struct ShowTimeRow: View {
    let showTime: ShowTime
    let isSelected: Bool
    let onTap: () -> Void
    
    var body: some View {
        Button(action: onTap) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text(formatDate(showTime.date))
                        .font(.headline)
                        .foregroundColor(.primary)
                    
                    Text(formatTime(showTime.time))
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                
                Spacer()
                
                VStack(alignment: .trailing, spacing: 4) {
                    Text("\(showTime.availableSeats ?? 0) seats")
                        .font(.caption)
                        .foregroundColor(availabilityColor)
                    
                    if isSelected {
                        Image(systemName: "checkmark.circle.fill")
                            .foregroundColor(.blue)
                    }
                }
            }
            .padding()
            .background(isSelected ? Color.blue.opacity(0.1) : Color(.systemGroupedBackground))
            .cornerRadius(8)
            .overlay(
                RoundedRectangle(cornerRadius: 8)
                    .stroke(isSelected ? Color.blue : Color.clear, lineWidth: 2)
            )
        }
        .buttonStyle(PlainButtonStyle())
    }
    
    private var availabilityColor: Color {
        guard let available = showTime.availableSeats else { return .secondary }
        if available > 10 { return .green }
        if available > 5 { return .orange }
        return .red
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
    
    private func formatTime(_ timeString: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm"
        if let time = formatter.date(from: timeString) {
            formatter.dateFormat = "h:mm a"
            return formatter.string(from: time)
        }
        return timeString
    }
}

// MARK: - Summary Row
struct SummaryRow: View {
    let label: String
    let value: String
    
    var body: some View {
        HStack {
            Text(label)
                .font(.subheadline)
                .foregroundColor(.secondary)
            Spacer()
            Text(value)
                .font(.subheadline)
                .fontWeight(.medium)
        }
    }
}

#Preview {
    let sampleShow = Show(
        id: "1",
        name: "Sample Show",
        title: "Romeo & Juliet",
        details: "A timeless tale of love and tragedy",
        price: 5000,
        totalTickets: 100,
        bookedTickets: 20,
        location: "Main Theater",
        showNumber: "RJ001",
        showDate: "2024-01-15",
        showTimes: [
            ShowTime(id: "1", date: "2024-01-15", time: "19:30", showId: "1", availableSeats: 80, totalSeats: 100),
            ShowTime(id: "2", date: "2024-01-16", time: "19:30", showId: "1", availableSeats: 75, totalSeats: 100)
        ],
        images: [],
        videos: [],
        createdAt: "2024-01-01T00:00:00Z",
        updatedAt: "2024-01-01T00:00:00Z"
    )
    
    return BookingView(show: sampleShow)
}
