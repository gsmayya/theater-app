import type { NextApiRequest, NextApiResponse } from 'next';
import { mockShows } from '../../lib/mockData';
import { Show } from '../../types/show';

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<Show[] | { error: string }>
) {
  if (req.method === 'GET') {
    // Simulate API delay
    setTimeout(() => {
      res.status(200).json(mockShows);
    }, 300);
  } else {
    res.setHeader('Allow', ['GET']);
    res.status(405).json({ error: `Method ${req.method} Not Allowed` });
  }
}
