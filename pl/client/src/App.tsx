import React, { useState, useEffect } from 'react';
import { Github, Twitter, Linkedin, ArrowRight } from 'lucide-react';
import TaskManager from './TaskManager';
import AuthService from './services/auth';

function App() {
  const [showAuth, setShowAuth] = useState(false);
  const [isLogin, setIsLogin] = useState(true);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [authError, setAuthError] = useState('');

  // Check for existing session on mount
  useEffect(() => {
    const checkAuth = () => {
      const isAuth = AuthService.isAuthenticated();
      setIsLoggedIn(isAuth);
    };

    // checkAuth();
    // Set up token refresh interval
    const refreshInterval = setInterval(async () => {
      if (isLoggedIn) {
        const success = await AuthService.refreshAccessToken();
        if (!success) {
          setIsLoggedIn(false);
        }
      }
    }, 4 * 60 * 1000); // Refresh every 4 minutes

    return () => clearInterval(refreshInterval);
  }, [isLoggedIn]);

  const handleLogin = async (email: string, password: string) => {
    try {
      const success = await AuthService.login(email, password);
      if (success) {
        setIsLoggedIn(true);
        setShowAuth(false);
        setAuthError('');
      } else {
        setAuthError('Invalid credentials');
      }
    } catch (error) {
      setAuthError('An error occurred during login');
    }
  };

  const handleRegister = async (email: string, password: string) => {
    try {
      await AuthService.register(email, password);
      setIsLoggedIn(true);
      setShowAuth(false);
      setAuthError('');
    } catch (error) {
      setAuthError(error instanceof Error ? error.message : 'Registration failed');
    }
  };

  const handleLogout = async () => {
    await AuthService.logout();
    setIsLoggedIn(false);
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-900 to-gray-800 text-white">
      {isLoggedIn ? (
        <>
          <nav className="bg-gray-800 shadow-lg">
            <div className="container mx-auto px-4 py-4 flex justify-between items-center">
              <div className="text-2xl font-bold">DevSpace</div>
              <button 
                onClick={handleLogout}
                className="px-4 py-2 rounded-lg bg-gray-700 hover:bg-gray-600 transition"
              >
                Logout
              </button>
            </div>
          </nav>
          <TaskManager />
        </>
      ) : (
        <>
          {/* Auth Modal */}
          {showAuth && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
              <div className="bg-gray-800 rounded-lg p-6 w-full max-w-md">
                <h2 className="text-2xl font-bold mb-6">{isLogin ? 'Login' : 'Register'}</h2>
                {authError && (
                  <div className="bg-red-500 bg-opacity-10 border border-red-500 text-red-500 px-4 py-2 rounded-lg mb-4">
                    {authError}
                  </div>
                )}
                <form onSubmit={(e) => {
                  e.preventDefault();
                  const form = e.target as HTMLFormElement;
                  const email = (form.elements.namedItem('email') as HTMLInputElement).value;
                  const password = (form.elements.namedItem('password') as HTMLInputElement).value;
                  if (isLogin) {
                    handleLogin(email, password);
                  } else {
                    handleRegister(email, password);
                  }
                }}>
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium mb-1">Email</label>
                      <input
                        type="email"
                        name="email"
                        className="w-full bg-gray-700 rounded-lg px-3 py-2 text-white"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium mb-1">Password</label>
                      <input
                        type="password"
                        name="password"
                        className="w-full bg-gray-700 rounded-lg px-3 py-2 text-white"
                        required
                      />
                    </div>
                  </div>
                  <div className="mt-6 flex flex-col gap-3">
                    <div className="flex justify-end gap-3">
                      <button
                        type="button"
                        onClick={() => setShowAuth(false)}
                        className="px-4 py-2 rounded-lg border border-gray-600 hover:border-gray-500"
                      >
                        Cancel
                      </button>
                      <button
                        type="submit"
                        className="bg-blue-500 hover:bg-blue-600 px-4 py-2 rounded-lg font-medium"
                      >
                        {isLogin ? 'Login' : 'Register'}
                      </button>
                    </div>
                    <button
                      type="button"
                      onClick={() => setIsLogin(!isLogin)}
                      className="text-sm text-gray-400 hover:text-white transition"
                    >
                      {isLogin ? "Don't have an account? Register" : 'Already have an account? Login'}
                    </button>
                  </div>
                </form>
              </div>
            </div>
          )}

          {/* Hero Section */}
          <div className="container mx-auto px-4 py-20">
            <nav className="flex justify-between items-center mb-16">
              <div className="text-2xl font-bold">DevSpace</div>
              <div className="flex gap-8">
                <a href="#features" className="hover:text-blue-400 transition">Features</a>
                <a href="#pricing" className="hover:text-blue-400 transition">Pricing</a>
                <a href="#contact" className="hover:text-blue-400 transition">Contact</a>
                <button 
                  onClick={() => {
                    setIsLogin(true);
                    setShowAuth(true);
                  }}
                  className="hover:text-blue-400 transition"
                >
                  Login
                </button>
              </div>
            </nav>

            <div className="max-w-4xl mx-auto text-center">
              <h1 className="text-6xl font-bold mb-6 bg-gradient-to-r from-blue-400 to-purple-500 text-transparent bg-clip-text">
                Build Better Software Together
              </h1>
              <p className="text-xl text-gray-300 mb-10">
                The all-in-one collaboration platform for developers. Code reviews, project management, and team communication in one place.
              </p>
              <div className="flex gap-4 justify-center">
                <button 
                  onClick={() => {
                    setIsLogin(false);
                    setShowAuth(true);
                  }}
                  className="bg-blue-500 hover:bg-blue-600 px-8 py-3 rounded-lg font-medium flex items-center gap-2 transition"
                >
                  Get Started <ArrowRight size={20} />
                </button>
                <button className="border border-gray-600 hover:border-gray-500 px-8 py-3 rounded-lg font-medium transition">
                  Live Demo
                </button>
              </div>
            </div>
          </div>

          {/* Features Section */}
          <div id="features" className="container mx-auto px-4 py-20">
            <h2 className="text-4xl font-bold text-center mb-16">Why Choose DevSpace?</h2>
            <div className="grid md:grid-cols-3 gap-8">
              {features.map((feature, index) => (
                <div key={index} className="bg-gray-800 p-6 rounded-lg hover:bg-gray-700 transition">
                  <div className="text-blue-400 mb-4">{feature.icon}</div>
                  <h3 className="text-xl font-semibold mb-3">{feature.title}</h3>
                  <p className="text-gray-400">{feature.description}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Social Proof */}
          <div className="bg-gray-800">
            <div className="container mx-auto px-4 py-16 text-center">
              <p className="text-2xl font-medium mb-8">Trusted by developers worldwide</p>
              <div className="flex justify-center gap-12">
                <a href="https://github.com/Adity-aprasad" target="_blank" rel="noopener noreferrer">
                  <Github size={32} className="text-gray-400 hover:text-white transition" />
                </a>
                <a href="https://twitter.com/your-profile" target="_blank" rel="noopener noreferrer">
                  <Twitter size={32} className="text-gray-400 hover:text-white transition" />
                </a>
                <a href="https://www.linkedin.com/in/aditya-prasad-081029228/" target="_blank" rel="noopener noreferrer">
                  <Linkedin size={32} className="text-gray-400 hover:text-white transition" />
                </a>
              </div>
            </div>
          </div>

          {/* Footer */}
          <footer className="container mx-auto px-4 py-8">
            <div className="flex justify-between items-center">
              <div className="text-gray-400">© 2025 DevSpace. All rights reserved.</div>
              <div className="flex gap-6">
                <a href="#" className="text-gray-400 hover:text-white transition">Privacy</a>
                <a href="#" className="text-gray-400 hover:text-white transition">Terms</a>
                <a href="#" className="text-gray-400 hover:text-white transition">Contact</a>
              </div>
            </div>
          </footer>
        </>
      )}
    </div>
  );
}

const features = [
  {
    icon: '🚀',
    title: 'Lightning Fast',
    description: 'Built for speed and efficiency. Get your work done faster than ever before.'
  },
  {
    icon: '🔒',
    title: 'Secure by Design',
    description: 'Enterprise-grade security with end-to-end encryption and compliance features.'
  },
  {
    icon: '🤝',
    title: 'Team Collaboration',
    description: 'Work together seamlessly with real-time collaboration tools and integrations.'
  }
];

export default App;