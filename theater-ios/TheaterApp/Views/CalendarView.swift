import SwiftUI
import Combine

struct CalendarView: View {
    @StateObject private var apiService = APIService.shared
    @State private var shows: [Show] = []
    @State private var currentDate = Date()
    @State private var selectedDate: Date?
    @State private var selectedEvents: [CalendarEvent] = []
    @State private var isLoading = false
    @State private var errorMessage: String?
    @State private var showingBooking = false
    @State private var selectedShow: Show?
    
    private var calendarEvents: [CalendarEvent] {
        shows.flatMap { show in
            show.showTimes.map { showTime in
                CalendarEvent(from: show, showTime: showTime)
            }
        }
    }
    
    private var eventsForSelectedDate: [CalendarEvent] {
        guard let selectedDate = selectedDate else { return [] }
        let dateString = dateFormatter.string(from: selectedDate)
        return calendarEvents.filter { $0.date == dateString }
    }
    
    private let dateFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd"
        return formatter
    }()
    
    private let displayFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        return formatter
    }()
    
    var body: some View {
        NavigationView {
            VStack(spacing: 0) {
                // Calendar Header
                calendarHeaderView
                
                // Calendar Grid
                calendarGridView
                
                // Selected Date Events
                if !eventsForSelectedDate.isEmpty {
                    selectedDateEventsView
                } else if selectedDate != nil {
                    noEventsView
                }
                
                Spacer()
            }
            .navigationTitle("Calendar")
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
    
    // MARK: - Calendar Header View
    private var calendarHeaderView: some View {
        HStack {
            Button(action: previousMonth) {
                Image(systemName: "chevron.left")
                    .font(.title2)
                    .foregroundColor(.blue)
            }
            
            Spacer()
            
            Text(monthYearString)
                .font(.title2)
                .fontWeight(.semibold)
            
            Spacer()
            
            Button(action: nextMonth) {
                Image(systemName: "chevron.right")
                    .font(.title2)
                    .foregroundColor(.blue)
            }
        }
        .padding()
        .background(Color(.systemGroupedBackground))
    }
    
    // MARK: - Calendar Grid View
    private var calendarGridView: some View {
        VStack(spacing: 0) {
            // Days of week header
            HStack {
                ForEach(dayHeaders, id: \.self) { day in
                    Text(day)
                        .font(.caption)
                        .fontWeight(.medium)
                        .foregroundColor(.secondary)
                        .frame(maxWidth: .infinity)
                }
            }
            .padding(.vertical, 8)
            
            // Calendar days
            LazyVGrid(columns: Array(repeating: GridItem(.flexible()), count: 7), spacing: 1) {
                ForEach(calendarDays, id: \.self) { day in
                    CalendarDayView(
                        day: day,
                        currentDate: currentDate,
                        selectedDate: selectedDate,
                        events: eventsForDay(day),
                        onTap: { selectDate(day) }
                    )
                }
            }
        }
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .padding(.horizontal)
        .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
    }
    
    // MARK: - Selected Date Events View
    private var selectedDateEventsView: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                Text("Events on \(displayFormatter.string(from: selectedDate!))")
                    .font(.headline)
                    .fontWeight(.semibold)
                
                Spacer()
                
                Text("\(eventsForSelectedDate.count) event\(eventsForSelectedDate.count == 1 ? "" : "s")")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            
            ScrollView {
                LazyVStack(spacing: 8) {
                    ForEach(eventsForSelectedDate) { event in
                        EventCardView(event: event) {
                            if let show = shows.first(where: { $0.id == event.showId }) {
                                selectedShow = show
                                showingBooking = true
                            }
                        }
                    }
                }
            }
            .frame(maxHeight: 300)
        }
        .padding()
        .background(Color(.systemGroupedBackground))
        .cornerRadius(12)
        .padding(.horizontal)
    }
    
    // MARK: - No Events View
    private var noEventsView: some View {
        VStack(spacing: 12) {
            Image(systemName: "calendar.badge.plus")
                .font(.system(size: 40))
                .foregroundColor(.secondary)
            
            Text("No events on this date")
                .font(.headline)
                .foregroundColor(.secondary)
            
            Text("Select another date to view available shows")
                .font(.caption)
                .foregroundColor(.secondary)
                .multilineTextAlignment(.center)
        }
        .padding()
        .frame(maxWidth: .infinity)
        .background(Color(.systemGroupedBackground))
        .cornerRadius(12)
        .padding(.horizontal)
    }
    
    // MARK: - Computed Properties
    private var monthYearString: String {
        let formatter = DateFormatter()
        formatter.dateFormat = "MMMM yyyy"
        return formatter.string(from: currentDate)
    }
    
    private var dayHeaders: [String] {
        ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]
    }
    
    private var calendarDays: [Int] {
        let calendar = Calendar.current
        let firstDayOfMonth = calendar.dateInterval(of: .month, for: currentDate)?.start ?? currentDate
        let firstWeekday = calendar.component(.weekday, from: firstDayOfMonth) - 1
        let daysInMonth = calendar.range(of: .day, in: .month, for: currentDate)?.count ?? 30
        
        var days: [Int] = []
        
        // Add empty days for the first week
        for _ in 0..<firstWeekday {
            days.append(0)
        }
        
        // Add days of the month
        for day in 1...daysInMonth {
            days.append(day)
        }
        
        return days
    }
    
    // MARK: - Methods
    private func previousMonth() {
        withAnimation(.easeInOut(duration: 0.3)) {
            currentDate = Calendar.current.date(byAdding: .month, value: -1, to: currentDate) ?? currentDate
        }
    }
    
    private func nextMonth() {
        withAnimation(.easeInOut(duration: 0.3)) {
            currentDate = Calendar.current.date(byAdding: .month, value: 1, to: currentDate) ?? currentDate
        }
    }
    
    private func selectDate(_ day: Int) {
        guard day > 0 else { return }
        
        let calendar = Calendar.current
        let year = calendar.component(.year, from: currentDate)
        let month = calendar.component(.month, from: currentDate)
        
        if let date = calendar.date(from: DateComponents(year: year, month: month, day: day)) {
            selectedDate = date
        }
    }
    
    private func eventsForDay(_ day: Int) -> [CalendarEvent] {
        guard day > 0 else { return [] }
        
        let calendar = Calendar.current
        let year = calendar.component(.year, from: currentDate)
        let month = calendar.component(.month, from: currentDate)
        
        guard let date = calendar.date(from: DateComponents(year: year, month: month, day: day)) else {
            return []
        }
        
        let dateString = dateFormatter.string(from: date)
        return calendarEvents.filter { $0.date == dateString }
    }
    
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

