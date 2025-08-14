export const generateUniqueId = (showId: string, sequenceNumber: number): string => {
    return `${showId}-${sequenceNumber.toString().padStart(5, '0')}`;
};