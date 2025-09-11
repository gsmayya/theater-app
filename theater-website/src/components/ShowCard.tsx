import React from 'react';
import { Show, ShowTime } from '../types/show';

interface ShowCardProps {
  show: Show;
  onBookShow?: (show: Show) => void;
}

const ShowCard: React.FC<ShowCardProps> = ({ show, onBookShow }) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const formatTime = (timeString: string) => {
    return new Date(`2000-01-01T${timeString}`).toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    });
  };

  const getAvailabilityColor = (availableSeats: number, totalSeats: number) => {
    const percentage = (availableSeats / totalSeats) * 100;
    if (percentage > 50) return 'text-green-600';
    if (percentage > 25) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition-shadow duration-300">
      {/* Show Image Placeholder */}
      <div className="h-48 bg-gradient-to-r from-theater-primary to-theater-secondary flex items-center justify-center">
        <div className="text-6xl">🎭</div>
      </div>
      
      {/* Show Details */}
      <div className="p-6">
        <div className="flex justify-between items-start mb-4">
          <div>
            <h3 className="text-xl font-bold text-theater-dark mb-2">{show.title || show.name}</h3>
            <span className="inline-block bg-theater-primary/10 text-theater-primary px-3 py-1 rounded-full text-sm font-medium">
              Theater Show
            </span>
          </div>
          <div className="text-right">
            <p className="text-2xl font-bold text-theater-primary">${(show.price / 100).toFixed(2)}</p>
            <p className="text-sm text-gray-500">per ticket</p>
          </div>
        </div>
        
        <p className="text-gray-600 mb-4 line-clamp-3">{show.details}</p>
        
        <div className="space-y-2 mb-4">
          <div className="flex items-center text-sm text-gray-600">
            <span className="font-medium mr-2">Location:</span> {show.location}
          </div>
          <div className="flex items-center text-sm text-gray-600">
            <span className="font-medium mr-2">Show Number:</span> {show.show_number}
          </div>
          <div className="flex items-center text-sm text-gray-600">
            <span className="font-medium mr-2">Available Tickets:</span> {show.total_tickets - show.booked_tickets} of {show.total_tickets}
          </div>
        </div>
        
        {/* Show Times */}
        <div className="border-t pt-4">
          <p className="text-sm font-medium text-gray-700 mb-3">Upcoming Shows:</p>
          <div className="space-y-2">
            {show.showTimes.slice(0, 3).map((showTime: ShowTime) => (
              <div
                key={showTime.id}
                className="flex justify-between items-center p-3 bg-gray-50 rounded-lg"
              >
                <div>
                  <p className="font-medium text-theater-dark">
                    {formatDate(showTime.date)}
                  </p>
                  <p className="text-sm text-gray-600">{formatTime(showTime.time)}</p>
                </div>
                <div className="text-right">
                  <p className={`text-sm font-medium ${getAvailabilityColor(showTime.availableSeats || (show.total_tickets - show.booked_tickets), showTime.totalSeats || show.total_tickets)}`}>
                    {showTime.availableSeats || (show.total_tickets - show.booked_tickets)} seats left
                  </p>
                  <p className="text-xs text-gray-500">
                    of {showTime.totalSeats || show.total_tickets} total
                  </p>
                </div>
              </div>
            ))}
            {show.showTimes.length > 3 && (
              <p className="text-sm text-gray-500 text-center py-2">
                +{show.showTimes.length - 3} more shows available
              </p>
            )}
          </div>
        </div>
        
        {/* Action Button */}
        <div className="mt-6">
          <button
            onClick={() => onBookShow?.(show)}
            className="w-full bg-theater-primary text-white py-3 px-4 rounded-lg font-medium hover:bg-theater-primary/90 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-theater-primary focus:ring-offset-2"
          >
            Book Tickets
          </button>
        </div>
      </div>
    </div>
  );
};

export default ShowCard;
