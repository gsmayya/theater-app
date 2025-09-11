import { AppConfig } from '@/types/show';

// Environment validation
function validateEnvVar(name: string, value: string | undefined, required: boolean = false): string {
  if (required && !value) {
    throw new Error(`Required environment variable ${name} is not set`);
  }
  return value || '';
}

function validateUrl(name: string, value: string): string {
  try {
    new URL(value);
    return value;
  } catch {
    throw new Error(`Invalid URL for environment variable ${name}: ${value}`);
  }
}

function validateJson(name: string, value: string | undefined): any {
  if (!value) return undefined;
  try {
    return JSON.parse(value);
  } catch {
    throw new Error(`Invalid JSON for environment variable ${name}: ${value}`);
  }
}

// Environment configuration
export const config: AppConfig = {
  apiUrl: validateUrl('NEXT_PUBLIC_API_URL', validateEnvVar('NEXT_PUBLIC_API_URL', process.env.NEXT_PUBLIC_API_URL, true)),
  isProduction: process.env.NODE_ENV === 'production',
  enableMockData: process.env.NEXT_PUBLIC_ENABLE_MOCK_DATA === 'true' || !process.env.NEXT_PUBLIC_API_URL,
  googleClientId: validateEnvVar('NEXT_PUBLIC_GOOGLE_CLIENT_ID', process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID),
  firebaseConfig: validateJson('NEXT_PUBLIC_FIREBASE_CONFIG', process.env.NEXT_PUBLIC_FIREBASE_CONFIG),
};

// Additional app configuration
export const appConfig = {
  name: process.env.NEXT_PUBLIC_APP_NAME || 'Theater Booking System',
  description: process.env.NEXT_PUBLIC_APP_DESCRIPTION || 'Book tickets for amazing theater shows',
  url: process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000',
  version: process.env.npm_package_version || '1.0.0',
  
  // Feature flags
  features: {
    analytics: process.env.NEXT_PUBLIC_ENABLE_ANALYTICS === 'true',
    debug: process.env.NEXT_PUBLIC_ENABLE_DEBUG === 'true',
    pwa: process.env.NEXT_PUBLIC_ENABLE_PWA === 'true',
  },
  
  // Analytics
  analytics: {
    gaTrackingId: process.env.NEXT_PUBLIC_GA_TRACKING_ID,
    gtmId: process.env.NEXT_PUBLIC_GTM_ID,
  },
  
  // Security
  security: {
    nextAuthSecret: process.env.NEXTAUTH_SECRET,
    nextAuthUrl: process.env.NEXTAUTH_URL || 'http://localhost:3000',
  },
  
  // API settings
  api: {
    timeout: 10000, // 10 seconds
    retryAttempts: 3,
    retryDelay: 1000, // 1 second
  },
  
  // UI settings
  ui: {
    itemsPerPage: 12,
    maxSearchResults: 100,
    debounceDelay: 300, // milliseconds
  },
};

// Validation function to check if all required config is present
export function validateConfig(): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];
  
  try {
    // Check API URL
    if (!config.apiUrl) {
      errors.push('API URL is required');
    }
    
    // Check Google Client ID if authentication is enabled
    if (process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID && !config.googleClientId) {
      errors.push('Google Client ID is invalid');
    }
    
    // Check Firebase config if provided
    if (process.env.NEXT_PUBLIC_FIREBASE_CONFIG && !config.firebaseConfig) {
      errors.push('Firebase configuration is invalid');
    }
    
    // Check NextAuth secret in production
    if (config.isProduction && !appConfig.security.nextAuthSecret) {
      errors.push('NextAuth secret is required in production');
    }
    
    return {
      isValid: errors.length === 0,
      errors
    };
  } catch (error) {
    return {
      isValid: false,
      errors: [error instanceof Error ? error.message : 'Unknown configuration error']
    };
  }
}

// Development helpers
export const isDevelopment = process.env.NODE_ENV === 'development';
export const isProduction = process.env.NODE_ENV === 'production';
export const isTest = process.env.NODE_ENV === 'test';

// Debug logging
export function debugLog(message: string, data?: any): void {
  if (appConfig.features.debug || isDevelopment) {
    console.log(`[DEBUG] ${message}`, data || '');
  }
}

// Error logging
export function errorLog(message: string, error?: Error): void {
  console.error(`[ERROR] ${message}`, error || '');
  
  // In production, you might want to send this to an error tracking service
  if (isProduction && appConfig.features.analytics) {
    // Example: Send to error tracking service
    // errorTrackingService.captureException(error);
  }
}
