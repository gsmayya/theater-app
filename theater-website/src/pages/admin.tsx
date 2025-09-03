import React from 'react';
import AdminLogin from '../components/AdminLogin';
import Layout from '../components/Layout';

const AdminPage: React.FC = () => {
  return (
    <Layout title="Admin Login - Illuminating Windows">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <AdminLogin />
      </div>
    </Layout>
  );
};

export default AdminPage;
