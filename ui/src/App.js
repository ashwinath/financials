import React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
} from "react-router-dom";

import LoginPage from "./components/LoginPage";
import BrandBar from "./components/BrandBar";

function App() {
  return (
    <div className="App">
      <BrandBar/>
      <Router>
        <Switch>
          <Route path="/create">
            <h1>hello create</h1>
          </Route>
          <Route path="/login">
            <LoginPage/>
          </Route>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
