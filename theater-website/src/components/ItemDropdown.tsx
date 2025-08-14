import React, { useEffect, useState } from 'react';
import axios from 'axios';

const ItemDropdown: React.FC<{ onSelect: (itemId: string) => void }> = ({ onSelect }) => {
  const [items, setItems] = useState<{ id: string; name: string }[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchItems = async () => {
      try {
        // Dummy data for testing UI
        const response = { data: [
          { id: '1', name: 'Item One' },
          { id: '2', name: 'Item Two' },
          { id: '3', name: 'Item Three' }
        ] };
        setItems(response.data as { id: string; name: string }[]);
      } catch (err) {
        setError('Failed to fetch items');
      } finally {
        setLoading(false);
      }
    };

    fetchItems();
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>{error}</div>;

  return (
    <select onChange={(e) => onSelect(e.target.value)}>
      <option value="">Select an item</option>
      {items.map((item) => (
        <option key={item.id} value={item.id}>
          {item.name}
        </option>
      ))}
    </select>
  );
};

export default ItemDropdown;