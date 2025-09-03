import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Layout from '../components/Layout';
import ShowCard from '../components/ShowCard';
import Calendar from '../components/Calendar';
import SearchBar from '../components/SearchBar';
import { Show, CalendarEvent } from '../types/show';

const HomePage = () => {
  const [shows, setShows] = useState<Show[]>([]);
  const [searchResults, setSearchResults] = useState<Show[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedEvents, setSelectedEvents] = useState<CalendarEvent[]>([]);
  const router = useRouter();

  useEffect(() => {
    fetchShows();
  }, []);

  const fetchShows = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/shows');
      if (!response.ok) {
        throw new Error('Failed to fetch shows');
      }
      const showsData = await response.json();
      setShows(showsData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleBookShow = (show: Show) => {
    router.push('/booking');
  };

  const handleDateSelect = (date: string, events: CalendarEvent[]) => {
    setSelectedDate(date);
    setSelectedEvents(events);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const formatTime = (timeString: string) => {
    return new Date(`2000-01-01T${timeString}`).toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    });
  };

  if (loading) {
    return (
      <Layout title="Current Shows - Illuminating Windows">
        <div className="flex justify-center items-center min-h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-theater-primary"></div>
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout title="Current Shows - Illuminating Windows">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            <strong className="font-bold">Error: </strong>
            <span className="block sm:inline">{error}</span>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout title="Current Shows - Illuminating Windows">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-theater-primary to-theater-secondary text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold mb-4">
              Welcome to Illuminating Windows
            </h1>
            <p className="text-xl md:text-2xl mb-8 opacity-90">
              Discover amazing live performances and book your tickets today
            </p>
            <button
              onClick={() => router.push('/booking')}
              className="bg-white text-theater-primary px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-colors duration-200"
            >
              Book Tickets Now
            </button>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 xl:grid-cols-4 gap-8">
          {/* Main Content */}
          <div className="xl:col-span-3">
            {/* Shows Section */}
            <div className="mb-8">
              <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-6">
                <h2 className="text-3xl font-bold text-theater-dark">
                  Current Shows
                </h2>
                <div className="lg:w-1/2">
                  <SearchBar onSearch={setSearchResults} />
                </div>
              </div>
              
              {((searchResults && searchResults.length === 0) || (!searchResults && shows.length === 0)) ? (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">🎭</div>
                  <h3 className="text-xl font-semibold text-gray-600 mb-2">
                    No shows available
                  </h3>
                  <p className="text-gray-500">
                    Check back soon for exciting new performances!
                  </p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {(searchResults ?? shows).map(show => (
                    <ShowCard
                      key={show.id}
                      show={show}
                      onBookShow={handleBookShow}
                    />
                  ))}
                </div>
              )}
            </div>

            {/* Selected Date Events */}
            {selectedDate && selectedEvents.length > 0 && (
              <div className="bg-white rounded-lg shadow-lg p-6">
                <h3 className="text-xl font-bold text-theater-dark mb-4">
                  Shows on {formatDate(selectedDate)}
                </h3>
                <div className="space-y-4">
                  {selectedEvents.map(event => {
                    const show = shows.find(s => s.id === event.showId);
                    return (
                      <div
                        key={event.id}
                        className="flex items-center justify-between p-4 bg-theater-light rounded-lg"
                      >
                        <div>
                          <h4 className="font-semibold text-theater-dark">
                            {event.title}
                          </h4>
                          <p className="text-sm text-gray-600">
                            {formatTime(event.time)} • {show?.venue}
                          </p>
                        </div>
                        <button
                          onClick={() => router.push('/booking')}
                          className="bg-theater-primary text-white px-4 py-2 rounded-lg hover:bg-theater-primary/90 transition-colors"
                        >
                          Book Now
                        </button>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </div>

          {/* Sidebar */}
          <div className="xl:col-span-1">
            <div className="sticky top-4">
              <h3 className="text-xl font-bold text-theater-dark mb-4">
                Show Calendar
              </h3>
              <Calendar
                shows={shows}
                onDateSelect={handleDateSelect}
              />
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default HomePage;
