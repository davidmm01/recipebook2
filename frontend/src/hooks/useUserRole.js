import { useState, useEffect } from 'react';
import { doc, onSnapshot } from 'firebase/firestore';
import { db } from '../firebase';
import { USER_ROLES, hasRole } from '../utils/userUtils';

/**
 * Hook to get and subscribe to the current user's role
 * @param {Object} user - Firebase Auth user object
 * @returns {Object} { role, loading, isViewer, isEditor, isAdmin, canEdit, canManageUsers }
 */
export function useUserRole(user) {
  const [role, setRole] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) {
      setRole(null);
      setLoading(false);
      return;
    }

    // Subscribe to real-time updates of the user's role
    const userRef = doc(db, 'users', user.uid);
    const unsubscribe = onSnapshot(
      userRef,
      (doc) => {
        if (doc.exists()) {
          setRole(doc.data().role || USER_ROLES.VIEWER);
        } else {
          setRole(USER_ROLES.VIEWER);
        }
        setLoading(false);
      },
      (error) => {
        console.error('Error fetching user role:', error);
        setRole(USER_ROLES.VIEWER); // Default to viewer on error
        setLoading(false);
      }
    );

    return () => unsubscribe();
  }, [user]);

  return {
    role,
    loading,
    // Convenience booleans
    isViewer: role === USER_ROLES.VIEWER,
    isEditor: role === USER_ROLES.EDITOR,
    isAdmin: role === USER_ROLES.ADMIN,
    // Permission checks
    canEdit: role && hasRole(role, USER_ROLES.EDITOR),
    canManageUsers: role && hasRole(role, USER_ROLES.ADMIN),
  };
}
