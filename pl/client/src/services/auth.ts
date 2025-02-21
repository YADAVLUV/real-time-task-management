class AuthService {
  private static API_BASE_URL = 'http://localhost:8080/auth';

  static async register(email: string, password: string): Promise<boolean> {
    try {
      const response = await fetch(`${this.API_BASE_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      return response.ok;
    } catch (error) {
      console.error('Registration error:', error);
      return false;
    }
  }

  static async login(email: string, password: string): Promise<boolean> {
      try {
        const response = await fetch(`${this.API_BASE_URL}/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password }),
          credentials: 'include' // Ensure cookies are sent with the request
        });

      return response.ok;
    } catch (error) {
      console.error('Login error:', error);
      return false;
    }
  }

  static async logout(): Promise<boolean> {
    try {
      await fetch(`${this.API_BASE_URL}/logout`, {
        method: 'POST',
        credentials: 'include' // Ensure cookies are cleared on logout
      });

      return true;
    } catch (error) {
      console.error('Logout error:', error);
      return false;
    }
  }

  static async isAuthenticated(): Promise<boolean> {
    try {
      return true;
      const response = await fetch('http://localhost:8080/auth/protected', {
        credentials: 'include'
      });

      return response.ok;
    } catch {
      return false;
    }
  }
}

export default AuthService;
