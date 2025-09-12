import Foundation
import Combine

// MARK: - API Service
class APIService: ObservableObject {
    static let shared = APIService()
    
    private let baseURL = "http://localhost:8080/api/v1"
    private let session = URLSession.shared
    private let mockDataService = MockDataService.shared
    private var useMockData = false
    
    private init() {
        // Check if backend is available on startup
        checkBackendAvailability()
    }
    
    // MARK: - Backend Availability Check
    private func checkBackendAvailability() {
        // For now, we'll use mock data by default
        // In a real app, you might want to ping the backend first
        useMockData = true
    }
    
    // MARK: - Configuration Methods
    func setUseMockData(_ useMock: Bool) {
        useMockData = useMock
    }
    
    var isUsingMockData: Bool {
        return useMockData
    }
    
    // MARK: - Show Endpoints
    
    func fetchShows() -> AnyPublisher<[Show], Error> {
        if useMockData {
            // Return mock data with a small delay to simulate network
            return Just(mockDataService.getMockShows())
                .delay(for: .milliseconds(500), scheduler: DispatchQueue.main)
                .setFailureType(to: Error.self)
                .eraseToAnyPublisher()
        }
        
        guard let url = URL(string: "\(baseURL)/shows") else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map(\.data)
            .decode(type: [Show].self, decoder: JSONDecoder())
            .catch { _ in
                // Fallback to mock data on error
                Just(self.mockDataService.getMockShows())
                    .setFailureType(to: Error.self)
            }
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func searchShows(params: SearchParams) -> AnyPublisher<[Show], Error> {
        if useMockData {
            // Use mock data for search
            let searchQuery = params.search ?? ""
            let location = params.location
            let results = mockDataService.searchMockShows(query: searchQuery, location: location)
            
            return Just(results)
                .delay(for: .milliseconds(300), scheduler: DispatchQueue.main)
                .setFailureType(to: Error.self)
                .eraseToAnyPublisher()
        }
        
        var components = URLComponents(string: "\(baseURL)/search")
        components?.queryItems = params.queryItems
        
        guard let url = components?.url else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map(\.data)
            .decode(type: [Show].self, decoder: JSONDecoder())
            .catch { _ in
                // Fallback to mock data on error
                let searchQuery = params.search ?? ""
                let location = params.location
                let results = self.mockDataService.searchMockShows(query: searchQuery, location: location)
                return Just(results)
                    .setFailureType(to: Error.self)
            }
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func getShow(by id: String) -> AnyPublisher<Show, Error> {
        var components = URLComponents(string: "\(baseURL)/shows/get")
        components?.queryItems = [URLQueryItem(name: "show_id", value: id)]
        
        guard let url = components?.url else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map(\.data)
            .decode(type: Show.self, decoder: JSONDecoder())
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func createShow(_ show: CreateShowRequest) -> AnyPublisher<Show, Error> {
        guard let url = URL(string: "\(baseURL)/shows/create") else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            request.httpBody = try JSONEncoder().encode(show)
        } catch {
            return Fail(error: error)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: request)
            .map(\.data)
            .decode(type: Show.self, decoder: JSONDecoder())
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    // MARK: - Booking Endpoints
    
    func createBooking(_ booking: BookingRequest) -> AnyPublisher<Booking, Error> {
        if useMockData {
            // Create mock booking
            let mockBooking = mockDataService.createMockBooking(
                showId: booking.showId,
                customerName: booking.customerName ?? "Guest",
                customerEmail: booking.contactValue,
                numberOfTickets: booking.numberOfTickets
            )
            
            return Just(mockBooking)
                .delay(for: .milliseconds(800), scheduler: DispatchQueue.main)
                .setFailureType(to: Error.self)
                .eraseToAnyPublisher()
        }
        
        guard let url = URL(string: "\(baseURL)/bookings/create") else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            request.httpBody = try JSONEncoder().encode(booking)
        } catch {
            return Fail(error: error)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: request)
            .map(\.data)
            .decode(type: Booking.self, decoder: JSONDecoder())
            .catch { _ in
                // Fallback to mock data on error
                let mockBooking = self.mockDataService.createMockBooking(
                    showId: booking.showId,
                    customerName: booking.customerName ?? "Guest",
                    customerEmail: booking.contactValue,
                    numberOfTickets: booking.numberOfTickets
                )
                return Just(mockBooking)
                    .setFailureType(to: Error.self)
            }
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func getBooking(by id: String) -> AnyPublisher<Booking, Error> {
        var components = URLComponents(string: "\(baseURL)/bookings/get")
        components?.queryItems = [URLQueryItem(name: "booking_id", value: id)]
        
        guard let url = components?.url else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map(\.data)
            .decode(type: Booking.self, decoder: JSONDecoder())
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func getBookingsByShow(showId: String) -> AnyPublisher<[Booking], Error> {
        var components = URLComponents(string: "\(baseURL)/bookings/by-show")
        components?.queryItems = [URLQueryItem(name: "show_id", value: showId)]
        
        guard let url = components?.url else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map(\.data)
            .decode(type: [Booking].self, decoder: JSONDecoder())
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func getBookingsByContact(contactValue: String) -> AnyPublisher<[Booking], Error> {
        if useMockData {
            // Return mock bookings filtered by contact
            let allBookings = mockDataService.getMockBookings()
            let filteredBookings = allBookings.filter { booking in
                booking.contactValue == contactValue
            }
            
            return Just(filteredBookings)
                .delay(for: .milliseconds(400), scheduler: DispatchQueue.main)
                .setFailureType(to: Error.self)
                .eraseToAnyPublisher()
        }
        
        var components = URLComponents(string: "\(baseURL)/bookings/by-contact")
        components?.queryItems = [URLQueryItem(name: "contact_value", value: contactValue)]
        
        guard let url = components?.url else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map(\.data)
            .decode(type: [Booking].self, decoder: JSONDecoder())
            .catch { _ in
                // Fallback to mock data on error
                let allBookings = self.mockDataService.getMockBookings()
                let filteredBookings = allBookings.filter { booking in
                    booking.contactValue == contactValue
                }
                return Just(filteredBookings)
                    .setFailureType(to: Error.self)
            }
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func confirmBooking(bookingId: String) -> AnyPublisher<Booking, Error> {
        guard let url = URL(string: "\(baseURL)/bookings/confirm") else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let confirmRequest = ["booking_id": bookingId]
        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: confirmRequest)
        } catch {
            return Fail(error: error)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: request)
            .map(\.data)
            .decode(type: Booking.self, decoder: JSONDecoder())
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    func cancelBooking(bookingId: String) -> AnyPublisher<Booking, Error> {
        guard let url = URL(string: "\(baseURL)/bookings/cancel") else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let cancelRequest = ["booking_id": bookingId]
        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: cancelRequest)
        } catch {
            return Fail(error: error)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: request)
            .map(\.data)
            .decode(type: Booking.self, decoder: JSONDecoder())
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
    
    // MARK: - Health Check
    
    func healthCheck() -> AnyPublisher<Bool, Error> {
        guard let url = URL(string: "\(baseURL)/health") else {
            return Fail(error: APIError.invalidURL)
                .eraseToAnyPublisher()
        }
        
        return session.dataTaskPublisher(for: url)
            .map { response in
                return (response.response as? HTTPURLResponse)?.statusCode == 200
            }
            .mapError { error in
                return error as Error
            }
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }
}

// MARK: - Create Show Request
struct CreateShowRequest: Codable {
    let name: String
    let details: String
    let price: Int
    let totalTickets: Int
    let location: String
    let showNumber: String
    let showDate: String
    
    enum CodingKeys: String, CodingKey {
        case name
        case details
        case price
        case totalTickets = "total_tickets"
        case location
        case showNumber = "show_number"
        case showDate = "show_date"
    }
}

// MARK: - API Errors
enum APIError: Error, LocalizedError {
    case invalidURL
    case noData
    case decodingError
    case networkError(Error)
    case serverError(Int)
    
    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid URL"
        case .noData:
            return "No data received"
        case .decodingError:
            return "Failed to decode response"
        case .networkError(let error):
            return "Network error: \(error.localizedDescription)"
        case .serverError(let code):
            return "Server error: \(code)"
        }
    }
}
