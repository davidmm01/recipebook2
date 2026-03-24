import React, { useState, useEffect } from 'react';
import { getAdminUsers, updateUserRole } from '../utils/api';
import { auth } from '../firebase';

const ROLES = ['viewer', 'editor', 'admin'];

// Icons for saved vs pending state
const SavedIcon = () => (
  <span title="Saved" style={{ color: '#28a745', fontWeight: 'bold', fontSize: '14px' }}>✓</span>
);
const PendingIcon = () => (
  <span title="Unsaved change" style={{ color: '#e67e00', fontWeight: 'bold', fontSize: '14px' }}>●</span>
);

function UserManagement({ onClose }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  // pendingChanges: { [firebaseUid]: newRole }
  const [pendingChanges, setPendingChanges] = useState({});

  const currentUid = auth.currentUser?.uid;
  const hasPendingChanges = Object.keys(pendingChanges).length > 0;

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getAdminUsers();
      setUsers(data);
      setPendingChanges({});
    } catch (err) {
      setError('Failed to load users');
    } finally {
      setLoading(false);
    }
  };

  const handleRoleSelect = (uid, newRole) => {
    const savedRole = users.find(u => u.firebaseUid === uid)?.role;
    setPendingChanges(prev => {
      // If reverted back to saved value, remove from pending
      if (newRole === savedRole) {
        const { [uid]: _, ...rest } = prev;
        return rest;
      }
      return { ...prev, [uid]: newRole };
    });
  };

  const handleSave = async () => {
    setSaving(true);
    setError(null);
    const entries = Object.entries(pendingChanges);
    const failed = [];

    for (const [uid, role] of entries) {
      try {
        await updateUserRole(uid, role);
      } catch (err) {
        failed.push(uid);
      }
    }

    if (failed.length > 0) {
      setError(`Failed to save changes for ${failed.length} user(s). Refreshing...`);
      await loadUsers();
    } else {
      // Apply pending changes to users state and clear pending
      setUsers(prev => prev.map(u =>
        pendingChanges[u.firebaseUid] !== undefined
          ? { ...u, role: pendingChanges[u.firebaseUid] }
          : u
      ));
      setPendingChanges({});
    }
    setSaving(false);
  };

  const getDisplayRole = (user) =>
    pendingChanges[user.firebaseUid] !== undefined
      ? pendingChanges[user.firebaseUid]
      : user.role;

  const isPending = (user) =>
    pendingChanges[user.firebaseUid] !== undefined &&
    pendingChanges[user.firebaseUid] !== user.role;

  return (
    <div style={styles.overlay} onClick={onClose}>
      <div style={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div style={styles.header}>
          <h2 style={{ margin: 0 }}>Manage Users</h2>
          <button onClick={onClose} style={styles.closeButton}>×</button>
        </div>

        {error && <div style={styles.errorMessage}>{error}</div>}

        <div style={styles.content}>
          {loading ? (
            <p>Loading...</p>
          ) : (
            <table style={styles.table}>
              <thead>
                <tr>
                  <th style={styles.th}>Email</th>
                  <th style={styles.th}>Display Name</th>
                  <th style={styles.th}>Role</th>
                  <th style={styles.th}></th>
                </tr>
              </thead>
              <tbody>
                {users.map(user => {
                  const isSelf = user.firebaseUid === currentUid;
                  const pending = isPending(user);
                  return (
                    <tr key={user.firebaseUid} style={styles.tr}>
                      <td style={styles.td}>{user.email}</td>
                      <td style={styles.td}>{user.displayName || <em style={{ color: '#999' }}>—</em>}</td>
                      <td style={styles.td}>
                        <select
                          value={getDisplayRole(user)}
                          disabled={isSelf || saving}
                          onChange={(e) => handleRoleSelect(user.firebaseUid, e.target.value)}
                          style={{
                            ...styles.select,
                            ...(isSelf ? styles.selectDisabled : {}),
                            ...(pending ? styles.selectPending : {}),
                          }}
                        >
                          {ROLES.map(r => (
                            <option key={r} value={r}>{r}</option>
                          ))}
                        </select>
                      </td>
                      <td style={{ ...styles.td, width: '24px', textAlign: 'center' }}>
                        {!isSelf && (pending ? <PendingIcon /> : <SavedIcon />)}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>

        <div style={styles.footer}>
          <button onClick={onClose} style={styles.closeBtn} disabled={saving}>Close</button>
          <button
            onClick={handleSave}
            disabled={!hasPendingChanges || saving}
            style={{
              ...styles.saveBtn,
              ...(!hasPendingChanges || saving ? styles.saveBtnDisabled : {}),
            }}
          >
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </div>
    </div>
  );
}

const styles = {
  overlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    backgroundColor: 'white',
    borderRadius: '8px',
    width: '90%',
    maxWidth: '700px',
    maxHeight: '80vh',
    display: 'flex',
    flexDirection: 'column',
    boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '20px',
    borderBottom: '1px solid #e0e0e0',
    flexShrink: 0,
  },
  closeButton: {
    background: 'none',
    border: 'none',
    fontSize: '32px',
    cursor: 'pointer',
    color: '#666',
    padding: '0',
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  content: {
    padding: '20px',
    overflow: 'auto',
    flex: 1,
  },
  table: {
    width: '100%',
    borderCollapse: 'collapse',
    fontSize: '14px',
  },
  th: {
    textAlign: 'left',
    padding: '8px 12px',
    borderBottom: '2px solid #e0e0e0',
    color: '#555',
    fontWeight: '600',
  },
  tr: {
    borderBottom: '1px solid #f0f0f0',
  },
  td: {
    padding: '10px 12px',
    verticalAlign: 'middle',
  },
  select: {
    padding: '4px 8px',
    fontSize: '14px',
    borderRadius: '4px',
    border: '1px solid #ddd',
    cursor: 'pointer',
  },
  selectDisabled: {
    backgroundColor: '#f5f5f5',
    color: '#999',
    cursor: 'not-allowed',
  },
  selectPending: {
    borderColor: '#e67e00',
    outline: '1px solid #e67e00',
  },
  footer: {
    padding: '16px 20px',
    borderTop: '1px solid #e0e0e0',
    display: 'flex',
    justifyContent: 'flex-end',
    gap: '10px',
    flexShrink: 0,
  },
  closeBtn: {
    padding: '8px 20px',
    fontSize: '14px',
    fontWeight: '500',
    color: '#666',
    backgroundColor: 'white',
    border: '1px solid #ddd',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  saveBtn: {
    padding: '8px 20px',
    fontSize: '14px',
    fontWeight: '500',
    color: 'white',
    backgroundColor: '#007bff',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
  saveBtnDisabled: {
    backgroundColor: '#b0c4de',
    cursor: 'not-allowed',
  },
  errorMessage: {
    margin: '16px 20px 0',
    padding: '10px',
    backgroundColor: '#fee',
    color: '#c33',
    borderRadius: '4px',
    fontSize: '14px',
  },
};

export default UserManagement;
