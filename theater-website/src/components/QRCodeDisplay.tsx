import React from 'react';
import { QRCode } from 'react-qrcode-logo';

interface QRCodeDisplayProps {
  ticketId: string;
  mobileNumber: string;
  numberOfTickets: number;
  name: string;
  paymentMobileNumber: string;
  dateTime: string;
  item: string;
}

const QRCodeDisplay: React.FC<QRCodeDisplayProps> = ({
  ticketId,
  mobileNumber,
  numberOfTickets,
  name,
  paymentMobileNumber,
  dateTime,
  item,
}) => {
  const ticketInfo = `Ticket ID: ${ticketId}\nMobile Number: ${mobileNumber}\nNumber of Tickets: ${numberOfTickets}\nName: ${name}\nPayment Mobile Number: ${paymentMobileNumber}\nDate and Time: ${dateTime}\nItem: ${item}`;

  return (
    <div>
      <h2>Your Ticket Information</h2>
      <pre>{ticketInfo}</pre>
      <QRCode value={ticketInfo} size={256} />
    </div>
  );
};

export default QRCodeDisplay;