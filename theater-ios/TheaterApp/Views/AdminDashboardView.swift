import SwiftUI
import Combine

struct AdminDashboardView: View {
    @StateObject private var authService = AuthService.shared
    @StateObject private var apiService = APIService.shared
    @Environment(\.dismiss) private var dismiss
    
    @State private var shows: [Show] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showingCreateShow = false
    @State private var selectedTab = 0
    
    private var totalShows: Int {
        shows.count
    }
    
    private var totalAvailableTickets: Int {
        shows.reduce(0) { $0 + $1.availableTickets }
    }
    
    private var averagePrice: Double {
        guard !shows.isEmpty else { return 0 }
        let totalPrice = shows.reduce(0) { $0 + $1.priceInDollars }
        return totalPrice / Double(shows.count)
    }
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Header
                headerView
                
                // Tab Picker
                tabPicker
                
                // Content
                TabView(selection: $selectedTab) {
                    // Shows List Tab
                    showsListView
                        .tag(0)
                    
                    // Statistics Tab
                    statisticsView
                        .tag(1)
                }
                .tabViewStyle(PageTabViewStyle(indexDisplayMode: .never))
            }
            .navigationTitle("Admin Dashboard")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Logout") {
                        authService.logout()
                        dismiss()
                    }
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Create Show") {
                        showingCreateShow = true
                    }
                    .disabled(isLoading)
                }
            }
            .onAppear {
                loadShows()
            }
            .sheet(isPresented: $showingCreateShow) {
                CreateShowView { newShow in
                    shows.append(newShow)
                }
            }
        }
    }
    
    // MARK: - Header View
    private var headerView: some View {
        VStack(spacing: 12) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("Welcome back!")
                        .font(.title2)
                        .fontWeight(.bold)
                    
                    Text(authService.currentUser?.name ?? "Admin")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                
                Spacer()
                
                // Refresh Button
                Button(action: loadShows) {
                    Image(systemName: "arrow.clockwise")
                        .font(.title2)
                        .foregroundColor(.blue)
                }
                .disabled(isLoading)
            }
            
            // Quick Stats
            HStack(spacing: 16) {
                StatCard(title: "Total Shows", value: "\(totalShows)", color: .blue)
                StatCard(title: "Available Tickets", value: "\(totalAvailableTickets)", color: .green)
                StatCard(title: "Avg Price", value: String(format: "$%.0f", averagePrice), color: .orange)
            }
        }
        .padding()
        .background(Color(.systemGroupedBackground))
    }
    
    // MARK: - Tab Picker
    private var tabPicker: some View {
        Picker("View", selection: $selectedTab) {
            Text("Shows").tag(0)
            Text("Statistics").tag(1)
        }
        .pickerStyle(SegmentedPickerStyle())
        .padding(.horizontal)
    }
    
    // MARK: - Shows List View
    private var showsListView: some View {
        VStack(spacing: 0) {
            if isLoading {
                loadingView
            } else if shows.isEmpty {
                emptyStateView
            } else {
                ScrollView {
                    LazyVStack(spacing: 12) {
                        ForEach(shows) { show in
                            AdminShowCard(show: show)
                        }
                    }
                    .padding()
                }
            }
        }
    }
    
    // MARK: - Statistics View
    private var statisticsView: some View {
        ScrollView {
            VStack(spacing: 20) {
                // Overview Cards
                VStack(spacing: 16) {
                    Text("Overview")
                        .font(.title2)
                        .fontWeight(.bold)
                        .frame(maxWidth: .infinity, alignment: .leading)
                    
                    LazyVGrid(columns: [
                        GridItem(.flexible()),
                        GridItem(.flexible())
                    ], spacing: 16) {
                        StatCard(
                            title: "Total Shows",
                            value: "\(totalShows)",
                            color: .blue,
                            icon: "theatermasks"
                        )
                        
                        StatCard(
                            title: "Available Tickets",
                            value: "\(totalAvailableTickets)",
                            color: .green,
                            icon: "ticket"
                        )
                        
                        StatCard(
                            title: "Average Price",
                            value: String(format: "$%.2f", averagePrice),
                            color: .orange,
                            icon: "dollarsign.circle"
                        )
                        
                        StatCard(
                            title: "Total Revenue",
                            value: String(format: "$%.0f", calculateTotalRevenue()),
                            color: .purple,
                            icon: "chart.bar"
                        )
                    }
                }
                
                // Shows by Location
                if !shows.isEmpty {
                    locationStatsView
                }
                
                // Recent Shows
                if !shows.isEmpty {
                    recentShowsView
                }
            }
            .padding()
        }
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
            
            Text("No Shows Created")
                .font(.title2)
                .fontWeight(.semibold)
            
            Text("Create your first show to get started")
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
            
            Button("Create Show") {
                showingCreateShow = true
            }
            .buttonStyle(.borderedProminent)
        }
        .padding()
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
    
    // MARK: - Location Stats View
    private var locationStatsView: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Shows by Location")
                .font(.headline)
                .fontWeight(.semibold)
            
            let locationGroups = Dictionary(grouping: shows) { $0.location }
            
            ForEach(Array(locationGroups.keys.sorted()), id: \.self) { location in
                let showsAtLocation = locationGroups[location] ?? []
                HStack {
                    Text(location)
                        .font(.subheadline)
                        .foregroundColor(.primary)
                    
                    Spacer()
                    
                    Text("\(showsAtLocation.count) show\(showsAtLocation.count == 1 ? "" : "s")")
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                .padding(.vertical, 4)
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Recent Shows View
    private var recentShowsView: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Recent Shows")
                .font(.headline)
                .fontWeight(.semibold)
            
            ForEach(shows.prefix(5)) { show in
                HStack {
                    VStack(alignment: .leading, spacing: 4) {
                        Text(show.title)
                            .font(.subheadline)
                            .fontWeight(.medium)
                        
                        Text(show.location)
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    
                    Spacer()
                    
                    VStack(alignment: .trailing, spacing: 4) {
                        Text(String(format: "$%.2f", show.priceInDollars))
                            .font(.subheadline)
                            .fontWeight(.semibold)
                            .foregroundColor(.blue)
                        
                        Text("\(show.availableTickets) available")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
                .padding(.vertical, 4)
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
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
    
    private func calculateTotalRevenue() -> Double {
        // This is a simplified calculation - in a real app, you'd calculate from actual bookings
        return shows.reduce(0) { $0 + ($1.priceInDollars * Double($1.totalTickets - $1.availableTickets)) }
    }
    
    @State private var cancellables = Set<AnyCancellable>()
}

// MARK: - Stat Card
struct StatCard: View {
    let title: String
    let value: String
    let color: Color
    let icon: String?
    
    init(title: String, value: String, color: Color, icon: String? = nil) {
        self.title = title
        self.value = value
        self.color = color
        self.icon = icon
    }
    
    var body: some View {
        VStack(spacing: 8) {
            if let icon = icon {
                Image(systemName: icon)
                    .font(.title2)
                    .foregroundColor(color)
            }
            
            Text(value)
                .font(.title2)
                .fontWeight(.bold)
                .foregroundColor(.primary)
            
            Text(title)
                .font(.caption)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
        }
        .padding()
        .frame(maxWidth: .infinity)
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
}

// MARK: - Admin Show Card
struct AdminShowCard: View {
    let show: Show
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text(show.title)
                        .font(.headline)
                        .fontWeight(.semibold)
                    
                    Text(show.location)
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                }
                
                Spacer()
                
                VStack(alignment: .trailing, spacing: 4) {
                    Text(String(format: "$%.2f", show.priceInDollars))
                        .font(.headline)
                        .fontWeight(.bold)
                        .foregroundColor(.blue)
                    
                    Text("per ticket")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
            }
            
            Text(show.details)
                .font(.body)
                .foregroundColor(.secondary)
                .lineLimit(2)
            
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("Show Number")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    Text(show.showNumber)
                        .font(.subheadline)
                        .fontWeight(.medium)
                }
                
                Spacer()
                
                VStack(alignment: .trailing, spacing: 4) {
                    Text("Available")
                        .font(.caption)
                        .foregroundColor(.secondary)
                    Text("\(show.availableTickets) of \(show.totalTickets)")
                        .font(.subheadline)
                        .fontWeight(.medium)
                        .foregroundColor(availabilityColor)
                }
            }
            
            // Show Times
            if !show.showTimes.isEmpty {
                VStack(alignment: .leading, spacing: 4) {
                    Text("Show Times")
                        .font(.caption)
                        .fontWeight(.semibold)
                        .foregroundColor(.secondary)
                    
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
                                .foregroundColor(.secondary)
                        }
                    }
                    
                    if show.showTimes.count > 3 {
                        Text("+\(show.showTimes.count - 3) more times")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                }
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    private var availabilityColor: Color {
        let percentage = Double(show.availableTickets) / Double(show.totalTickets)
        if percentage > 0.5 { return .green }
        if percentage > 0.25 { return .orange }
        return .red
    }
    
    private func formatDate(_ dateString: String) -> String {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        if let date = formatter.date(from: dateString) {
            formatter.dateStyle = .short
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

#Preview {
    AdminDashboardView()
}
