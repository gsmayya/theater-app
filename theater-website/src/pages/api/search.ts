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
      const { location, title, limit = 10, offset = 0 } = req.query;

      if (apiService.isBackendConfigured) {
        // Try to search from backend
        const shows = await apiService.searchShows({
          location: location as string,
          search: title as string,
          page_size: parseInt(limit as string),
          page: Math.floor(parseInt(offset as string) / parseInt(limit as string)) + 1
        });
        res.status(200).json(shows.data || []);
      } else {
        // Fallback to filtering mock data
        let filteredShows = mockShows;
        
        if (location) {
          filteredShows = filteredShows.filter(show => 
            show.location.toLowerCase().includes((location as string).toLowerCase())
          );
        }
        
        if (title) {
          filteredShows = filteredShows.filter(show => 
            (show.title || show.name).toLowerCase().includes((title as string).toLowerCase())
          );
        }
        
        // Apply pagination
        const startIndex = parseInt(offset as string) || 0;
        const limitNum = parseInt(limit as string) || 10;
        const paginatedShows = filteredShows.slice(startIndex, startIndex + limitNum);
        
        // Simulate API delay
        setTimeout(() => {
          res.status(200).json(paginatedShows);
        }, 300);
      }
    } catch (error) {
      console.error('Search API failed, falling back to mock data:', error);
      
      // Fallback search on mock data
      const { location, title, limit = 10, offset = 0 } = req.query;
      let filteredShows = mockShows;
      
      if (location) {
        filteredShows = filteredShows.filter(show => 
          show.location.toLowerCase().includes((location as string).toLowerCase())
        );
      }
      
      if (title) {
        filteredShows = filteredShows.filter(show => 
          show.title.toLowerCase().includes((title as string).toLowerCase())
        );
      }
      
      const startIndex = parseInt(offset as string) || 0;
      const limitNum = parseInt(limit as string) || 10;
      const paginatedShows = filteredShows.slice(startIndex, startIndex + limitNum);
      
      setTimeout(() => {
        res.status(200).json(paginatedShows);
      }, 300);
    }
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).json({ error: `Method ${req.method} Not Allowed` });
  }
}