// MARK: - Calendar Day View
struct CalendarDayView: View {
    let day: Int
    let currentDate: Date
    let selectedDate: Date?
    let events: [CalendarEvent]
    let onTap: () -> Void
    
    private var isToday: Bool {
        guard day > 0 else { return false }
        let calendar = Calendar.current
        let today = Date()
        let year = calendar.component(.year, from: currentDate)
        let month = calendar.component(.month, from: currentDate)
        
        guard let date = calendar.date(from: DateComponents(year: year, month: month, day: day)) else {
            return false
        }
        
        return calendar.isDate(date, inSameDayAs: today)
    }
    
    private var isSelected: Bool {
        guard day > 0, let selectedDate = selectedDate else { return false }
        let calendar = Calendar.current
        let year = calendar.component(.year, from: currentDate)
        let month = calendar.component(.month, from: currentDate)
        
        guard let date = calendar.date(from: DateComponents(year: year, month: month, day: day)) else {
            return false
        }
        
        return calendar.isDate(date, inSameDayAs: selectedDate)
    }
    
    private var isPast: Bool {
        guard day > 0 else { return false }
        let calendar = Calendar.current
        let today = Date()
        let year = calendar.component(.year, from: currentDate)
        let month = calendar.component(.month, from: currentDate)
        
        guard let date = calendar.date(from: DateComponents(year: year, month: month, day: day)) else {
            return false
        }
        
        return date < today
    }
    
    var body: some View {
        Button(action: onTap) {
            VStack(spacing: 2) {
                Text(day > 0 ? "\(day)" : "")
                    .font(.system(size: 16, weight: .medium))
                    .foregroundColor(textColor)
                
                if !events.isEmpty {
                    HStack(spacing: 1) {
                        ForEach(events.prefix(3), id: \.id) { _ in
                            Circle()
                                .fill(eventDotColor)
                                .frame(width: 4, height: 4)
                        }
                        
                        if events.count > 3 {
                            Circle()
                                .fill(eventDotColor.opacity(0.5))
                                .frame(width: 4, height: 4)
                        }
                    }
                }
            }
            .frame(height: 44)
            .frame(maxWidth: .infinity)
            .background(backgroundColor)
            .cornerRadius(8)
        }
        .disabled(day == 0 || isPast)
        .buttonStyle(PlainButtonStyle())
    }
    
    private var textColor: Color {
        if day == 0 { return .clear }
        if isPast { return .secondary }
        if isSelected { return .white }
        if isToday { return .blue }
        return .primary
    }
    
    private var backgroundColor: Color {
        if day == 0 { return .clear }
        if isSelected { return .blue }
        if isToday { return .blue.opacity(0.1) }
        if !events.isEmpty { return .blue.opacity(0.05) }
        return .clear
    }
    
    private var eventDotColor: Color {
        if isSelected || isToday { return .white }
        return .blue
    }
}

// MARK: - Event Card View
struct EventCardView: View {
    let event: CalendarEvent
    let onTap: () -> Void
    
    var body: some View {
        Button(action: onTap) {
            HStack(spacing: 12) {
                // Time
                VStack(alignment: .leading, spacing: 2) {
                    Text(formatTime(event.time))
                        .font(.headline)
                        .foregroundColor(.primary)
                    
                    Text("Show Time")
                        .font(.caption)
                        .foregroundColor(.secondary)
                }
                .frame(width: 80, alignment: .leading)
                
                // Event Details
                VStack(alignment: .leading, spacing: 4) {
                    Text(event.title)
                        .font(.headline)
                        .foregroundColor(.primary)
                        .lineLimit(2)
                    
                    Text(event.location)
                        .font(.subheadline)
                        .foregroundColor(.secondary)
                    
                    HStack {
                        Text(String(format: "$%.2f", Double(event.price) / 100.0))
                            .font(.subheadline)
                            .fontWeight(.semibold)
                            .foregroundColor(.blue)
                        
                        Spacer()
                        
                        Text("\(event.availableTickets) seats")
                            .font(.caption)
                            .foregroundColor(availabilityColor)
                    }
                }
                
                Spacer()
                
                // Arrow
                Image(systemName: "chevron.right")
                    .font(.caption)
                    .foregroundColor(.secondary)
            }
            .padding()
            .background(Color(.systemBackground))
            .cornerRadius(12)
            .shadow(color: .black.opacity(0.1), radius: 2, x: 0, y: 1)
        }
        .buttonStyle(PlainButtonStyle())
    }
    
    private var availabilityColor: Color {
        if event.availableTickets > 10 { return .green }
        if event.availableTickets > 5 { return .orange }
        return .red
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
    CalendarView()
}
