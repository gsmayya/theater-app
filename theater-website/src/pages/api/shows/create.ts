import type { NextApiRequest, NextApiResponse } from 'next';
import { mockShows, generateId } from '../../../lib/mockData';
import { apiService } from '../../../lib/apiService';
import { Show } from '../../../types/show';

interface CreateShowRequest {
  title: string;
  description: string;
  location: string;
  date: string;
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Show | { error: string }>
) {
  if (req.method === 'POST') {
    try {
      const { title, description, location, date } = req.body as CreateShowRequest;

      // Validate required fields
      if (!title || !description || !location || !date) {
        return res.status(400).json({ 
          error: 'Missing required fields: title, description, location, date' 
        });
      }

      // Create show using mock data (backend integration can be added later)
        const newShow: Show = {
          id: generateId(),
          name: title,
          title,
          details: description,
          price: 5000, // $50.00 in cents
          total_tickets: 100,
          booked_tickets: 0,
          location,
          show_number: `SH${Date.now()}`,
          show_date: date,
          showTimes: [
            {
              id: generateId(),
              show_id: generateId(),
              date: date,
              time: '19:00', // Default time
              availableSeats: 100,
              totalSeats: 100
            }
          ],
          images: ['/images/default-show.jpg'],
          videos: [],
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        };

      // Add to mock data array (in a real app, this would be persisted)
      mockShows.push(newShow);

      // Simulate API delay
      setTimeout(() => {
        res.status(201).json(newShow);
      }, 300);
    } catch (error) {
      console.error('Create show API failed:', error);
      
      if (error instanceof Error && error.message.includes('Backend API error')) {
        res.status(500).json({ error: `Backend error: ${error.message}` });
      } else {
        res.status(500).json({ error: 'Internal server error' });
      }
    }
  } else {
    res.setHeader('Allow', ['POST']);
    res.status(405).json({ error: `Method ${req.method} Not Allowed` });
  }
}
