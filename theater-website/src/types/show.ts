export interface Show {
  id: string;
  title: string;
  description: string;
  imageUrl: string;
  duration: number; // in minutes
  genre: string;
  director: string;
  cast: string[];
  showTimes: ShowTime[];
  ticketPrice: number;
  availableTickets: number;
  venue: string;
  rating: string; // e.g., "PG-13", "R", "G"
}

export interface ShowTime {
  id: string;
  showId: string;
  date: string; // ISO date string
  time: string; // HH:MM format
  availableSeats: number;
  totalSeats: number;
}

export interface Booking {
  id: string;
  showId: string;
  showTimeId: string;
  userId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  numberOfTickets: number;
  totalPrice: number;
  bookingDate: string; // ISO date string
  status: 'confirmed' | 'cancelled' | 'pending';
  qrCode?: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  phone?: string;
}

export interface CalendarEvent {
  id: string;
  title: string;
  date: string;
  time: string;
  showId: string;
}
