import {Route, Routes} from "react-router-dom";
import Login from "./pages/login/login";
import Main from "./pages/main/main";

function App() {
  return (
      <div className="App">
        <Routes>
          <Route path="/login" element={<Login/>}/>
          <Route path="/" element={<Main/>}/>
        </Routes>
      </div>
  );
}

export default App;
