import type { NextApiRequest, NextApiResponse } from 'next';
import { mockShows } from '../../lib/mockData';
import { apiService } from '../../lib/apiService';
import { Show } from '../../types/show';

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Show[] | { error: string }>
) {
  if (req.method === 'GET') {
    try {
      if (apiService.isBackendConfigured) {
        // Try to fetch from backend
        const shows = await apiService.getShows();
        res.status(200).json(shows);
      } else {
        // Fallback to mock data
        setTimeout(() => {
          res.status(200).json(mockShows);
        }, 300);
      }
    } catch (error) {
      console.error('Backend API failed, falling back to mock data:', error);
      // Fallback to mock data on error
      setTimeout(() => {
        res.status(200).json(mockShows);
      }, 300);
    }
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).json({ error: `Method ${req.method} Not Allowed` });
  }
}
