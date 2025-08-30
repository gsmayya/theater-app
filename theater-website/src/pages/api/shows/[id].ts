import type { NextApiRequest, NextApiResponse } from 'next';
import { mockShows } from '../../../lib/mockData';
import { Show } from '../../../types/show';

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<Show | { error: string }>
) {
  const { id } = req.query;

  if (req.method === 'GET') {
    const show = mockShows.find(s => s.id === id);
    
    if (!show) {
      return res.status(404).json({ error: 'Show not found' });
    }

    // Simulate API delay
    setTimeout(() => {
      res.status(200).json(show);
    }, 200);
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).json({ error: `Method ${req.method} Not Allowed` });
  }
}
