// Mock user storage
let MOCK_USERS = [
  {
    id: '1',
    email: 'demo@example.com',
    password: 'demo123', // In a real app, passwords would be hashed
  }
];

// Mock token generation
const generateTokens = (user: { id: string; email: string }) => {
  const now = Date.now();
  const accessToken = btoa(JSON.stringify({
    user,
    exp: now + 15 * 60 * 1000, // 15 minutes
    iat: now,
  }));
  const refreshToken = btoa(JSON.stringify({
    user: { id: user.id },
    exp: now + 7 * 24 * 60 * 60 * 1000, // 7 days
    iat: now,
  }));
  return { accessToken, refreshToken };
};

export const mockAuthEndpoints = {
  '/api/auth/register': async (email: string, password: string) => {
    await new Promise(resolve => setTimeout(resolve, 500)); // Simulate network delay

    // Check if user already exists
    if (MOCK_USERS.some(user => user.email === email)) {
      throw new Error('Email already registered');
    }

    // Create new user
    const newUser = {
      id: String(MOCK_USERS.length + 1),
      email,
      password,
    };
    MOCK_USERS.push(newUser);

    // Return tokens for automatic login
    return generateTokens({
      id: newUser.id,
      email: newUser.email,
    });
  },

  '/api/auth/login': async (email: string, password: string) => {
    await new Promise(resolve => setTimeout(resolve, 500)); // Simulate network delay

    const user = MOCK_USERS.find(u => u.email === email && u.password === password);
    if (user) {
      return generateTokens({
        id: user.id,
        email: user.email,
      });
    }
    throw new Error('Invalid credentials');
  },

  '/api/auth/refresh': async (refreshToken: string) => {
    await new Promise(resolve => setTimeout(resolve, 500)); // Simulate network delay
    
    try {
      const decoded = JSON.parse(atob(refreshToken));
      if (decoded.exp > Date.now()) {
        const user = MOCK_USERS.find(u => u.id === decoded.user.id);
        if (!user) throw new Error('User not found');
        return generateTokens({
          id: user.id,
          email: user.email,
        });
      }
      throw new Error('Refresh token expired');
    } catch {
      throw new Error('Invalid refresh token');
    }
  },

  '/api/auth/logout': async () => {
    await new Promise(resolve => setTimeout(resolve, 500)); // Simulate network delay
    return { success: true };
  },
};