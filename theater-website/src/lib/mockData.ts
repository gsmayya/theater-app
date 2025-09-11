import { Show, ShowTime, Booking, User } from '../types/show';

export const mockShows: Show[] = [
  {
    id: '1',
    name: 'Romeo and Juliet',
    title: 'Romeo and Juliet',
    details: 'Shakespeare\'s timeless tragedy of star-crossed lovers. A tale of passion, family rivalry, and destiny that continues to captivate audiences worldwide.',
    price: 4500, // $45.00 in cents
    total_tickets: 100,
    booked_tickets: 55,
    location: 'Grand Theater',
    show_number: 'RJ001',
    show_date: '2024-09-15',
    showTimes: [
      {
        id: 'st1',
        show_id: '1',
        date: '2024-09-15',
        time: '19:00',
        availableSeats: 45,
        totalSeats: 100
      },
      {
        id: 'st2',
        show_id: '1',
        date: '2024-09-16',
        time: '14:00',
        availableSeats: 32,
        totalSeats: 100
      },
      {
        id: 'st3',
        show_id: '1',
        date: '2024-09-16',
        time: '19:00',
        availableSeats: 67,
        totalSeats: 100
      }
    ],
    images: ['/images/romeo-juliet.jpg'],
    videos: [],
    created_at: '2024-08-01T10:00:00Z',
    updated_at: '2024-08-01T10:00:00Z'
  },
  {
    id: '2',
    name: 'The Lion King',
    title: 'The Lion King',
    details: 'Disney\'s beloved musical brings the African savanna to life with stunning costumes, innovative puppetry, and unforgettable music.',
    price: 6500, // $65.00 in cents
    total_tickets: 120,
    booked_tickets: 8,
    location: 'Royal Opera House',
    show_number: 'LK002',
    show_date: '2024-09-20',
    showTimes: [
      {
        id: 'st4',
        show_id: '2',
        date: '2024-09-20',
        time: '15:00',
        availableSeats: 89,
        totalSeats: 120
      },
      {
        id: 'st5',
        show_id: '2',
        date: '2024-09-21',
        time: '19:30',
        availableSeats: 112,
        totalSeats: 120
      }
    ],
    images: ['/images/lion-king.jpg'],
    videos: [],
    created_at: '2024-08-01T10:00:00Z',
    updated_at: '2024-08-01T10:00:00Z'
  },
  {
    id: '3',
    name: 'Hamilton',
    title: 'Hamilton',
    details: 'Lin-Manuel Miranda\'s revolutionary musical biography of Alexander Hamilton, featuring hip-hop, R&B, and traditional show tunes.',
    price: 9500, // $95.00 in cents
    total_tickets: 150,
    booked_tickets: 119,
    location: 'Broadway Theater',
    show_number: 'HM003',
    show_date: '2024-09-25',
    showTimes: [
      {
        id: 'st6',
        show_id: '3',
        date: '2024-09-25',
        time: '20:00',
        availableSeats: 23,
        totalSeats: 150
      },
      {
        id: 'st7',
        show_id: '3',
        date: '2024-09-26',
        time: '20:00',
        availableSeats: 8,
        totalSeats: 150
      }
    ],
    images: ['/images/hamilton.jpg'],
    videos: [],
    created_at: '2024-08-01T10:00:00Z',
    updated_at: '2024-08-01T10:00:00Z'
  },
  {
    id: '4',
    name: 'A Midsummer Night\'s Dream',
    title: 'A Midsummer Night\'s Dream',
    details: 'Shakespeare\'s magical comedy filled with fairies, lovers, and mischief in an enchanted forest setting.',
    price: 4000, // $40.00 in cents
    total_tickets: 90,
    booked_tickets: 5,
    location: 'Garden Theater',
    show_number: 'MS004',
    show_date: '2024-09-30',
    showTimes: [
      {
        id: 'st8',
        show_id: '4',
        date: '2024-09-30',
        time: '18:00',
        availableSeats: 78,
        totalSeats: 90
      },
      {
        id: 'st9',
        show_id: '4',
        date: '2024-10-01',
        time: '18:00',
        availableSeats: 85,
        totalSeats: 90
      }
    ],
    images: ['/images/midsummer.jpg'],
    videos: [],
    created_at: '2024-08-01T10:00:00Z',
    updated_at: '2024-08-01T10:00:00Z'
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
    booking_id: 'booking1',
    show_id: '1',
    contact_type: 'email',
    contact_value: 'john.doe@example.com',
    number_of_tickets: 2,
    customer_name: 'John Doe',
    total_amount: 9000, // $90.00 in cents
    booking_date: '2024-08-30T10:00:00Z',
    status: 'confirmed',
    created_at: '2024-08-30T10:00:00Z',
    updated_at: '2024-08-30T10:00:00Z'
  },
  {
    booking_id: 'booking2',
    show_id: '2',
    contact_type: 'email',
    contact_value: 'jane.smith@example.com',
    number_of_tickets: 4,
    customer_name: 'Jane Smith',
    total_amount: 26000, // $260.00 in cents
    booking_date: '2024-08-30T14:30:00Z',
    status: 'confirmed',
    created_at: '2024-08-30T14:30:00Z',
    updated_at: '2024-08-30T14:30:00Z'
  }
];

// Helper function to generate unique IDs
export const generateId = (): string => {
  return Math.random().toString(36).substr(2, 9);
};
