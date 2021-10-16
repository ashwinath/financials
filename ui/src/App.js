import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
} from "react-router-dom";

import { Header, PublicRoute, PrivateRoute } from "./components";
import { HomePage, LoginPage, CreateAccountPage } from "./pages";
import { InvestmentsTradesPage, InvestmentsMainPage } from "./pages/investments";

function App() {
  return (
    <div className="App">
      <Router>
        <Header/>
        <Switch>
          <PublicRoute path="/create">
            <CreateAccountPage/>
          </PublicRoute>
          <PublicRoute path="/login">
            <LoginPage/>
          </PublicRoute>
          <PrivateRoute path="/investments/trades">
            <InvestmentsTradesPage/>
          </PrivateRoute>
          <PrivateRoute path="/investments">
            <InvestmentsMainPage/>
          </PrivateRoute>
          <PrivateRoute path="/">
            <HomePage/>
          </PrivateRoute>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
