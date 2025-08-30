# Theater Hub - Theater Show Management System

A Next.js application for managing theater shows with calendar views, ticket booking, and QR code generation.

## Features

- 🎭 **Current Shows & Calendar**: Public page showing current theater performances with interactive calendar
- 🎫 **Ticket Booking**: Authenticated booking system with QR code generation
- 🔐 **Multiple Login Options**: Mock login and Google Sign-In support
- 📱 **Responsive Design**: Fully responsive with Tailwind CSS
- 🎨 **Modern UI**: Clean, theater-themed design with custom color palette

## Quick Start

1. **Install dependencies**
   ```bash
   npm install
   ```

2. **Start the development server**
   ```bash
   npm run dev
   ```

3. **Visit the application**
- Home page: http://localhost:3000 (shows + search + calendar)
   - Booking: http://localhost:3000/booking (login + book tickets)
   - Admin: http://localhost:3000/admin (login page)
   - Admin Dashboard: http://localhost:3000/admin/dashboard (create shows)

## Authentication Options

### Mock Login (Default)
- Works out of the box with no setup required
- Use any email and password to log in
- Perfect for development and testing

### Google Sign-In (Optional)
To enable Google Sign-In:

1. **Set up Firebase project**
   - Go to [Firebase Console](https://console.firebase.google.com)
   - Create a new project or use an existing one
   - Enable Authentication > Sign-in method > Google

2. **Get Firebase config**
   - Go to Project Settings > General
   - Add a web app if you haven't already
   - Copy the config values

3. **Create environment file**
   ```bash
   cp .env.local.example .env.local
   ```

4. **Add your Firebase config to `.env.local`**
   ```bash
   NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key-here
   NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
   NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
   ```

5. **Restart the development server**
   ```bash
   npm run dev
   ```

Now you'll see both "Login" and "Continue with Google" options on login pages.

### Backend Integration (Optional)
To connect to your own backend API:

1. **Set backend URL in environment file**
   ```bash
   # Add to .env.local
   NEXT_PUBLIC_API_BASE_URL=https://your-backend-domain.com
   API_KEY=your-backend-api-key
   ```

2. **Backend API Endpoints**
   Your backend should implement these endpoints:
   - `GET /api/v1/get` - List all shows
   - `GET /api/v1/search?location=&title=` - Search shows by location/title
   - `POST /api/v1/shows` - Create new show
   - `POST /api/v1/bookings` - Create booking (optional)

3. **Show Data Format**
   ```json
   {
     "id": "string",
     "title": "string",
     "description": "string",
     "venue": "string",
     "showTimes": [...],
     "ticketPrice": 50,
     "availableTickets": 100
   }
   ```

4. **Authentication**
   - API requests include `Authorization: Bearer {API_KEY}` header
   - Customize auth in `src/lib/apiService.ts`

Without backend configuration, the app uses mock data automatically.

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main layout with navigation
│   ├── Navigation.tsx  # Navigation bar
│   ├── ShowCard.tsx    # Individual show display
│   └── Calendar.tsx    # Interactive calendar
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication state management
├── lib/               # Utilities and configuration
│   ├── firebaseClient.ts # Safe Firebase initialization
│   └── mockData.ts    # Sample data for development
├── pages/             # Next.js pages
│   ├── api/           # API routes
│   ├── index.tsx      # Home page (shows + calendar)
│   ├── booking.tsx    # Ticket booking page
│   └── admin.tsx      # Admin login page
├── styles/            # Global styles
└── types/             # TypeScript type definitions
```

## API Routes

### Frontend Routes
- `GET /api/shows` - List all shows (with backend fallback)
- `GET /api/search` - Search shows by location/title
- `GET /api/shows/[id]` - Get show details
- `POST /api/shows/create` - Create new show
- `POST /api/bookings` - Create booking with QR code
- `GET /api/bookings` - List bookings

### Backend Integration
When `NEXT_PUBLIC_API_BASE_URL` is configured, calls:
- `GET {baseUrl}/api/v1/get` - List shows
- `GET {baseUrl}/api/v1/search` - Search shows
- `POST {baseUrl}/api/v1/shows` - Create show

## Technologies Used

- **Next.js 15** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Firebase Auth** - Google Sign-In (optional)
- **QRCode library** - QR code generation
- **React Context** - State management

## Development

The app uses mock data by default, stored in `src/lib/mockData.ts`. This includes:
- Sample theater shows with showtimes
- Mock users and bookings
- Realistic theater data (Romeo & Juliet, Lion King, Hamilton, etc.)

All booking operations work with this mock data, including:
- Seat availability tracking
- QR code generation
- Booking confirmation

## Building for Production

```bash
npm run build
npm start
```

## License

MIT License - feel free to use this project for your theater management needs!

### Test 

Backend Integration Features

🔧 Configurable Backend Connection
•  Added NEXT_PUBLIC_API_BASE_URL and API_KEY environment variables
•  Created src/lib/apiService.ts - centralized service for backend API calls
•  Automatic fallback to mock data when backend is not configured

📡 API Endpoints Implemented
•  GET /api/v1/search - Search shows by location and/or title with pagination
•  GET /api/v1/get - Get all shows (mapped from your backend)  
•  POST /api/v1/shows - Create new show with title, description, location, date

🎯 Frontend Integration
•  Search functionality: Added SearchBar component on homepage that calls /api/search
•  Show creation: Added CreateShowForm component for admin dashboard
•  Admin dashboard: New /admin/dashboard page with show management
•  Smart fallback: All endpoints gracefully fall back to mock data if backend fails

🛠 How It Works

1. Without Backend (Default):
•  Uses mock data from src/lib/mockData.ts
•  All functionality works immediately
2. With Backend (Configured):
•  Set NEXT_PUBLIC_API_BASE_URL=https://your-backend-domain.com
•  Calls your backend at /api/v1/ endpoints
•  Includes Authorization: Bearer {API_KEY} headers

🎨 UI Enhancements
•  Added search bar to homepage with location/title filtering
•  Created admin dashboard at /admin/dashboard (requires login)
•  Enhanced navigation with Admin link
•  Professional forms with loading states and error handling

📋 Usage Examples

Search shows:

  ```bash
  curl "http://localhost:3000/api/search?location=Broadway&title=Romeo"
```

Create show:
```bash
curl -X POST "http://localhost:3000/api/shows/create" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hamlet",
    "description": "Shakespeare's masterpiece",
    "location": "Globe Theater", 
    "date": "2024-12-01"
  }'
  ```


