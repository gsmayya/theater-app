import React, { useState } from 'react';
import { Show, ShowTime, CalendarEvent } from '../types/show';

interface CalendarProps {
  shows: Show[];
  onDateSelect?: (date: string, events: CalendarEvent[]) => void;
}

const Calendar: React.FC<CalendarProps> = ({ shows, onDateSelect }) => {
  const [currentDate, setCurrentDate] = useState(new Date());

  // Create calendar events from shows
  const calendarEvents: CalendarEvent[] = shows.flatMap(show =>
    show.showTimes.map(showTime => ({
      id: `${show.id}-${showTime.id}`,
      title: show.title,
      date: showTime.date,
      time: showTime.time,
      show_id: show.id,
      location: show.location,
      price: show.price,
      available_tickets: show.total_tickets - show.booked_tickets
    }))
  );

  const getDaysInMonth = (date: Date) => {
    return new Date(date.getFullYear(), date.getMonth() + 1, 0).getDate();
  };

  const getFirstDayOfMonth = (date: Date) => {
    return new Date(date.getFullYear(), date.getMonth(), 1).getDay();
  };

  const formatMonthYear = (date: Date) => {
    return date.toLocaleDateString('en-US', { 
      month: 'long', 
      year: 'numeric' 
    });
  };

  const navigateMonth = (direction: 'prev' | 'next') => {
    const newDate = new Date(currentDate);
    if (direction === 'prev') {
      newDate.setMonth(newDate.getMonth() - 1);
    } else {
      newDate.setMonth(newDate.getMonth() + 1);
    }
    setCurrentDate(newDate);
  };

  const getEventsForDate = (day: number): CalendarEvent[] => {
    const dateString = new Date(currentDate.getFullYear(), currentDate.getMonth(), day)
      .toISOString().split('T')[0];
    
    return calendarEvents.filter(event => event.date === dateString);
  };

  const handleDateClick = (day: number) => {
    const dateString = new Date(currentDate.getFullYear(), currentDate.getMonth(), day)
      .toISOString().split('T')[0];
    const events = getEventsForDate(day);
    
    if (events.length > 0 && onDateSelect && dateString) {
      onDateSelect(dateString, events);
    }
  };

  const isToday = (day: number) => {
    const today = new Date();
    const checkDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), day);
    return checkDate.toDateString() === today.toDateString();
  };

  const isPastDate = (day: number) => {
    const today = new Date();
    const checkDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), day);
    return checkDate < today;
  };

  const daysInMonth = getDaysInMonth(currentDate);
  const firstDayOfMonth = getFirstDayOfMonth(currentDate);
  const daysArray = Array.from({ length: daysInMonth }, (_, i) => i + 1);
  const emptyDays = Array.from({ length: firstDayOfMonth }, (_, i) => i);

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      {/* Calendar Header */}
      <div className="flex justify-between items-center mb-6">
        <button
          onClick={() => navigateMonth('prev')}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        
        <h2 className="text-xl font-bold text-theater-dark">
          {formatMonthYear(currentDate)}
        </h2>
        
        <button
          onClick={() => navigateMonth('next')}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>

      {/* Days of Week Header */}
      <div className="grid grid-cols-7 gap-1 mb-2">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div
            key={day}
            className="p-2 text-center text-sm font-medium text-gray-500"
          >
            {day}
          </div>
        ))}
      </div>

      {/* Calendar Grid */}
      <div className="grid grid-cols-7 gap-1">
        {/* Empty days for first week */}
        {emptyDays.map(day => (
          <div key={`empty-${day}`} className="h-12"></div>
        ))}
        
        {/* Days of month */}
        {daysArray.map(day => {
          const events = getEventsForDate(day);
          const hasEvents = events.length > 0;
          const isCurrentDay = isToday(day);
          const isPast = isPastDate(day);
          
          return (
            <div
              key={day}
              onClick={() => handleDateClick(day)}
              className={`
                h-12 flex flex-col items-center justify-center text-sm relative rounded cursor-pointer
                transition-colors duration-200
                ${isPast ? 'text-gray-400 cursor-not-allowed' : 'hover:bg-theater-primary/10'}
                ${isCurrentDay ? 'bg-theater-primary text-white' : ''}
                ${hasEvents && !isCurrentDay ? 'bg-theater-secondary/20 text-theater-dark font-medium' : ''}
                ${!hasEvents && !isCurrentDay && !isPast ? 'hover:bg-gray-100' : ''}
              `}
            >
              <span className="text-xs">{day}</span>
              {hasEvents && (
                <div className="flex gap-0.5 mt-0.5">
                  {events.slice(0, 3).map((_, index) => (
                    <div
                      key={index}
                      className={`w-1 h-1 rounded-full ${
                        isCurrentDay ? 'bg-white' : 'bg-theater-primary'
                      }`}
                    />
                  ))}
                  {events.length > 3 && (
                    <div className={`w-1 h-1 rounded-full ${
                      isCurrentDay ? 'bg-white' : 'bg-theater-accent'
                    }`} />
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* Legend */}
      <div className="mt-6 flex flex-wrap gap-4 text-xs">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-theater-primary rounded"></div>
          <span className="text-gray-600">Today</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-theater-secondary/20 rounded"></div>
          <span className="text-gray-600">Shows Available</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 bg-gray-300 rounded"></div>
          <span className="text-gray-600">Past Dates</span>
        </div>
      </div>
    </div>
  );
};

export default Calendar;
