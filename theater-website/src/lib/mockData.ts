import { Show, ShowTime, Booking, User } from '../types/show';

export const mockShows: Show[] = [
  {
    id: '1',
    title: 'Romeo and Juliet',
    description: 'Shakespeare\'s timeless tragedy of star-crossed lovers. A tale of passion, family rivalry, and destiny that continues to captivate audiences worldwide.',
    imageUrl: '/images/romeo-juliet.jpg',
    duration: 150,
    genre: 'Drama',
    director: 'Sarah Mitchell',
    cast: ['Emma Watson', 'Tom Holland', 'Benedict Cumberbatch'],
    showTimes: [
      {
        id: 'st1',
        showId: '1',
        date: '2024-09-15',
        time: '19:00',
        availableSeats: 45,
        totalSeats: 100
      },
      {
        id: 'st2',
        showId: '1',
        date: '2024-09-16',
        time: '14:00',
        availableSeats: 32,
        totalSeats: 100
      },
      {
        id: 'st3',
        showId: '1',
        date: '2024-09-16',
        time: '19:00',
        availableSeats: 67,
        totalSeats: 100
      }
    ],
    ticketPrice: 45,
    availableTickets: 144,
    venue: 'Grand Theater',
    rating: 'PG-13'
  },
  {
    id: '2',
    title: 'The Lion King',
    description: 'Disney\'s beloved musical brings the African savanna to life with stunning costumes, innovative puppetry, and unforgettable music.',
    imageUrl: '/images/lion-king.jpg',
    duration: 165,
    genre: 'Musical',
    director: 'James Rodriguez',
    cast: ['Michael Johnson', 'Lupita Nyong\'o', 'Idris Elba'],
    showTimes: [
      {
        id: 'st4',
        showId: '2',
        date: '2024-09-20',
        time: '15:00',
        availableSeats: 89,
        totalSeats: 120
      },
      {
        id: 'st5',
        showId: '2',
        date: '2024-09-21',
        time: '19:30',
        availableSeats: 112,
        totalSeats: 120
      }
    ],
    ticketPrice: 65,
    availableTickets: 201,
    venue: 'Royal Opera House',
    rating: 'G'
  },
  {
    id: '3',
    title: 'Hamilton',
    description: 'Lin-Manuel Miranda\'s revolutionary musical biography of Alexander Hamilton, featuring hip-hop, R&B, and traditional show tunes.',
    imageUrl: '/images/hamilton.jpg',
    duration: 175,
    genre: 'Musical',
    director: 'Thomas Kail',
    cast: ['Lin-Manuel Miranda', 'Daveed Diggs', 'Phillipa Soo'],
    showTimes: [
      {
        id: 'st6',
        showId: '3',
        date: '2024-09-25',
        time: '20:00',
        availableSeats: 23,
        totalSeats: 150
      },
      {
        id: 'st7',
        showId: '3',
        date: '2024-09-26',
        time: '20:00',
        availableSeats: 8,
        totalSeats: 150
      }
    ],
    ticketPrice: 95,
    availableTickets: 31,
    venue: 'Broadway Theater',
    rating: 'PG-13'
  },
  {
    id: '4',
    title: 'A Midsummer Night\'s Dream',
    description: 'Shakespeare\'s magical comedy filled with fairies, lovers, and mischief in an enchanted forest setting.',
    imageUrl: '/images/midsummer.jpg',
    duration: 135,
    genre: 'Comedy',
    director: 'Helena Price',
    cast: ['Keira Knightley', 'Oscar Isaac', 'Tilda Swinton'],
    showTimes: [
      {
        id: 'st8',
        showId: '4',
        date: '2024-09-30',
        time: '18:00',
        availableSeats: 78,
        totalSeats: 90
      },
      {
        id: 'st9',
        showId: '4',
        date: '2024-10-01',
        time: '18:00',
        availableSeats: 85,
        totalSeats: 90
      }
    ],
    ticketPrice: 40,
    availableTickets: 163,
    venue: 'Garden Theater',
    rating: 'PG'
  }
];

export const mockUsers: User[] = [
  {
    id: 'user1',
    email: 'john.doe@example.com',
    name: 'John Doe',
    phone: '+1234567890'
  },
  {
    id: 'user2',
    email: 'jane.smith@example.com',
    name: 'Jane Smith',
    phone: '+1987654321'
  }
];

export const mockBookings: Booking[] = [
  {
    id: 'booking1',
    showId: '1',
    showTimeId: 'st1',
    userId: 'user1',
    customerName: 'John Doe',
    customerEmail: 'john.doe@example.com',
    customerPhone: '+1234567890',
    numberOfTickets: 2,
    totalPrice: 90,
    bookingDate: '2024-08-30T10:00:00Z',
    status: 'confirmed',
    qrCode: 'QR_CODE_DATA_1'
  },
  {
    id: 'booking2',
    showId: '2',
    showTimeId: 'st4',
    userId: 'user2',
    customerName: 'Jane Smith',
    customerEmail: 'jane.smith@example.com',
    customerPhone: '+1987654321',
    numberOfTickets: 4,
    totalPrice: 260,
    bookingDate: '2024-08-30T14:30:00Z',
    status: 'confirmed',
    qrCode: 'QR_CODE_DATA_2'
  }
];

// Helper function to generate unique IDs
export const generateId = (): string => {
  return Math.random().toString(36).substr(2, 9);
};
