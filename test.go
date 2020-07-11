package main

import (
    "fmt"
    "net"
    "os"
    "log"
    "encoding/csv"
    "encoding/json"
    "strconv"
    "bufio"
    "strings"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "4040"
    CONN_TYPE = "tcp"
)

type CovidData struct {
  Date      string
  Positive  int
  Tests     int
  Expired   int
  Admitted  int
  Discharge int
  Region    string
}

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  // Reading query data from connection (sent from User).
  query, err := bufio.NewReader(conn).ReadString('\n')
  // Check if there was any error.
  if err != nil {
    fmt.Println(err)
    return
  }
  queryJson := strings.TrimSpace(string(query))

  fmt.Println(queryJson) // Prints query data
  var result map[string]interface{}
  json.Unmarshal([]byte(queryJson), &result)
  // If user sent query based on Date, get date.
  date := result["query"].(map[string]interface{})["date"]
  // OR
  // If user sent query based on Region, get region.
  region := result["query"].(map[string]interface{})["region"]

  // Reading covid data from CSV file.
  filepath := "covid.csv"
  openfile, err := os.Open(filepath)
  checkError("Error in opening file\n",err)
  filedata, err := csv.NewReader(openfile).ReadAll()
  checkError("Error in reading the file\n",err)

  var dataPerDay CovidData
  var matchedData []CovidData
  // In (key,value) pair, If you only need the second item in the range (the value), 
  // Then use the blank identifier, an underscore, to discard the first (the key).
  for _, value:=range filedata{
    if (date==value[2] && region==nil) || (date==nil && region==value[5]){
      dataPerDay.Positive, _ = strconv.Atoi(value[0])
      dataPerDay.Tests, _ = strconv.Atoi(value[1])
      dataPerDay.Date = value[2]
      dataPerDay.Discharge, _ = strconv.Atoi(value[3])
      dataPerDay.Expired, _ = strconv.Atoi(value[4])
      dataPerDay.Region = value[5]
      dataPerDay.Admitted, _ = strconv.Atoi(value[6])
      //Append matched data
      matchedData = append(matchedData, dataPerDay)
    }
  }

  match := map[string]interface{}{
    "response" : matchedData,
  }
  // Convert Matched Data to JSON format with proper indentation for response.
  jsonData, err := json.MarshalIndent(match,"","  ")
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  //fmt.Println(string(jsonData))
  // Send a JSON Response back to person contacting us.
  conn.Write([]byte(string(jsonData)))
  // Close the connection when you're done with it.
  conn.Close()
}

// Function to handle errors
func checkError(msg string, err error){
  if err != nil{
    log.Fatal(msg, err)
  }
}