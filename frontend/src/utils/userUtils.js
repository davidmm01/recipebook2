/**
 * User roles in the system
 * - viewer: Can only view recipes
 * - editor: Can create and edit recipes
 * - admin: Full access to all recipes and user management
 */
export const USER_ROLES = {
  VIEWER: 'viewer',
  EDITOR: 'editor',
  ADMIN: 'admin',
};

/**
 * Checks if a user has a specific role or higher privileges
 * @param {string} userRole - The user's current role
 * @param {string} requiredRole - The required role to check against
 * @returns {boolean} Whether the user has sufficient privileges
 */
export function hasRole(userRole, requiredRole) {
  const roleHierarchy = {
    [USER_ROLES.VIEWER]: 1,
    [USER_ROLES.EDITOR]: 2,
    [USER_ROLES.ADMIN]: 3,
  };

  return roleHierarchy[userRole] >= roleHierarchy[requiredRole];
}
