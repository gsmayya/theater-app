import Foundation

// MARK: - Show Models
struct Show: Codable, Identifiable {
    let id: String
    let name: String
    let title: String
    let details: String
    let price: Int // Price in cents
    let totalTickets: Int
    let bookedTickets: Int
    let location: String
    let showNumber: String
    let showDate: String
    let showTimes: [ShowTime]
    let images: [String]
    let videos: [String]
    let createdAt: String
    let updatedAt: String
    
    enum CodingKeys: String, CodingKey {
        case id
        case name
        case title
        case details
        case price
        case totalTickets = "total_tickets"
        case bookedTickets = "booked_tickets"
        case location
        case showNumber = "show_number"
        case showDate = "show_date"
        case showTimes = "showTimes"
        case images
        case videos
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
    
    var availableTickets: Int {
        return totalTickets - bookedTickets
    }
    
    var priceInDollars: Double {
        return Double(price) / 100.0
    }
}

struct ShowTime: Codable, Identifiable {
    let id: String
    let date: String
    let time: String
    let showId: String
    let availableSeats: Int?
    let totalSeats: Int?
    
    enum CodingKeys: String, CodingKey {
        case id
        case date
        case time
        case showId = "show_id"
        case availableSeats = "availableSeats"
        case totalSeats = "totalSeats"
    }
}

// MARK: - Booking Models
struct Booking: Codable, Identifiable {
    let id: String
    let showId: String
    let contactType: String
    let contactValue: String
    let numberOfTickets: Int
    let customerName: String?
    let totalAmount: Int
    let bookingDate: String
    let status: String
    let createdAt: String
    let updatedAt: String
    
    enum CodingKeys: String, CodingKey {
        case id = "booking_id"
        case showId = "show_id"
        case contactType = "contact_type"
        case contactValue = "contact_value"
        case numberOfTickets = "number_of_tickets"
        case customerName = "customer_name"
        case totalAmount = "total_amount"
        case bookingDate = "booking_date"
        case status
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
    
    var totalAmountInDollars: Double {
        return Double(totalAmount) / 100.0
    }
    
    // Computed property for backward compatibility
    var bookingId: String {
        return id
    }
}

struct BookingRequest: Codable {
    let showId: String
    let contactType: String
    let contactValue: String
    let numberOfTickets: Int
    let customerName: String?
    let bookingDate: String
    
    enum CodingKeys: String, CodingKey {
        case showId = "show_id"
        case contactType = "contact_type"
        case contactValue = "contact_value"
        case numberOfTickets = "number_of_tickets"
        case customerName = "customer_name"
        case bookingDate = "booking_date"
    }
}

// MARK: - User Models
struct User: Codable, Identifiable {
    let id: String
    let email: String
    let name: String
    let phone: String?
    let avatar: String?
    let provider: String?
}

// MARK: - API Response Models
struct APIResponse<T: Codable>: Codable {
    let success: Bool
    let message: String
    let data: T?
    let error: String?
    let statusCode: Int
    
    enum CodingKeys: String, CodingKey {
        case success
        case message
        case data
        case error
        case statusCode = "status_code"
    }
}

struct PaginatedResponse<T: Codable>: Codable {
    let success: Bool
    let message: String
    let data: [T]
    let pagination: PaginationInfo
    let statusCode: Int
    
    enum CodingKeys: String, CodingKey {
        case success
        case message
        case data
        case pagination
        case statusCode = "status_code"
    }
}

struct PaginationInfo: Codable {
    let page: Int
    let pageSize: Int
    let total: Int
    let totalPages: Int
    
    enum CodingKeys: String, CodingKey {
        case page
        case pageSize = "page_size"
        case total
        case totalPages = "total_pages"
    }
}

// MARK: - Search Models
struct SearchParams {
    let location: String?
    let minPrice: Int?
    let maxPrice: Int?
    let minAvailable: Int?
    let search: String?
    let onlyAvailable: Bool?
    let page: Int?
    let pageSize: Int?
    
    var queryItems: [URLQueryItem] {
        var items: [URLQueryItem] = []
        
        if let location = location, !location.isEmpty {
            items.append(URLQueryItem(name: "location", value: location))
        }
        if let minPrice = minPrice {
            items.append(URLQueryItem(name: "min_price", value: String(minPrice)))
        }
        if let maxPrice = maxPrice {
            items.append(URLQueryItem(name: "max_price", value: String(maxPrice)))
        }
        if let minAvailable = minAvailable {
            items.append(URLQueryItem(name: "min_available", value: String(minAvailable)))
        }
        if let search = search, !search.isEmpty {
            items.append(URLQueryItem(name: "search", value: search))
        }
        if let onlyAvailable = onlyAvailable {
            items.append(URLQueryItem(name: "only_available", value: String(onlyAvailable)))
        }
        if let page = page {
            items.append(URLQueryItem(name: "page", value: String(page)))
        }
        if let pageSize = pageSize {
            items.append(URLQueryItem(name: "page_size", value: String(pageSize)))
        }
        
        return items
    }
}

// MARK: - Calendar Event Model
struct CalendarEvent: Identifiable {
    let id: String
    let title: String
    let date: String
    let time: String
    let showId: String
    let location: String
    let price: Int
    let availableTickets: Int
    
    init(from show: Show, showTime: ShowTime) {
        self.id = "\(show.id)-\(showTime.id)"
        self.title = show.title
        self.date = showTime.date
        self.time = showTime.time
        self.showId = show.id
        self.location = show.location
        self.price = show.price
        self.availableTickets = showTime.availableSeats ?? show.availableTickets
    }
}
