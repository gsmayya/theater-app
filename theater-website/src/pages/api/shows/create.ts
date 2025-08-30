import type { NextApiRequest, NextApiResponse } from 'next';
import { mockShows, generateId } from '../../../lib/mockData';
import { apiService, CreateShowRequest } from '../../../lib/apiService';
import { Show } from '../../../types/show';

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

      if (apiService.isBackendConfigured) {
        // Try to create show via backend
        const newShow = await apiService.createShow({
          title,
          description,
          location,
          date
        });
        res.status(201).json(newShow);
      } else {
        // Fallback to mock data creation
        const newShow: Show = {
          id: generateId(),
          title,
          description,
          imageUrl: '/images/default-show.jpg',
          duration: 120, // Default duration
          genre: 'Drama', // Default genre
          director: 'Unknown Director',
          cast: ['TBD'],
          showTimes: [
            {
              id: generateId(),
              showId: generateId(),
              date: date,
              time: '19:00', // Default time
              availableSeats: 100,
              totalSeats: 100
            }
          ],
          ticketPrice: 50, // Default price
          availableTickets: 100,
          venue: location,
          rating: 'PG'
        };

        // Add to mock data array (in a real app, this would be persisted)
        mockShows.push(newShow);

        // Simulate API delay
        setTimeout(() => {
          res.status(201).json(newShow);
        }, 300);
      }
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
