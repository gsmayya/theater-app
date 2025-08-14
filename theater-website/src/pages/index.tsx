import { useState } from 'react';
import TicketForm from '../components/TicketForm';
import { TicketDetails } from '../utils/TicketDetails';

const HomePage = () => {
  const [ticketDetails, setTicketDetails] = useState<TicketDetails | null>(null);

  const handleTicketSubmission = (details: TicketDetails) => {
    setTicketDetails(details);
  };

  return (
    <div>
      <h1>Welcome to the Ticket Booking System</h1>
      <TicketForm onSubmit={handleTicketSubmission} />
      {ticketDetails && (
        <div>
          <h2>Ticket Details</h2>
          <p>Mobile Number: {ticketDetails.mobileNumber}</p>
          <p>Number of Tickets: {ticketDetails.numberOfTickets}</p>
          <p>Name: {ticketDetails.name}</p>
          <p>Payment Mobile Number: {ticketDetails.paymentMobileNumber}</p>
          <p>Date and Time: {ticketDetails.dateTime}</p>
          <p>Item ID: {ticketDetails.itemId}</p>
          <p>Unique ID: {ticketDetails.uniqueId}</p>
        </div>
      )}
    </div>
  );
};

export default HomePage;
