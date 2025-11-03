import { useState, useEffect } from 'react';
import { getUserProfile } from '../utils/api';
import { USER_ROLES, hasRole } from '../utils/userUtils';

/**
 * Hook to get the current user's role from backend
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

    // Fetch user profile from backend API
    const fetchUserRole = async () => {
      try {
        const profile = await getUserProfile();
        setRole(profile.role || USER_ROLES.VIEWER);
      } catch (error) {
        console.error('Error fetching user role:', error);
        setRole(USER_ROLES.VIEWER); // Default to viewer on error
      } finally {
        setLoading(false);
      }
    };

    fetchUserRole();
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
