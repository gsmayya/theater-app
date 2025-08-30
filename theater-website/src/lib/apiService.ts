import { Show, Booking } from '../types/show';

export interface CreateShowRequest {
  title: string;
  description: string;
  location: string;
  date: string;
}

export interface SearchShowsParams {
  location?: string;
  title?: string;
  limit?: number;
  offset?: number;
}

class ApiService {
  private baseUrl: string;
  private apiKey: string | null;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '';
    this.apiKey = process.env.API_KEY || null;
  }

  private get headers() {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    
    if (this.apiKey) {
      headers['Authorization'] = `Bearer ${this.apiKey}`;
      // Or if your API uses a different auth header format:
      // headers['X-API-Key'] = this.apiKey;
    }
    
    return headers;
  }

  private async fetchFromBackend(endpoint: string, options: RequestInit = {}): Promise<Response> {
    if (!this.baseUrl) {
      throw new Error('Backend not configured');
    }

    const url = `${this.baseUrl}/api/v1${endpoint}`;
    
    const response = await fetch(url, {
      ...options,
      headers: {
        ...this.headers,
        ...(options.headers || {})
      }
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Backend API error: ${response.status} ${errorText}`);
    }

    return response;
  }

  public get isBackendConfigured(): boolean {
    return !!this.baseUrl;
  }

  // Search shows by location and/or title
  public async searchShows(params: SearchShowsParams): Promise<Show[]> {
    const searchParams = new URLSearchParams();
    
    if (params.location) searchParams.append('location', params.location);
    if (params.title) searchParams.append('title', params.title);
    if (params.limit) searchParams.append('limit', params.limit.toString());
    if (params.offset) searchParams.append('offset', params.offset.toString());

    const queryString = searchParams.toString();
    const endpoint = `/search${queryString ? `?${queryString}` : ''}`;
    
    const response = await this.fetchFromBackend(endpoint);
    return response.json();
  }

  // Get all shows
  public async getShows(): Promise<Show[]> {
    const response = await this.fetchFromBackend('/get');
    return response.json();
  }

  // Create a new show
  public async createShow(showData: CreateShowRequest): Promise<Show> {
    const response = await this.fetchFromBackend('/shows', {
      method: 'POST',
      body: JSON.stringify(showData)
    });
    return response.json();
  }

  // Book tickets (if your backend has this endpoint)
  public async bookTickets(bookingData: any): Promise<Booking> {
    const response = await this.fetchFromBackend('/bookings', {
      method: 'POST',
      body: JSON.stringify(bookingData)
    });
    return response.json();
  }
}

export const apiService = new ApiService();
