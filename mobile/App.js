import { StatusBar } from 'expo-status-bar';
import React, {useEffect, useState} from 'react';
import {StyleSheet, Text, View, Image, TouchableWithoutFeedback} from 'react-native';
import {useKeepAwake} from "expo-keep-awake";

export default function App() {
  useKeepAwake();

  let images = [
      "",
      "https://mir-s3-cdn-cf.behance.net/project_modules/fs/9afe0493484903.5e66500f8dea4.gif"
  ]

  let [image, setImage] = useState(0)
  let [status, setStatus] = useState("unknow")
  let [time, setTime] = useState(null)
  let [applications, setApplications] = useState([])

  useEffect(() => {
    let socket = new WebSocket("ws://192.168.0.108:3010/websocket");
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
    <View style={styles.container}>
      <View>
        <StatusBar hidden={true} />
        <Text style={{...styles.text, ...styles.header}}>Server status: <Text style={styles.bold}>{status}</Text> | {time}</Text>

        <View style={styles.list}>
          {
            applications.map(app => {
              let appParsed = JSON.parse(app)
              return(
                  <Text style={styles.text}>[ <Text
                      style={
                        appParsed.status === "SUCCESS" ? styles.success : styles.failed
                      }
                  >{appParsed.status}</Text> ]<Text style={styles.info}>[{appParsed.statusCode}]</Text>: <Text>{appParsed.link}</Text></Text>

              )
            })
          }
        </View>
      </View>

      <TouchableWithoutFeedback onPress={() => {
        let processable = image + 1
        if (processable > images.length){
          setImage(0)
          return
        }
        setImage(processable)
      }}>
        <View
            style={{flex: 1, justifyContent: "center", alignItems: "center"}}
        >
          <Image
              style ={{width: "80%", height:"80%", marginTop: 15, borderRadius: 10, overflow: "hidden"}}
              source={{ uri : images[image]}}
          />
        </View>
      </TouchableWithoutFeedback>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: 'black',
    justifyContent: 'center',
    alignItems: "center",
    flexDirection: "row",
    padding: 60
  },
  header: {
    paddingBottom: 15
  },
  text: {
    color: "white",
    fontSize: 18
  },
  list: {
    alignItems: 'flex-start',
  },
  bold: {
    fontWeight: "bold"
  },
  success: {
    backgroundColor: "green",
  },
  failed: {
    backgroundColor: "red"
  },
  info: {
    backgroundColor: "blue"
  }
});
