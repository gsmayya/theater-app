export const createShareText = (ticketDetails: { id?: any; mobileNumber: any; numberOfTickets: any; name: any; paymentMobileNumber?: string; selectedItem?: string; uniqueId?: any; itemBooked?: any; bookedBy?: any; }) => {
  const { uniqueId, mobileNumber, numberOfTickets, itemBooked, name, bookedBy } = ticketDetails;

  return `Ticket Details:
  Unique ID: ${uniqueId}
  Mobile Number: ${mobileNumber}
  Number of Tickets: ${numberOfTickets}
  Item Booked: ${itemBooked}
  Name: ${name}
  Booked By: ${bookedBy}`;
};