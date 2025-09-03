import React, { useState } from 'react';
import Layout from '../components/Layout';
import { useAuth } from '../contexts/AuthContext';
import { Show, ShowTime } from '../types/show';

const BookingPage: React.FC = () => {
  const { isAuthenticated, user, login, loginWithGoogle, logout, loading, googleAvailable } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [shows, setShows] = useState<Show[]>([]);
  const [selectedShow, setSelectedShow] = useState<string>('');
  const [selectedShowTime, setSelectedShowTime] = useState<string>('');
  const [tickets, setTickets] = useState<number>(1);
  const [customerName, setCustomerName] = useState('');
  const [customerPhone, setCustomerPhone] = useState('');
  const [bookingResult, setBookingResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  React.useEffect(() => {
    const fetchShows = async () => {
      try {
        const response = await fetch('/api/shows');
        const data = await response.json();
        setShows(data);
      } catch (err) {
        setError('Failed to load shows');
      }
    };
    fetchShows();
  }, []);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    await login(email, password);
  };

  const handleBooking = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      const selectedShowObj = shows.find(s => s.id === selectedShow);
      if (!selectedShowObj) {
        setError('Please select a show');
        return;
      }
      const selectedShowTimeObj = selectedShowObj.showTimes.find(st => st.id === selectedShowTime);
      if (!selectedShowTimeObj) {
        setError('Please select a show time');
        return;
      }

      const response = await fetch('/api/bookings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          showId: selectedShow,
          showTimeId: selectedShowTime,
          customerName: customerName || user?.name || 'Guest',
          customerEmail: user?.email || 'guest@example.com',
          customerPhone: customerPhone || user?.phone || '',
          numberOfTickets: tickets
        })
      });
      
      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.error || 'Booking failed');
      }

      const data = await response.json();
      setBookingResult(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    }
  };

  return (
    <Layout title="Book Tickets - Illuminating Windows">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {!isAuthenticated ? (
          <div className="max-w-md mx-auto bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-2xl font-bold mb-6 text-theater-dark">Login to Book Tickets</h2>
            <form onSubmit={handleLogin} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Password</label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                  placeholder="Any value works in mock"
                />
              </div>
              <button
                type="submit"
                className="w-full bg-theater-primary text-white py-2 rounded-md hover:bg-theater-primary/90 transition-colors"
                disabled={loading}
              >
                {loading ? 'Logging in...' : 'Login'}
              </button>
            </form>
            
            {googleAvailable && (
              <>
                <div className="my-6 flex items-center">
                  <div className="flex-1 border-t border-gray-300"></div>
                  <div className="px-4 text-sm text-gray-500">Or</div>
                  <div className="flex-1 border-t border-gray-300"></div>
                </div>
                
                <button
                  onClick={loginWithGoogle}
                  disabled={loading}
                  className="w-full bg-white border border-gray-300 text-gray-700 py-2 px-4 rounded-md hover:bg-gray-50 transition-colors flex items-center justify-center space-x-2"
                >
                  <svg className="w-5 h-5" viewBox="0 0 24 24">
                    <path fill="#4285f4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                    <path fill="#34a853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                    <path fill="#fbbc05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                    <path fill="#ea4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                  </svg>
                  <span>Continue with Google</span>
                </button>
              </>
            )}
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 bg-white rounded-lg shadow-lg p-6">
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-2xl font-bold text-theater-dark">Book Your Tickets</h2>
                <button onClick={logout} className="text-sm text-red-600 hover:underline">Logout</button>
              </div>

              {error && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
                  {error}
                </div>
              )}

              <form onSubmit={handleBooking} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Select Show</label>
                  <select
                    value={selectedShow}
                    onChange={(e) => { setSelectedShow(e.target.value); setSelectedShowTime(''); }}
                    className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                    required
                  >
                    <option value="">-- Choose a show --</option>
                    {shows.map((show) => (
                      <option key={show.id} value={show.id}>{show.title}</option>
                    ))}
                  </select>
                </div>

                {selectedShow && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Select Show Time</label>
                    <select
                      value={selectedShowTime}
                      onChange={(e) => setSelectedShowTime(e.target.value)}
                      className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                      required
                    >
                      <option value="">-- Choose a time --</option>
                      {shows.find(s => s.id === selectedShow)?.showTimes.map((st) => (
                        <option key={st.id} value={st.id}>
                          {st.date} at {st.time} ({st.availableSeats} left)
                        </option>
                      ))}
                    </select>
                  </div>
                )}

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Your Name</label>
                    <input
                      type="text"
                      value={customerName}
                      onChange={(e) => setCustomerName(e.target.value)}
                      placeholder={user?.name}
                      className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Phone Number</label>
                    <input
                      type="tel"
                      value={customerPhone}
                      onChange={(e) => setCustomerPhone(e.target.value)}
                      className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Number of Tickets</label>
                  <input
                    type="number"
                    min={1}
                    value={tickets}
                    onChange={(e) => setTickets(parseInt(e.target.value, 10))}
                    className="mt-1 w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-theater-primary focus:border-theater-primary"
                    required
                  />
                </div>

                <button
                  type="submit"
                  className="w-full bg-theater-primary text-white py-3 rounded-md hover:bg-theater-primary/90 transition-colors"
                >
                  Confirm Booking
                </button>
              </form>
            </div>

            <div className="lg:col-span-1">
              <div className="bg-white rounded-lg shadow-lg p-6">
                <h3 className="text-xl font-bold text-theater-dark mb-4">Booking Summary</h3>
                <div className="space-y-2 text-sm">
                  <p><span className="font-medium">User:</span> {user?.email}</p>
                  <p><span className="font-medium">Show:</span> {shows.find(s => s.id === selectedShow)?.title || '-'}</p>
                  <p><span className="font-medium">Show Time:</span> {shows.find(s => s.id === selectedShow)?.showTimes.find(st => st.id === selectedShowTime)?.date || '-'} {shows.find(s => s.id === selectedShow)?.showTimes.find(st => st.id === selectedShowTime)?.time || ''}</p>
                  <p><span className="font-medium">Tickets:</span> {tickets}</p>
                  <p><span className="font-medium">Total Price:</span> ${
                    (shows.find(s => s.id === selectedShow)?.ticketPrice || 0) * tickets
                  }</p>
                </div>
              </div>

              {bookingResult && (
                <div className="bg-white rounded-lg shadow-lg p-6 mt-6">
                  <h3 className="text-xl font-bold text-theater-dark mb-4">Booking Confirmed!</h3>
                  <p className="mb-2 text-sm">Booking ID: <span className="font-mono">{bookingResult.id}</span></p>
                  <p className="mb-4 text-sm">Show: {shows.find(s => s.id === selectedShow)?.title}</p>
                  {bookingResult.qrCode && (
                    <img src={bookingResult.qrCode} alt="Ticket QR Code" className="w-full" />
                  )}
                  <p className="mt-4 text-sm text-gray-600">Show this QR code at the entrance for scanning.</p>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default BookingPage;

