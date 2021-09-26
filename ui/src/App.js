import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";

import LoginPage from "./components/LoginPage";
import CreateAccountPage from "./components/CreateAccountPage";
import BrandBar from "./components/BrandBar";
import HomePage from "./components/HomePage";

function App() {
  return (
    <div className="App">
      <BrandBar/>
      <Router>
        <Switch>
          <Route path="/create">
            <CreateAccountPage/>
          </Route>
          <Route path="/login">
            <LoginPage/>
          </Route>
          <Route path="/">
            <HomePage/>
          </Route>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
