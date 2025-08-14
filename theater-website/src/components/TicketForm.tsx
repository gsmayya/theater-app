import React, { useState } from 'react';
import ItemDropdown from './ItemDropdown';
import { generateUniqueId } from '../utils/generateId';
import { createShareText } from '../utils/shareText';
import { TicketDetails } from '@/utils/TicketDetails';

interface TicketFormProps {
  onSubmit: (details: TicketDetails) => void;
}

const TicketForm: React.FC<TicketFormProps> = ({ onSubmit }) => {
  const [mobileNumber, setMobileNumber] = useState('');
  const [numberOfTickets, setNumberOfTickets] = useState(1);
  const [name, setName] = useState('');
  const [paymentMobileNumber, setPaymentMobileNumber] = useState('');
  const [selectedItem, setSelectedItem] = useState('');
  const [ticketId, setTicketId] = useState('');
  const [qrCodeData, setQrCodeData] = useState('');
  const [bookedBy, setBookedBy] = useState('');

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    var start = 1
    const id = generateUniqueId(selectedItem, start);
    start = start + 1
    setTicketId(id);

    const ticketDetails: TicketDetails = {
      mobileNumber,
      numberOfTickets,
      name,
      paymentMobileNumber,
      itemId: selectedItem,
      uniqueId: id,
      bookedBy,
      dateTime: new Date().toISOString(), // Or get the value from a controlled input if needed
    };

    setQrCodeData(createShareText(ticketDetails));

    // Pass ticketDetails to parent
    onSubmit(ticketDetails);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Mobile Number:</label>
        <input
          type="text"
          value={mobileNumber}
          onChange={(e) => setMobileNumber(e.target.value)}
          required
        />
      </div>
      <div>
        <label>Number of Tickets:</label>
        <input
          type="number"
          value={numberOfTickets}
          onChange={(e) => setNumberOfTickets(Number(e.target.value))}
          min="1"
          required
        />
      </div>
      <div>
        <label>Name:</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>
      <div>
        <label>Payment Mobile Number:</label>
        <input
          type="text"
          value={paymentMobileNumber}
          onChange={(e) => setPaymentMobileNumber(e.target.value)}
          required
        />
      </div>
       <div>
        <label>Booked by:</label>
        <input
          type="text"
          value={bookedBy}
          onChange={(e) => setBookedBy(e.target.value)}
          required
        />
      </div>
      <div>
        <label>Select Item:</label>
        <ItemDropdown onSelect={setSelectedItem} />
      </div>
      <button type="submit">Submit Ticket</button>
      {ticketId && (
        <div>
          <h3>Ticket ID: {ticketId}</h3>
          <p>Shareable Info: {qrCodeData}</p>
        </div>
      )}
    </form>
  );
};

export default TicketForm;
