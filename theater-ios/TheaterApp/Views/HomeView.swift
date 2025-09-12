import SwiftUI
import Combine

struct HomeView: View {
    @StateObject private var apiService = APIService.shared
    @StateObject private var authService = AuthService.shared
    @State private var shows: [Show] = []
    @State private var searchText = ""
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var selectedShow: Show?
    @State private var showingBooking = false
    
    private var filteredShows: [Show] {
        if searchText.isEmpty {
            return shows
        } else {
            return shows.filter { show in
                show.title.localizedCaseInsensitiveContains(searchText) ||
                show.details.localizedCaseInsensitiveContains(searchText) ||
                show.location.localizedCaseInsensitiveContains(searchText)
            }
        }
    }
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Header
                headerView
                
                // Search Bar
                searchBar
                
                // Content
                if isLoading {
                    loadingView
                } else if shows.isEmpty {
                    emptyStateView
                } else {
                    showsListView
                }
            }
            .navigationTitle("Theater Shows")
            .navigationBarTitleDisplayMode(.large)
            .onAppear {
                loadShows()
            }
            .sheet(isPresented: $showingBooking) {
                if let selectedShow = selectedShow {
                    BookingView(show: selectedShow)
                }
            }
        }
    }
    
    // MARK: - Header View
    private var headerView: some View {
        VStack(spacing: 16) {
            // Welcome Section
            VStack(spacing: 8) {
                Text("Welcome to")
                    .font(.title2)
                    .foregroundColor(.secondary)
                
                Text("Illuminating Windows")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .foregroundColor(.primary)
                
                Text("Discover amazing live performances")
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                
                // Debug indicator for mock data
                if apiService.isUsingMockData {
                    Text("📱 Using Sample Data")
                        .font(.caption)
                        .foregroundColor(.blue)
                        .padding(.horizontal, 8)
                        .padding(.vertical, 4)
                        .background(Color.blue.opacity(0.1))
                        .cornerRadius(8)
                }
            }
            .padding(.horizontal)
            
            // Quick Action Button
            Button(action: {
                showingBooking = true
            }) {
                HStack {
                    Image(systemName: "ticket.fill")
                    Text("Book Tickets Now")
                        .fontWeight(.semibold)
                }
                .foregroundColor(.white)
                .padding()
                .background(
                    LinearGradient(
                        colors: [.blue, .purple],
                        startPoint: .leading,
                        endPoint: .trailing
                    )
                )
                .cornerRadius(12)
            }
            .padding(.horizontal)
        }
        .padding(.vertical)
        .background(Color(.systemGroupedBackground))
    }
    
    // MARK: - Search Bar
    private var searchBar: some View {
        HStack {
            Image(systemName: "magnifyingglass")
                .foregroundColor(.secondary)
            
            TextField("Search shows, locations...", text: $searchText)
                .textFieldStyle(PlainTextFieldStyle())
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .padding(.horizontal)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Loading View
    private var loadingView: some View {
        VStack(spacing: 16) {
            ProgressView()
                .scaleEffect(1.2)
            Text("Loading shows...")
                .foregroundColor(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
    
    // MARK: - Empty State View
    private var emptyStateView: some View {
        VStack(spacing: 16) {
            Image(systemName: "theatermasks")
                .font(.system(size: 60))
                .foregroundColor(.secondary)
            
            Text("No Shows Available")
                .font(.title2)
                .fontWeight(.semibold)
            
            Text("Check back later for upcoming performances")
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
            
            Button("Refresh") {
                loadShows()
            }
            .buttonStyle(.borderedProminent)
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
    
    // MARK: - Shows List View
    private var showsListView: some View {
        ScrollView {
            LazyVStack(spacing: 16) {
                ForEach(filteredShows) { show in
                    ShowCardView(show: show) {
                        selectedShow = show
                        showingBooking = true
                    }
                }
            }
            .padding()
        }
    }
    
    // MARK: - Methods
    private func loadShows() {
        isLoading = true
        errorMessage = nil
        
        apiService.fetchShows()
            .sink(
                receiveCompletion: { completion in
                    isLoading = false
                    if case .failure(let error) = completion {
                        errorMessage = error.localizedDescription
                    }
                },
                receiveValue: { shows in
                    self.shows = shows
                }
            )
            .store(in: &cancellables)
    }
    
    @State private var cancellables = Set<AnyCancellable>()
}

// MARK: - Show Card View
struct ShowCardView: View {
    let show: Show
    let onBookTapped: () -> Void
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            // Show Image Placeholder
            ZStack {
                Rectangle()
                    .fill(
                        LinearGradient(
                            colors: [.blue, .purple],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(height: 200)
                
                Image(systemName: "theatermasks.fill")
                    .font(.system(size: 60))
                    .foregroundColor(.white)
            }
            .cornerRadius(12, corners: [.topLeft, .topRight])
            
            VStack(alignment: .leading, spacing: 8) {
                // Title and Price
                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text(show.title)
                            .font(.title2)
                            .fontWeight(.bold)
                            .lineLimit(2)
                        
                        Text("Theater Show")
                            .font(.caption)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 4)
                            .background(Color.blue.opacity(0.1))
                            .foregroundColor(.blue)
                            .cornerRadius(8)
                    }
                    
                    Spacer()
                    
                    VStack(alignment: .trailing) {
                        Text(String(format: "$%.2f", show.priceInDollars))
                            .font(.title2)
                            .fontWeight(.bold)
                            .foregroundColor(.blue)
                        Text("per ticket")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
                
                // Description
                Text(show.details)
                    .font(.body)
                    .foregroundColor(.secondary)
                    .lineLimit(3)
                
                // Show Details
                VStack(alignment: .leading, spacing: 4) {
                    DetailRow(icon: "location", text: show.location)
                    DetailRow(icon: "ticket", text: "\(show.availableTickets) of \(show.totalTickets) available")
                    DetailRow(icon: "calendar", text: formatDate(show.showDate))
                }
                .font(.caption)
                .foregroundColor(.secondary)
                
                // Show Times
                if !show.showTimes.isEmpty {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Upcoming Shows:")
                            .font(.caption)
                            .fontWeight(.semibold)
                            .foregroundColor(.primary)
                        
                        ForEach(show.showTimes.prefix(3)) { showTime in
                            HStack {
                                Text(formatDate(showTime.date))
                                    .font(.caption)
                                Text("•")
                                    .font(.caption)
                                Text(formatTime(showTime.time))
                                    .font(.caption)
                                
                                Spacer()
                                
                                Text("\(showTime.availableSeats ?? show.availableTickets) seats")
                                    .font(.caption)
                                    .foregroundColor(availabilityColor(showTime.availableSeats ?? show.availableTickets, show.totalTickets))
                            }
                            .padding(.vertical, 2)
                        }
                        
                        if show.showTimes.count > 3 {
                            Text("+\(show.showTimes.count - 3) more shows")
                                .font(.caption)
                                .foregroundColor(.secondary)
                        }
                    }
                    .padding(.top, 4)
                }
                
                // Book Button
                Button(action: onBookTapped) {
                    HStack {
                        Image(systemName: "ticket.fill")
                        Text("Book Tickets")
                            .fontWeight(.semibold)
                    }
                    .foregroundColor(.white)
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(Color.blue)
                    .cornerRadius(12)
                }
                .padding(.top, 8)
            }
            .padding()
        }
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 4, x: 0, y: 2)
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
    
    private func availabilityColor(_ available: Int, _ total: Int) -> Color {
        let percentage = Double(available) / Double(total)
        if percentage > 0.5 { return .green }
        if percentage > 0.25 { return .orange }
        return .red
    }
}

// MARK: - Detail Row
struct DetailRow: View {
    let icon: String
    let text: String
    
    var body: some View {
        HStack(spacing: 8) {
            Image(systemName: icon)
                .foregroundColor(.blue)
                .frame(width: 16)
            Text(text)
        }
    }
}

// MARK: - Corner Radius Extension
extension View {
    func cornerRadius(_ radius: CGFloat, corners: UIRectCorner) -> some View {
        clipShape(RoundedCorner(radius: radius, corners: corners))
    }
}

struct RoundedCorner: Shape {
    var radius: CGFloat = .infinity
    var corners: UIRectCorner = .allCorners

    func path(in rect: CGRect) -> Path {
        let path = UIBezierPath(
            roundedRect: rect,
            byRoundingCorners: corners,
            cornerRadii: CGSize(width: radius, height: radius)
        )
        return Path(path.cgPath)
    }
}

#Preview {
    HomeView()
}
