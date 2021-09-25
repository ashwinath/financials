import React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
} from "react-router-dom";

import LoginPage from "./components/LoginPage";
import CreatePage from "./components/CreatePage";
import BrandBar from "./components/BrandBar";

function App() {
  return (
    <div className="App">
      <BrandBar/>
      <Router>
        <Switch>
          <Route path="/create">
            <CreatePage/>
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
