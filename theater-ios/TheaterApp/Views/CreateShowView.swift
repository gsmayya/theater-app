import SwiftUI
import Combine

struct CreateShowView: View {
    @Environment(\.dismiss) private var dismiss
    @StateObject private var apiService = APIService.shared
    
    @State private var name = ""
    @State private var details = ""
    @State private var price = ""
    @State private var totalTickets = ""
    @State private var location = ""
    @State private var showNumber = ""
    @State private var showDate = Date()
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showingSuccess = false
    
    let onShowCreated: (Show) -> Void
    
    private var isFormValid: Bool {
        !name.isEmpty &&
        !details.isEmpty &&
        !price.isEmpty &&
        !totalTickets.isEmpty &&
        !location.isEmpty &&
        !showNumber.isEmpty &&
        Double(price) != nil &&
        Int(totalTickets) != nil
    }
    
    var body: some View {
        NavigationView {
            ScrollView {
                VStack(spacing: 20) {
                    // Header
                    VStack(spacing: 8) {
                        Image(systemName: "plus.circle.fill")
                            .font(.system(size: 60))
                            .foregroundColor(.blue)
                        
                        Text("Create New Show")
                            .font(.title)
                            .fontWeight(.bold)
                        
                        Text("Add a new theater show to the system")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                            .multilineTextAlignment(.center)
                    }
                    .padding(.top, 20)
                    
                    // Form
                    VStack(spacing: 16) {
                        // Basic Information
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Basic Information")
                                .font(.headline)
                                .fontWeight(.semibold)
                            
                            VStack(alignment: .leading, spacing: 8) {
                                Text("Show Name")
                                    .font(.subheadline)
                                    .fontWeight(.medium)
                                
                                TextField("Enter show name", text: $name)
                                    .textFieldStyle(RoundedBorderTextFieldStyle())
                            }
                            
                            VStack(alignment: .leading, spacing: 8) {
                                Text("Description")
                                    .font(.subheadline)
                                    .fontWeight(.medium)
                                
                                TextField("Enter show description", text: $details, axis: .vertical)
                                    .textFieldStyle(RoundedBorderTextFieldStyle())
                                    .lineLimit(3...6)
                            }
                            
                            VStack(alignment: .leading, spacing: 8) {
                                Text("Location")
                                    .font(.subheadline)
                                    .fontWeight(.medium)
                                
                                TextField("Enter venue location", text: $location)
                                    .textFieldStyle(RoundedBorderTextFieldStyle())
                            }
                            
                            VStack(alignment: .leading, spacing: 8) {
                                Text("Show Number")
                                    .font(.subheadline)
                                    .fontWeight(.medium)
                                
                                TextField("Enter show number (e.g., RJ001)", text: $showNumber)
                                    .textFieldStyle(RoundedBorderTextFieldStyle())
                            }
                        }
                        
                        // Pricing and Tickets
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Pricing & Tickets")
                                .font(.headline)
                                .fontWeight(.semibold)
                            
                            HStack(spacing: 16) {
                                VStack(alignment: .leading, spacing: 8) {
                                    Text("Price per Ticket")
                                        .font(.subheadline)
                                        .fontWeight(.medium)
                                    
                                    HStack {
                                        Text("$")
                                            .foregroundColor(.secondary)
                                        TextField("0.00", text: $price)
                                            .textFieldStyle(RoundedBorderTextFieldStyle())
                                            .keyboardType(.decimalPad)
                                    }
                                }
                                
                                VStack(alignment: .leading, spacing: 8) {
                                    Text("Total Tickets")
                                        .font(.subheadline)
                                        .fontWeight(.medium)
                                    
                                    TextField("100", text: $totalTickets)
                                        .textFieldStyle(RoundedBorderTextFieldStyle())
                                        .keyboardType(.numberPad)
                                }
                            }
                        }
                        
                        // Show Date
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Show Date")
                                .font(.headline)
                                .fontWeight(.semibold)
                            
                            DatePicker(
                                "Select show date",
                                selection: $showDate,
                                in: Date()...,
                                displayedComponents: .date
                            )
                            .datePickerStyle(CompactDatePickerStyle())
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
                        
                        // Create Button
                        Button(action: createShow) {
                            if isLoading {
                                ProgressView()
                                    .progressViewStyle(CircularProgressViewStyle(tint: .white))
                            } else {
                                Text("Create Show")
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
                    .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
                }
                .padding()
            }
            .navigationTitle("Create Show")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        dismiss()
                    }
                }
            }
            .alert("Show Created!", isPresented: $showingSuccess) {
                Button("OK") {
                    dismiss()
                }
            } message: {
                Text("The show has been successfully created and added to the system.")
            }
        }
    }
    
    private func createShow() {
        guard let priceValue = Double(price),
              let totalTicketsValue = Int(totalTickets) else {
            errorMessage = "Please enter valid price and ticket numbers"
            return
        }
        
        isLoading = true
        errorMessage = nil
        
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "yyyy-MM-dd"
        let showDateString = dateFormatter.string(from: showDate)
        
        let createRequest = CreateShowRequest(
            name: name,
            details: details,
            price: Int(priceValue * 100), // Convert to cents
            totalTickets: totalTicketsValue,
            location: location,
            showNumber: showNumber,
            showDate: showDateString
        )
        
        apiService.createShow(createRequest)
            .sink(
                receiveCompletion: { completion in
                    isLoading = false
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { show in
                    onShowCreated(show)
                    showingSuccess = true
                }
            )
            .store(in: &cancellables)
    }
    
    @State private var cancellables = Set<AnyCancellable>()
}

#Preview {
    CreateShowView { _ in }
}
