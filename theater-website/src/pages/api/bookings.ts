import type { NextApiRequest, NextApiResponse } from 'next';
import { mockBookings, mockShows, generateId } from '../../lib/mockData';
import { Booking, BookingRequest } from '../../types/show';
import QRCode from 'qrcode';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Booking | Booking[] | { error: string }>
) {
  if (req.method === 'POST') {
    try {
      const bookingData: BookingRequest = req.body;
      
      // Find the show and show time
      const show = mockShows.find(s => s.id === bookingData.show_id);
      if (!show) {
        return res.status(404).json({ error: 'Show not found' });
      }

      // For now, use the first show time since the BookingRequest doesn't have showTimeId
      const showTime = show.showTimes[0];
      if (!showTime) {
        return res.status(404).json({ error: 'Show time not found' });
      }

      // Check availability
      const availableSeats = showTime.availableSeats || (show.total_tickets - show.booked_tickets);
      if (availableSeats < bookingData.number_of_tickets) {
        return res.status(400).json({ error: 'Not enough available seats' });
      }

      // Generate booking
      const bookingId = generateId();
      const totalPrice = (show.price / 100) * bookingData.number_of_tickets;
      
      // Generate QR code data
      const qrData = JSON.stringify({
        bookingId,
        showTitle: show.title || show.name,
        showDate: showTime.date,
        showTime: showTime.time,
        venue: show.location,
        numberOfTickets: bookingData.number_of_tickets,
        customerName: bookingData.customer_name
      });

      const qrCode = await QRCode.toDataURL(qrData);

      const newBooking: Booking = {
        booking_id: bookingId,
        show_id: bookingData.show_id,
        contact_type: bookingData.contact_type,
        contact_value: bookingData.contact_value,
        number_of_tickets: bookingData.number_of_tickets,
        customer_name: bookingData.customer_name,
        total_amount: Math.round(totalPrice * 100), // Convert to cents
        booking_date: new Date().toISOString(),
        status: 'confirmed',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      };

      // Update available seats (in real app, this would be in database)
      if (showTime.availableSeats !== undefined) {
        showTime.availableSeats -= bookingData.number_of_tickets;
      }
      
      // Add to mock bookings
      mockBookings.push(newBooking);

      // Simulate API delay
      setTimeout(() => {
        res.status(201).json(newBooking);
      }, 500);

    } catch (error) {
      console.error('Booking error:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  } else if (req.method === 'GET') {
    // Get all bookings (in real app, this would be filtered by user)
    res.status(200).json(mockBookings);
  } else {
    res.setHeader('Allow', ['GET', 'POST']);
    res.status(405).json({ error: `Method ${req.method} Not Allowed` });
  }
}
