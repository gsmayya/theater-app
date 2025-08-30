import React, { useState } from 'react';
import { Show } from '../types/show';

interface SearchBarProps {
  onSearch: (results: Show[]) => void;
  loading?: boolean;
  className?: string;
}

const SearchBar: React.FC<SearchBarProps> = ({ onSearch, loading = false, className = '' }) => {
  const [location, setLocation] = useState('');
  const [title, setTitle] = useState('');
  const [isSearching, setIsSearching] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!location.trim() && !title.trim()) {
      return;
    }

    setIsSearching(true);

    try {
      const searchParams = new URLSearchParams();
      if (location.trim()) searchParams.append('location', location.trim());
      if (title.trim()) searchParams.append('title', title.trim());

      const response = await fetch(`/api/search?${searchParams.toString()}`);
      if (!response.ok) {
        throw new Error('Search failed');
      }

      const results = await response.json();
      onSearch(results);
    } catch (error) {
      console.error('Search error:', error);
      // You could show an error message to the user here
    } finally {
      setIsSearching(false);
    }
  };

  const handleClear = () => {
    setLocation('');
    setTitle('');
    onSearch([]); // Clear search results
  };

  return (
    <div className={`bg-white rounded-lg shadow-lg p-6 ${className}`}>
      <h3 className="text-lg font-semibold text-theater-dark mb-4">
        Search Shows
      </h3>
      
      <form onSubmit={handleSearch} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Location/Venue
            </label>
            <input
              type="text"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              placeholder="e.g., Grand Theater, Broadway"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Show Title
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="e.g., Romeo and Juliet, Hamilton"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
            />
          </div>
        </div>
        
        <div className="flex gap-2">
          <button
            type="submit"
            disabled={isSearching || loading || (!location.trim() && !title.trim())}
            className="flex-1 bg-theater-primary text-white py-2 px-4 rounded-md hover:bg-theater-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center"
          >
            {isSearching ? (
              <>
                <svg className="animate-spin -ml-1 mr-3 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Searching...
              </>
            ) : (
              <>
                <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                Search
              </>
            )}
          </button>
          
          <button
            type="button"
            onClick={handleClear}
            disabled={isSearching || loading}
            className="bg-gray-100 text-gray-600 py-2 px-4 rounded-md hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Clear
          </button>
        </div>
      </form>
    </div>
  );
};

export default SearchBar;
