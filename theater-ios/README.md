# Theater iOS App

A SwiftUI-based iOS application for the Illuminating Windows Theater Hub, providing a complete theater show management and booking system.

## Features

### 🎭 **Public Features**
- **Show Discovery**: Browse all available theater shows with detailed information
- **Search & Filter**: Search shows by title, location, or description
- **Calendar View**: Interactive calendar showing shows by date
- **Show Details**: Comprehensive show information including pricing, availability, and show times
- **Ticket Booking**: Complete booking flow with customer information and payment
- **QR Code Tickets**: Generate QR codes for ticket validation

### 🔐 **User Authentication**
- **Mock Login**: Simple email/password authentication for development
- **Google Sign-In**: Optional Google authentication integration
- **User Profile**: Manage user information and view booking history
- **Booking Management**: View, confirm, and cancel bookings

### 👨‍💼 **Admin Features**
- **Admin Dashboard**: Comprehensive admin interface for show management
- **Create Shows**: Add new theater performances with full details
- **Show Management**: View and manage all shows in the system
- **Statistics**: View booking statistics and revenue reports
- **Location Analytics**: Shows grouped by venue location

## Architecture

### **MVVM Pattern**
- **Models**: Data structures for Show, Booking, User, and API responses
- **Views**: SwiftUI views for all screens and components
- **Services**: API service, authentication service, and QR code generation
- **ViewModels**: Observable objects managing state and business logic

### **Key Components**

#### **Data Models** (`Models/`)
- `Show.swift`: Theater show data structure with show times and pricing
- `Booking.swift`: Booking information and status tracking
- `User.swift`: User profile and authentication data
- `APIResponse.swift`: Standardized API response handling

#### **Services** (`Services/`)
- `APIService.swift`: RESTful API communication with theater backend
- `AuthService.swift`: Authentication and user session management
- `QRCodeService.swift`: QR code generation for tickets and shows

#### **Views** (`Views/`)
- `HomeView.swift`: Main show listing with search functionality
- `BookingView.swift`: Complete ticket booking flow
- `CalendarView.swift`: Interactive calendar for show dates
- `BookingsView.swift`: User booking history and management
- `ProfileView.swift`: User profile and settings
- `AdminLoginView.swift`: Admin authentication
- `AdminDashboardView.swift`: Admin show management interface
- `CreateShowView.swift`: Form for creating new shows

## API Integration

The app integrates with the theater backend API endpoints:

### **Show Management**
- `GET /api/v1/shows` - Fetch all shows
- `GET /api/v1/search` - Search shows with filters
- `GET /api/v1/shows/get` - Get specific show details
- `POST /api/v1/shows/create` - Create new show (admin)

### **Booking Management**
- `POST /api/v1/bookings/create` - Create new booking
- `GET /api/v1/bookings/get` - Get booking details
- `GET /api/v1/bookings/by-contact` - Get user bookings
- `PUT /api/v1/bookings/confirm` - Confirm booking
- `PUT /api/v1/bookings/cancel` - Cancel booking

### **Authentication**
- Mock authentication for development
- Google Sign-In integration ready
- Admin authentication with role-based access

## Setup Instructions

### **Prerequisites**
- Xcode 15.0 or later
- iOS 17.0 or later
- Theater backend running on `http://localhost:8080`

### **Installation**
1. Open `TheaterApp.xcodeproj` in Xcode
2. Select your target device or simulator
3. Build and run the project (⌘+R)

### **Configuration**
1. **Backend URL**: Update `baseURL` in `APIService.swift` if needed
2. **Authentication**: Modify `AuthService.swift` for production authentication
3. **QR Codes**: QR code generation works out of the box

## Usage

### **For Users**
1. **Browse Shows**: Open the app to see all available shows
2. **Search**: Use the search bar to find specific shows
3. **View Calendar**: Tap Calendar tab to see shows by date
4. **Book Tickets**: Tap "Book Tickets" on any show
5. **Sign In**: Create account or sign in to manage bookings
6. **View Bookings**: Check "My Bookings" tab for your tickets

### **For Admins**
1. **Admin Login**: Go to Profile → Admin Login
2. **Use Demo Credentials**:
   - Email: `admin@theater.com`
   - Password: `admin123`
3. **Create Shows**: Use the admin dashboard to add new shows
4. **Manage Shows**: View and manage all shows in the system
5. **View Statistics**: Check booking statistics and revenue

## Features in Detail

### **Show Discovery**
- Grid layout with show cards
- Search by title, location, or description
- Filter by availability and price
- Show times and availability indicators
- Price display in dollars

### **Booking Flow**
1. Select show and show time
2. Choose number of tickets
3. Enter customer information
4. Confirm booking
5. Receive QR code ticket

### **Calendar Integration**
- Monthly calendar view
- Shows marked on available dates
- Tap date to see events
- Event cards with booking options

### **Admin Dashboard**
- Create new shows with full details
- View all shows in organized list
- Statistics cards for quick overview
- Location-based show grouping
- Revenue and booking analytics

### **QR Code System**
- Generate QR codes for bookings
- Include booking ID and show information
- QR code display for ticket validation
- Scanner functionality (placeholder)

## Technical Details

### **Dependencies**
- **SwiftUI**: Modern UI framework
- **Combine**: Reactive programming for API calls
- **CoreImage**: QR code generation
- **Foundation**: Core data handling

### **State Management**
- `@StateObject` for service instances
- `@Published` properties for reactive updates
- `Combine` publishers for API communication
- `@Environment` for dependency injection

### **Error Handling**
- Comprehensive error handling for API calls
- User-friendly error messages
- Loading states and empty states
- Network error recovery

### **Performance**
- Lazy loading for show lists
- Image placeholders and caching
- Efficient calendar rendering
- Memory-conscious data handling

## Development Notes

### **Mock Data**
- Authentication uses mock login for development
- All API calls work with real backend
- Fallback to mock data if backend unavailable

### **Testing**
- Preview providers for all views
- Mock data for development
- Error state testing
- Loading state validation

### **Future Enhancements**
- Push notifications for booking confirmations
- Apple Pay integration
- Offline mode support
- Advanced search filters
- Social sharing features
- Real-time availability updates

## Troubleshooting

### **Common Issues**
1. **API Connection**: Ensure backend is running on port 8080
2. **Authentication**: Use demo credentials for testing
3. **QR Codes**: QR code generation requires iOS 17.0+
4. **Calendar**: Date formatting may vary by locale

### **Debug Mode**
- Enable debug logging in services
- Check console for API errors
- Verify network connectivity
- Test with different show data

## Contributing

1. Follow SwiftUI best practices
2. Maintain MVVM architecture
3. Add proper error handling
4. Include preview providers
5. Document new features

## License

This project is part of the Illuminating Windows Theater Hub system.