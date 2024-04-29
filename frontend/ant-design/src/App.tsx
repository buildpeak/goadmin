import { BrowserRouter } from "react-router-dom";
import MainRouter from "./routes/Router";

import "./App.css";

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <MainRouter />
      </BrowserRouter>
    </div>
  );
}

export default App;
