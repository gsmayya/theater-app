import Foundation

// MARK: - Mock Data Service
class MockDataService {
    static let shared = MockDataService()
    
    private init() {}
    
    // MARK: - Mock Shows Data
    func getMockShows() -> [Show] {
        return [
            Show(
                id: "1",
                name: "Romeo and Juliet",
                title: "Romeo and Juliet",
                details: "Shakespeare's timeless tragedy of star-crossed lovers. A tale of passion, family rivalry, and destiny that continues to captivate audiences worldwide.",
                price: 4500, // $45.00 in cents
                totalTickets: 100,
                bookedTickets: 55,
                location: "Grand Theater",
                showNumber: "RJ001",
                showDate: "2024-09-15",
                showTimes: [
                    ShowTime(id: "st1", date: "2024-09-15", time: "19:00", showId: "1", availableSeats: 45, totalSeats: 100),
                    ShowTime(id: "st2", date: "2024-09-16", time: "14:00", showId: "1", availableSeats: 32, totalSeats: 100),
                    ShowTime(id: "st3", date: "2024-09-16", time: "19:00", showId: "1", availableSeats: 67, totalSeats: 100)
                ],
                images: ["/images/romeo-juliet.jpg"],
                videos: [],
                createdAt: "2024-08-01T10:00:00Z",
                updatedAt: "2024-08-01T10:00:00Z"
            ),
            Show(
                id: "2",
                name: "The Lion King",
                title: "The Lion King",
                details: "Disney's beloved musical brings the African savanna to life with stunning costumes, innovative puppetry, and unforgettable music.",
                price: 6500, // $65.00 in cents
                totalTickets: 120,
                bookedTickets: 8,
                location: "Royal Opera House",
                showNumber: "LK002",
                showDate: "2024-09-20",
                showTimes: [
                    ShowTime(id: "st4", date: "2024-09-20", time: "15:00", showId: "2", availableSeats: 89, totalSeats: 120),
                    ShowTime(id: "st5", date: "2024-09-21", time: "19:30", showId: "2", availableSeats: 112, totalSeats: 120)
                ],
                images: ["/images/lion-king.jpg"],
                videos: [],
                createdAt: "2024-08-01T10:00:00Z",
                updatedAt: "2024-08-01T10:00:00Z"
            ),
            Show(
                id: "3",
                name: "Hamilton",
                title: "Hamilton",
                details: "Lin-Manuel Miranda's revolutionary musical biography of Alexander Hamilton, featuring hip-hop, R&B, and traditional show tunes.",
                price: 9500, // $95.00 in cents
                totalTickets: 150,
                bookedTickets: 119,
                location: "Broadway Theater",
                showNumber: "HM003",
                showDate: "2024-09-25",
                showTimes: [
                    ShowTime(id: "st6", date: "2024-09-25", time: "20:00", showId: "3", availableSeats: 23, totalSeats: 150),
                    ShowTime(id: "st7", date: "2024-09-26", time: "20:00", showId: "3", availableSeats: 8, totalSeats: 150)
                ],
                images: ["/images/hamilton.jpg"],
                videos: [],
                createdAt: "2024-08-01T10:00:00Z",
                updatedAt: "2024-08-01T10:00:00Z"
            ),
            Show(
                id: "4",
                name: "A Midsummer Night's Dream",
                title: "A Midsummer Night's Dream",
                details: "Shakespeare's magical comedy filled with fairies, lovers, and mischief in an enchanted forest setting.",
                price: 4000, // $40.00 in cents
                totalTickets: 90,
                bookedTickets: 5,
                location: "Garden Theater",
                showNumber: "MS004",
                showDate: "2024-09-30",
                showTimes: [
                    ShowTime(id: "st8", date: "2024-09-30", time: "18:00", showId: "4", availableSeats: 78, totalSeats: 90),
                    ShowTime(id: "st9", date: "2024-10-01", time: "18:00", showId: "4", availableSeats: 85, totalSeats: 90)
                ],
                images: ["/images/midsummer.jpg"],
                videos: [],
                createdAt: "2024-08-01T10:00:00Z",
                updatedAt: "2024-08-01T10:00:00Z"
            ),
            Show(
                id: "5",
                name: "The Phantom of the Opera",
                title: "The Phantom of the Opera",
                details: "Andrew Lloyd Webber's masterpiece about a masked figure who lurks beneath the catacombs of the Paris Opera House.",
                price: 7500, // $75.00 in cents
                totalTickets: 200,
                bookedTickets: 45,
                location: "Majestic Theater",
                showNumber: "PO005",
                showDate: "2024-10-05",
                showTimes: [
                    ShowTime(id: "st10", date: "2024-10-05", time: "19:30", showId: "5", availableSeats: 120, totalSeats: 200),
                    ShowTime(id: "st11", date: "2024-10-06", time: "14:00", showId: "5", availableSeats: 135, totalSeats: 200),
                    ShowTime(id: "st12", date: "2024-10-06", time: "19:30", showId: "5", availableSeats: 110, totalSeats: 200)
                ],
                images: ["/images/phantom.jpg"],
                videos: [],
                createdAt: "2024-08-01T10:00:00Z",
                updatedAt: "2024-08-01T10:00:00Z"
            ),
            Show(
                id: "6",
                name: "Les Misérables",
                title: "Les Misérables",
                details: "Victor Hugo's epic tale of redemption set against the backdrop of 19th-century France, featuring unforgettable songs.",
                price: 8500, // $85.00 in cents
                totalTickets: 180,
                bookedTickets: 92,
                location: "Imperial Theater",
                showNumber: "LM006",
                showDate: "2024-10-10",
                showTimes: [
                    ShowTime(id: "st13", date: "2024-10-10", time: "20:00", showId: "6", availableSeats: 65, totalSeats: 180),
                    ShowTime(id: "st14", date: "2024-10-11", time: "20:00", showId: "6", availableSeats: 23, totalSeats: 180)
                ],
                images: ["/images/les-mis.jpg"],
                videos: [],
                createdAt: "2024-08-01T10:00:00Z",
                updatedAt: "2024-08-01T10:00:00Z"
            )
        ]
    }
    
