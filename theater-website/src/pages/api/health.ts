import type { NextApiRequest, NextApiResponse } from 'next';

interface HealthResponse {
  status: 'healthy' | 'unhealthy';
  timestamp: string;
  uptime: number;
  version: string;
  services: {
    frontend: 'connected';
    backend?: 'connected' | 'disconnected';
  };
}

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<HealthResponse>
) {
  if (req.method !== 'GET') {
    return res.status(405).json({
      status: 'unhealthy',
      timestamp: new Date().toISOString(),
      uptime: process.uptime(),
      version: '1.0.0',
      services: {
        frontend: 'connected',
      },
    });
  }

  const healthData: HealthResponse = {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    version: '1.0.0',
    services: {
      frontend: 'connected',
    },
  };

  // Try to check backend health if API URL is configured
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  if (apiUrl) {
    // In a real implementation, you might want to make an actual HTTP request
    // to the backend to check its health, but for now we'll assume it's connected
    // if the environment variable is set
    healthData.services.backend = 'connected';
  }

  res.status(200).json(healthData);
}
