import type { NextApiRequest, NextApiResponse } from 'next';
import { mockBookings, mockShows, generateId } from '../../lib/mockData';
import { Booking } from '../../types/show';
import QRCode from 'qrcode';

interface BookingRequest {
  showId: string;
  showTimeId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  numberOfTickets: number;
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Booking | Booking[] | { error: string }>
) {
  if (req.method === 'POST') {
    try {
      const bookingData: BookingRequest = req.body;
      
      // Find the show and show time
      const show = mockShows.find(s => s.id === bookingData.showId);
      if (!show) {
        return res.status(404).json({ error: 'Show not found' });
      }

      const showTime = show.showTimes.find(st => st.id === bookingData.showTimeId);
      if (!showTime) {
        return res.status(404).json({ error: 'Show time not found' });
      }

      // Check availability
      if (showTime.availableSeats < bookingData.numberOfTickets) {
        return res.status(400).json({ error: 'Not enough available seats' });
      }

      // Generate booking
      const bookingId = generateId();
      const totalPrice = show.ticketPrice * bookingData.numberOfTickets;
      
      // Generate QR code data
      const qrData = JSON.stringify({
        bookingId,
        showTitle: show.title,
        showDate: showTime.date,
        showTime: showTime.time,
        venue: show.venue,
        numberOfTickets: bookingData.numberOfTickets,
        customerName: bookingData.customerName
      });

      const qrCode = await QRCode.toDataURL(qrData);

      const newBooking: Booking = {
        id: bookingId,
        showId: bookingData.showId,
        showTimeId: bookingData.showTimeId,
        userId: generateId(), // In real app, this would come from auth
        customerName: bookingData.customerName,
        customerEmail: bookingData.customerEmail,
        customerPhone: bookingData.customerPhone,
        numberOfTickets: bookingData.numberOfTickets,
        totalPrice,
        bookingDate: new Date().toISOString(),
        status: 'confirmed',
        qrCode
      };

      // Update available seats (in real app, this would be in database)
      showTime.availableSeats -= bookingData.numberOfTickets;
      
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
