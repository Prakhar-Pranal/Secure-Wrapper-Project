// API Configuration
const API_BASE_URL = 'http://localhost:8080';

// Health Check
export async function checkAPIHealth() {
  try {
    const response = await fetch(`${API_BASE_URL}/health`);
    return response.ok;
  } catch {
    return false;
  }
}

// Login
export async function login(email, password) {
  try {
    const response = await fetch(`${API_BASE_URL}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      throw new Error('Invalid credentials');
    }

    const data = await response.json();
    if (data.token) {
      localStorage.setItem('authToken', data.token);
      return data;
    }
    throw new Error(data.error || 'Login failed');
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
}

// Register Access Rule with Backend
export async function registerAccessRule(fileId, password, ip, mac) {
  try {
    const response = await fetch(`${API_BASE_URL}/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        file_id: fileId,
        password,
        ip,
        mac,
      }),
    });

    if (!response.ok) {
      throw new Error('Failed to register access rule');
    }

    return await response.json();
  } catch (error) {
    console.error('Register error:', error);
    throw error;
  }
}

// Verify Access
export async function verifyAccess(fileId, password, ip, mac) {
  try {
    const response = await fetch(`${API_BASE_URL}/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        file_id: fileId,
        password,
        ip,
        mac,
      }),
    });

    const data = await response.json();
    if (data.status === 'allowed') {
      return true;
    }
    return false;
  } catch (error) {
    console.error('Verify error:', error);
    throw error;
  }
}

// Pre-check (check if environment matches before asking for password)
export async function preCheckEnvironment(fileId, ip, mac) {
  try {
    const response = await fetch(`${API_BASE_URL}/pre-check`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        file_id: fileId,
        ip,
        mac,
      }),
    });

    const data = await response.json();
    return data.status === 'allowed';
  } catch (error) {
    console.error('Pre-check error:', error);
    throw error;
  }
}

// Logout
export function logout() {
  localStorage.removeItem('authToken');
}

// Get Auth Token
export function getAuthToken() {
  return localStorage.getItem('authToken');
}

// Check if authenticated
export function isAuthenticated() {
  return !!getAuthToken();
}
