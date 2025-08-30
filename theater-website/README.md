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
   - Home page: http://localhost:3000
   - Booking: http://localhost:3000/booking
   - Admin: http://localhost:3000/admin

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

- `GET /api/shows` - List all shows
- `GET /api/shows/[id]` - Get show details
- `POST /api/bookings` - Create booking with QR code
- `GET /api/bookings` - List bookings

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
