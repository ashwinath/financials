import React from 'react';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from "react-router-dom";

import { BrandBar } from "./components/BrandBar";
import { HomePage, LoginPage, CreateAccountPage } from "./pages";
import { InvestmentsTradesPage, InvestmentsMainPage } from "./pages/investments";

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
          <Route path="/investments/trades">
            <InvestmentsTradesPage/>
          </Route>
          <Route path="/investments">
            <InvestmentsMainPage/>
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
