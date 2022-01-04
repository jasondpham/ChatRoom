import React from 'react';
import './App.css';

function App() {
  return (
    <div className="App">
      <main className="form-signin">
        <form>
          <h1 className="h3 mb-3 fw-normal">Please sign in</h1>
            <input type="email" className="form-control" id="floatingInput" placeholder="Email" />
            <input type="password" className="form-control" id="floatingPassword" placeholder="Password"/>
          <button className="w-100 btn btn-lg btn-primary" type="submit">Sign in</button>
          <p className="mt-5 mb-3 text-muted">&copy; 2017–2021</p>
        </form>
      </main>
    </div>
  );
}

export default App;
