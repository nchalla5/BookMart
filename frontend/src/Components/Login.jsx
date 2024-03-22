import React, { useState } from 'react';
import SignIn from './SignIn';
import SignUp from './SignUp';
import './Login.css';

function Login() {
  const [isSignUp, setIsSignUp] = useState(false); // State to toggle between forms

  const toggleForm = () => setIsSignUp(!isSignUp); // Toggle function

  return (
    <div className="auth-page">
      {isSignUp ? <SignUp onToggle={toggleForm} /> : <SignIn onToggle={toggleForm} />}
    </div>
  );
}

export default Login;
