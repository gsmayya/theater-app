import { 
  Show, 
  Booking, 
  ApiResponse, 
  PaginatedResponse, 
  SearchParams, 
  BookingRequest, 
  ShowStats,
  ApiError,
  AppConfig
} from '@/types/show';
import { mockShows, mockBookings } from './mockData';

class ApiService {
  private baseUrl: string;
  private apiKey: string | null;
  private enableMockData: boolean;
  private retryAttempts: number = 3;
  private retryDelay: number = 1000;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    this.apiKey = process.env.NEXT_PUBLIC_API_KEY || null;
    this.enableMockData = process.env.NEXT_PUBLIC_ENABLE_MOCK_DATA === 'true' || !this.baseUrl;
  }

  private get headers(): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };
    
    if (this.apiKey) {
      headers['Authorization'] = `Bearer ${this.apiKey}`;
    }
    
    return headers;
  }

  private async sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  private async fetchWithRetry(
    url: string, 
    options: RequestInit = {}, 
    attempt: number = 1
  ): Promise<Response> {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          ...this.headers,
          ...(options.headers || {})
        }
      });

      if (!response.ok) {
        if (response.status >= 500 && attempt < this.retryAttempts) {
          await this.sleep(this.retryDelay * attempt);
          return this.fetchWithRetry(url, options, attempt + 1);
        }
        
        const errorText = await response.text();
        let errorData: ApiError;
        
        try {
          errorData = JSON.parse(errorText);
        } catch {
          errorData = {
            message: errorText || `HTTP ${response.status}`,
            code: `HTTP_${response.status}`,
            details: { status: response.status, statusText: response.statusText }
          };
        }
        
        throw new Error(JSON.stringify(errorData));
      }

      return response;
    } catch (error) {
      if (attempt < this.retryAttempts && error instanceof TypeError) {
        await this.sleep(this.retryDelay * attempt);
        return this.fetchWithRetry(url, options, attempt + 1);
      }
      throw error;
    }
  }

  private async fetchFromBackend<T = any>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    if (this.enableMockData) {
      throw new Error('Backend not available, using mock data');
    }

    const url = `${this.baseUrl}/api/v1${endpoint}`;
    const response = await this.fetchWithRetry(url, options);
    return response.json();
  }

  private handleApiError(error: unknown): ApiError {
    if (error instanceof Error) {
      try {
        return JSON.parse(error.message);
      } catch {
        return {
          message: error.message,
          code: 'UNKNOWN_ERROR',
          details: { originalError: error.message }
        };
      }
    }
    
    return {
      message: 'An unknown error occurred',
      code: 'UNKNOWN_ERROR',
      details: { error }
    };
  }

  public get isBackendConfigured(): boolean {
    return !this.enableMockData && !!this.baseUrl;
  }

  public get config(): AppConfig {
    return {
      apiUrl: this.baseUrl,
      isProduction: process.env.NODE_ENV === 'production',
      enableMockData: this.enableMockData,
      googleClientId: process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID,
      firebaseConfig: process.env.NEXT_PUBLIC_FIREBASE_CONFIG ? 
        JSON.parse(process.env.NEXT_PUBLIC_FIREBASE_CONFIG) : undefined
    };
  }

  // Health check
  public async healthCheck(): Promise<ApiResponse> {
    try {
      return await this.fetchFromBackend('/health');
    } catch (error) {
      return {
        success: false,
        message: 'Backend is not available',
        error: this.handleApiError(error).message,
        status_code: 503
      };
    }
  }

  // Get all shows
  public async getShows(): Promise<Show[]> {
    try {
      const response = await this.fetchFromBackend<Show[]>('/shows/get');
      return response.data || [];
    } catch (error) {
      if (this.enableMockData) {
        return mockShows;
      }
      throw this.handleApiError(error);
    }
  }

  // Get show by ID
  public async getShowById(id: string): Promise<Show | null> {
    try {
      const response = await this.fetchFromBackend<Show>(`/shows/get?id=${id}`);
      return response.data || null;
    } catch (error) {
      if (this.enableMockData) {
        return mockShows.find(show => show.id === id) || null;
      }
      throw this.handleApiError(error);
    }
  }

  // Search shows
  public async searchShows(params: SearchParams): Promise<PaginatedResponse<Show>> {
    try {
      const searchParams = new URLSearchParams();
      
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          searchParams.append(key, value.toString());
        }
      });

      const queryString = searchParams.toString();
      const endpoint = `/search${queryString ? `?${queryString}` : ''}`;
      
      const response = await this.fetchFromBackend<PaginatedResponse<Show>>(endpoint);
      return response.data || { success: true, message: 'No results', data: [], pagination: { page: 1, page_size: 10, total: 0, total_pages: 0 } };
    } catch (error) {
      if (this.enableMockData) {
        let filteredShows = mockShows;
        
        if (params.search) {
          filteredShows = filteredShows.filter(show => 
            show.name.toLowerCase().includes(params.search!.toLowerCase()) ||
            show.details.toLowerCase().includes(params.search!.toLowerCase())
          );
        }
        
        if (params.location) {
          filteredShows = filteredShows.filter(show => 
            show.location.toLowerCase().includes(params.location!.toLowerCase())
          );
        }
        
        if (params.min_price !== undefined) {
          filteredShows = filteredShows.filter(show => show.price >= params.min_price!);
        }
        
        if (params.max_price !== undefined) {
          filteredShows = filteredShows.filter(show => show.price <= params.max_price!);
        }
        
        if (params.only_available) {
          filteredShows = filteredShows.filter(show => 
            (show.total_tickets - show.booked_tickets) > 0
          );
        }

        return {
          success: true,
          message: 'Mock data results',
          data: filteredShows,
          pagination: {
            page: params.page || 1,
            page_size: params.page_size || 10,
            total: filteredShows.length,
            total_pages: Math.ceil(filteredShows.length / (params.page_size || 10))
          }
        };
      }
      throw this.handleApiError(error);
    }
  }

  // Create booking
  public async createBooking(bookingData: BookingRequest): Promise<Booking> {
    try {
      const response = await this.fetchFromBackend<Booking>('/bookings/create', {
        method: 'POST',
        body: JSON.stringify(bookingData)
      });
      return response.data!;
    } catch (error) {
      if (this.enableMockData) {
        const mockBooking: Booking = {
          booking_id: `BK-${Date.now()}`,
          show_id: bookingData.show_id,
          contact_type: bookingData.contact_type,
          contact_value: bookingData.contact_value,
          number_of_tickets: bookingData.number_of_tickets,
          customer_name: bookingData.customer_name,
          total_amount: 0, // Would be calculated by backend
          booking_date: bookingData.booking_date,
          status: 'pending',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        };
        return mockBooking;
      }
      throw this.handleApiError(error);
    }
  }

  // Get booking by ID
  public async getBookingById(bookingId: string): Promise<Booking | null> {
    try {
      const response = await this.fetchFromBackend<Booking>(`/bookings/get?booking_id=${bookingId}`);
      return response.data || null;
    } catch (error) {
      if (this.enableMockData) {
        return mockBookings.find(booking => booking.booking_id === bookingId) || null;
      }
      throw this.handleApiError(error);
    }
  }

  // Get bookings by contact
  public async getBookingsByContact(contactType: string, contactValue: string): Promise<Booking[]> {
    try {
      const response = await this.fetchFromBackend<Booking[]>(`/bookings/by-contact?contact_type=${contactType}&contact_value=${contactValue}`);
      return response.data || [];
    } catch (error) {
      if (this.enableMockData) {
        return mockBookings.filter(booking => 
          booking.contact_type === contactType && booking.contact_value === contactValue
        );
      }
      throw this.handleApiError(error);
    }
  }

  // Get show statistics
  public async getShowStats(showId: string): Promise<ShowStats | null> {
    try {
      const response = await this.fetchFromBackend<ShowStats>(`/shows/booking-summary?show_id=${showId}`);
      return response.data || null;
    } catch (error) {
      if (this.enableMockData) {
        const show = mockShows.find(s => s.id === showId);
        const bookings = mockBookings.filter(b => b.show_id === showId);
        
        if (!show) return null;
        
        return {
          show_id: showId,
          name: show.name,
          total_bookings: bookings.length,
          confirmed_revenue: bookings
            .filter(b => b.status === 'confirmed')
            .reduce((sum, b) => sum + b.total_amount, 0),
          pending_revenue: bookings
            .filter(b => b.status === 'pending')
            .reduce((sum, b) => sum + b.total_amount, 0),
          confirmed_bookings: bookings.filter(b => b.status === 'confirmed').length,
          pending_bookings: bookings.filter(b => b.status === 'pending').length,
          cancelled_bookings: bookings.filter(b => b.status === 'cancelled').length,
          available_tickets: show.total_tickets - show.booked_tickets
        };
      }
      throw this.handleApiError(error);
    }
  }

  // Update booking status
  public async updateBookingStatus(bookingId: string, status: 'confirmed' | 'cancelled'): Promise<Booking> {
    try {
      const response = await this.fetchFromBackend<Booking>('/bookings/update-status', {
        method: 'PUT',
        body: JSON.stringify({ booking_id: bookingId, status })
      });
      return response.data!;
    } catch (error) {
      if (this.enableMockData) {
        const booking = mockBookings.find(b => b.booking_id === bookingId);
        if (booking) {
          booking.status = status;
          booking.updated_at = new Date().toISOString();
        }
        return booking!;
      }
      throw this.handleApiError(error);
    }
  }
}

export const apiService = new ApiService();
