import React from 'react';
import { useUserRole } from '../hooks/useUserRole';
import { auth } from '../firebase';

function RoleDisplay() {
  const user = auth.currentUser;
  const { role, loading, isAdmin, isEditor, isViewer } = useUserRole(user);

  if (loading) {
    return null;
  }

  if (!role) {
    return null;
  }

  const getRoleBadgeColor = () => {
    if (isAdmin) return '#dc3545';
    if (isEditor) return '#28a745';
    if (isViewer) return '#6c757d';
    return '#6c757d';
  };

  return (
    <div style={{
      display: 'inline-block',
      padding: '4px 12px',
      backgroundColor: getRoleBadgeColor(),
      color: '#fff',
      borderRadius: '12px',
      fontSize: '12px',
      fontWeight: '600',
      textTransform: 'uppercase',
      marginLeft: '12px'
    }}>
      {role}
    </div>
  );
}

export default RoleDisplay;