    // MARK: - Mock Users Data
    func getMockUsers() -> [User] {
        return [
            User(
                id: "user1",
                email: "john.doe@example.com",
                name: "John Doe",
                phone: "+1234567890",
                avatar: nil,
                provider: "local"
            ),
            User(
                id: "user2",
                email: "jane.smith@example.com",
                name: "Jane Smith",
                phone: "+1987654321",
                avatar: nil,
                provider: "local"
            )
        ]
    }
    
    // MARK: - Mock Bookings Data
    func getMockBookings() -> [Booking] {
        return [
            Booking(
                id: "booking1",
                showId: "1",
                contactType: "email",
                contactValue: "john.doe@example.com",
                numberOfTickets: 2,
                customerName: "John Doe",
                totalAmount: 9000, // $90.00 in cents
                bookingDate: "2024-08-30",
                status: "confirmed",
                createdAt: "2024-08-30T10:00:00Z",
                updatedAt: "2024-08-30T10:00:00Z"
            ),
            Booking(
                id: "booking2",
                showId: "2",
                contactType: "email",
                contactValue: "jane.smith@example.com",
                numberOfTickets: 4,
                customerName: "Jane Smith",
                totalAmount: 26000, // $260.00 in cents
                bookingDate: "2024-08-30",
                status: "confirmed",
                createdAt: "2024-08-30T14:30:00Z",
                updatedAt: "2024-08-30T14:30:00Z"
            ),
            Booking(
                id: "booking3",
                showId: "3",
                contactType: "mobile",
                contactValue: "+1234567890",
                numberOfTickets: 1,
                customerName: "John Doe",
                totalAmount: 9500, // $95.00 in cents
                bookingDate: "2024-09-01",
                status: "pending",
                createdAt: "2024-09-01T09:00:00Z",
                updatedAt: "2024-09-01T09:00:00Z"
            )
        ]
    }
    
    // MARK: - Search Mock Data
    func searchMockShows(query: String, location: String? = nil) -> [Show] {
        let shows = getMockShows()
        
        var filteredShows = shows
        
        if !query.isEmpty {
            filteredShows = filteredShows.filter { show in
                show.title.localizedCaseInsensitiveContains(query) ||
                show.details.localizedCaseInsensitiveContains(query) ||
                show.name.localizedCaseInsensitiveContains(query)
            }
        }
        
        if let location = location, !location.isEmpty {
            filteredShows = filteredShows.filter { show in
                show.location.localizedCaseInsensitiveContains(location)
            }
        }
        
        return filteredShows
    }
    
    // MARK: - Generate Mock Booking
    func createMockBooking(showId: String, customerName: String, customerEmail: String, numberOfTickets: Int) -> Booking {
        let shows = getMockShows()
        guard let show = shows.first(where: { $0.id == showId }) else {
            fatalError("Show not found")
        }
        
        let totalAmount = show.price * numberOfTickets
        let bookingId = "booking_\(Int.random(in: 1000...9999))"
        
        return Booking(
            id: bookingId,
            showId: showId,
            contactType: "email",
            contactValue: customerEmail,
            numberOfTickets: numberOfTickets,
            customerName: customerName,
            totalAmount: totalAmount,
            bookingDate: DateFormatter.iso8601.string(from: Date()),
            status: "confirmed",
            createdAt: DateFormatter.iso8601.string(from: Date()),
            updatedAt: DateFormatter.iso8601.string(from: Date())
        )
    }
}

// MARK: - Date Formatter Extension
extension DateFormatter {
    static let iso8601: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ss'Z'"
        formatter.timeZone = TimeZone(abbreviation: "UTC")
        return formatter
    }()
}
