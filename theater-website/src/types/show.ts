// ShowTime interface for individual show times
export interface ShowTime {
  id: string;
  date: string; // ISO date string
  time: string; // Time string (e.g., "19:30")
  show_id: string;
  availableSeats?: number; // Optional for backward compatibility
  totalSeats?: number; // Optional for backward compatibility
}

// Core Show interface matching Go backend structure
export interface Show {
  id: string;
  name: string;
  title: string; // Added for Calendar component compatibility
  details: string;
  price: number; // Price in cents
  total_tickets: number;
  booked_tickets: number;
  location: string;
  show_number: string;
  show_date: string; // ISO date string
  showTimes: ShowTime[]; // Array of show times for Calendar component
  images: string[]; // Array of CMS image IDs
  videos: string[]; // Array of CMS video IDs
  created_at: string;
  updated_at: string;
}

// Booking interface matching Go backend structure
export interface Booking {
  booking_id: string;
  show_id: string;
  contact_type: 'mobile' | 'email';
  contact_value: string;
  number_of_tickets: number;
  customer_name?: string | undefined;
  total_amount: number;
  booking_date: string; // ISO date string
  status: 'pending' | 'confirmed' | 'cancelled';
  created_at: string;
  updated_at: string;
}

// User interface for authentication
export interface User {
  id: string;
  email: string;
  name: string;
  phone?: string;
  avatar?: string;
  provider?: 'google' | 'local';
}

// API Response interfaces
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
  status_code: number;
}

export interface PaginatedResponse<T = any> {
  success: boolean;
  message: string;
  data: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

// Search and filter interfaces
export interface SearchParams {
  location?: string;
  min_price?: number;
  max_price?: number;
  min_available?: number;
  search?: string;
  only_available?: boolean;
  page?: number;
  page_size?: number;
}

export interface BookingRequest {
  show_id: string;
  contact_type: 'mobile' | 'email';
  contact_value: string;
  number_of_tickets: number;
  customer_name?: string;
  booking_date: string;
}

// Calendar event interface
export interface CalendarEvent {
  id: string;
  title: string;
  date: string;
  time: string;
  show_id: string;
  location: string;
  price: number;
  available_tickets: number;
}

// Ticket interface for QR code generation
export interface Ticket {
  booking_id: string;
  show_name: string;
  show_date: string;
  show_location: string;
  customer_name: string;
  number_of_tickets: number;
  total_amount: number;
  status: string;
  qr_data: string;
}

// Show statistics interface
export interface ShowStats {
  show_id: string;
  name: string;
  total_bookings: number;
  confirmed_revenue: number;
  pending_revenue: number;
  confirmed_bookings: number;
  pending_bookings: number;
  cancelled_bookings: number;
  available_tickets: number;
}

// Error interface
export interface ApiError {
  message: string;
  code: string;
  details?: Record<string, any>;
}

// Form validation interfaces
export interface BookingFormData {
  contact_type: 'mobile' | 'email';
  contact_value: string;
  number_of_tickets: number;
  customer_name?: string;
}

export interface SearchFormData {
  location: string;
  min_price: string;
  max_price: string;
  search: string;
  only_available: boolean;
}

// Environment configuration
export interface AppConfig {
  apiUrl: string;
  isProduction: boolean;
  enableMockData: boolean;
  googleClientId?: string | undefined;
  firebaseConfig?: {
    apiKey: string;
    authDomain: string;
    projectId: string;
    storageBucket: string;
    messagingSenderId: string;
    appId: string;
  } | undefined;
}
