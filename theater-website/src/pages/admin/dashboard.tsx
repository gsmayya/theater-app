import React, { useState, useEffect } from 'react';
import Layout from '../../components/Layout';
import CreateShowForm from '../../components/CreateShowForm';
import { useAuth } from '../../contexts/AuthContext';
import { useRouter } from 'next/router';
import { Show } from '../../types/show';

const AdminDashboard: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const [shows, setShows] = useState<Show[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/admin');
      return;
    }
    fetchShows();
  }, [isAuthenticated, router]);

  const fetchShows = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/shows');
      if (response.ok) {
        const data = await response.json();
        setShows(data);
      }
    } catch (error) {
      console.error('Failed to fetch shows:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleShowCreated = (newShow: Show) => {
    setShows(prev => [...prev, newShow]);
  };

  if (!isAuthenticated) {
    return (
      <Layout title="Admin Dashboard - Illuminating Windows">
        <div className="flex justify-center items-center min-h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-theater-primary"></div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout title="Admin Dashboard - Illuminating Windows">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-theater-dark">Admin Dashboard</h1>
            <p className="text-gray-600 mt-2">Welcome back, {user?.name}!</p>
          </div>
          <button
            onClick={logout}
            className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors"
          >
            Logout
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Create Show Form */}
          <div>
            <CreateShowForm 
              onShowCreated={handleShowCreated}
              className="sticky top-4"
            />
          </div>

          {/* Shows List */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-xl font-bold text-theater-dark">
                Current Shows ({shows.length})
              </h3>
              <button
                onClick={fetchShows}
                className="text-theater-primary hover:text-theater-primary/80 transition-colors"
                disabled={loading}
              >
                {loading ? (
                  <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                ) : (
                  <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                )}
              </button>
            </div>

            {loading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-theater-primary"></div>
              </div>
            ) : shows.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                <div className="text-4xl mb-2">🎭</div>
                <p>No shows created yet</p>
              </div>
            ) : (
              <div className="space-y-4 max-h-96 overflow-y-auto">
                {shows.map(show => (
                  <div
                    key={show.id}
                    className="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
                  >
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h4 className="font-semibold text-theater-dark mb-1">
                          {show.title}
                        </h4>
                        <p className="text-sm text-gray-600 mb-2 line-clamp-2">
                          {show.details}
                        </p>
                        <div className="flex items-center text-xs text-gray-500 space-x-4">
                          <span>📍 {show.location}</span>
                          <span>💰 ${(show.price / 100).toFixed(2)}</span>
                          <span>🎫 {show.total_tickets - show.booked_tickets} available</span>
                        </div>
                      </div>
                      <div className="flex items-center space-x-2 ml-4">
                        <span className="px-2 py-1 rounded-full text-xs bg-gray-100 text-gray-800">
                          Theater Show
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Stats Cards */}
        <div className="mt-12 grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-gradient-to-r from-theater-primary to-theater-secondary text-white rounded-lg p-6">
            <h4 className="text-lg font-semibold mb-2">Total Shows</h4>
            <p className="text-3xl font-bold">{shows.length}</p>
          </div>
          
          <div className="bg-gradient-to-r from-theater-secondary to-theater-accent text-white rounded-lg p-6">
            <h4 className="text-lg font-semibold mb-2">Available Tickets</h4>
            <p className="text-3xl font-bold">
              {shows.reduce((sum, show) => sum + (show.total_tickets - show.booked_tickets), 0)}
            </p>
          </div>
          
          <div className="bg-gradient-to-r from-theater-accent to-orange-400 text-white rounded-lg p-6">
            <h4 className="text-lg font-semibold mb-2">Average Price</h4>
            <p className="text-3xl font-bold">
              ${shows.length > 0 ? Math.round(shows.reduce((sum, show) => sum + (show.price / 100), 0) / shows.length) : 0}
            </p>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default AdminDashboard;
