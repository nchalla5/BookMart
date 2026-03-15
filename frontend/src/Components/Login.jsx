import React, { useState } from 'react';
import SignIn from './SignIn';
import SignUp from './SignUp';
import './Login.css';

function Login() {
  const [isSignUp, setIsSignUp] = useState(false); // State to toggle between forms

  const toggleForm = () => setIsSignUp(!isSignUp); // Toggle function

  return (
    <div className="auth-page">
      <div className="auth-shell">
        <section className="auth-hero">
          <p className="auth-eyebrow">Book Mart</p>
          <h1>Build your next bookshelf from someone else&apos;s last chapter.</h1>
          <p className="auth-copy">
            A full-stack marketplace for buying and selling used books, built with Go and React.
          </p>
          <div className="auth-highlights">
            <span>Secure login</span>
            <span>Local demo mode</span>
            <span>Buy and sell flow</span>
          </div>
        </section>
        <section className="auth-panel">
          {isSignUp ? <SignUp onToggle={toggleForm} /> : <SignIn onToggle={toggleForm} />}
        </section>
      </div>
    </div>
  );
}

export default Login;
