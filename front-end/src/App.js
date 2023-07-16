import logo from './logo.svg';
import './App.css';
import {useEffect, useState} from "react";

function App() {

  let [status, setStatus] = useState("unknow")
  let [time, setTime] = useState(null)
  let [applications, setApplications] = useState([])

  useEffect(() => {
    let socket = new WebSocket(process.env.REACT_APP_SERVER);
    socket.onopen = () => {
      setStatus("connected")

      let date = new Date(Date.now());
      setTime(date.toLocaleTimeString("en-US"))
    }

    socket.onmessage = function(event) {
      let applications = JSON.parse(event.data)
      setApplications(applications)

      let date = new Date(Date.now());
      setTime(date.toLocaleTimeString("en-US"))
    };

    socket.onclose = function(event) {
      setStatus("connection closed")
    };

    socket.onerror = function(error) {
      setStatus("error: " + error)
    };
  }, [])

  return (
    <div>
      <span>Status: <b>{status}</b> | {time}</span>
      <ul>
        {
          applications.map(app => {
            let appParsed = JSON.parse(app)
            return(
                <li>[ <b
                  className={
                    appParsed.status === "SUCCESS" ? 'green' : "red"
                  }
                >{appParsed.status}</b> ][{appParsed.statusCode}]: <a href={appParsed.link}>{appParsed.link}</a> ({appParsed.lastCheck})</li>
            )
          })
        }
      </ul>
    </div>
  );
}

export default App;
